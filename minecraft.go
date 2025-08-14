package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func writeVarInt(value int) []byte {
	buf := make([]byte, 0, 5)
	for {
		if (value & ^0x7F) == 0 {
			buf = append(buf, byte(value))
			break
		}
		buf = append(buf, byte((value&0x7F)|0x80))
		value = int(uint32(value) >> 7)
	}
	return buf
}

func readVarInt(r io.ByteReader) (int, error) {
	var value, position int
	for {
		currentByte, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		value |= (int(currentByte) & 0x7F) << position
		if (currentByte & 0x80) == 0 {
			break
		}
		position += 7
		if position >= 32 {
			return 0, errors.New("VarInt is too large")
		}
	}
	return value, nil
}

func createPacket(packetID int, data []byte) []byte {
	packetIDBuffer := writeVarInt(packetID)
	payload := append(packetIDBuffer, data...)
	lengthBuffer := writeVarInt(len(payload))
	return append(lengthBuffer, payload...)
}

func createHandshakePacket(host string, port uint16) []byte {
	const protocolVersion = 765
	var data bytes.Buffer
	data.Write(writeVarInt(protocolVersion))
	data.Write(writeVarInt(len(host)))
	data.WriteString(host)
	binary.Write(&data, binary.BigEndian, port)
	data.Write(writeVarInt(1))
	return createPacket(0x00, data.Bytes())
}

func getPingResult(host string, port uint16, timeout time.Duration, debugMode bool) (string, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	if debugMode {
		log.Printf("DEBUG [%s]: Initiating connection...", address)
	}

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		if debugMode {
			log.Printf("DEBUG [%s]: Connection error: %v", address, err)
		}
		return "", err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))
	if debugMode {
		log.Printf("DEBUG [%s]: Connection successful.", address)
	}

	handshakePacket := createHandshakePacket(host, port)
	if _, err = conn.Write(handshakePacket); err != nil {
		return "", err
	}

	statusRequestPacket := createPacket(0x00, []byte{})
	if _, err = conn.Write(statusRequestPacket); err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)

	_, err = readVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("error reading packet length: %w", err)
	}

	packetID, err := readVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("error reading packet id: %w", err)
	}
	if packetID != 0x00 {
		return "", fmt.Errorf("received unexpected packet ID: 0x%X", packetID)
	}

	jsonLength, err := readVarInt(reader)
	if err != nil {
		return "", fmt.Errorf("error reading json length: %w", err)
	}
	if jsonLength <= 0 {
		return "", errors.New("invalid JSON length")
	}

	jsonBytes := make([]byte, jsonLength)
	if _, err := io.ReadFull(reader, jsonBytes); err != nil {
		return "", fmt.Errorf("error reading JSON content: %w", err)
	}

	var statusObject map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &statusObject); err != nil {
		if debugMode {
			log.Printf("DEBUG [%s]: JSON unmarshal error: %v", address, err)
		}
		return "", fmt.Errorf("error unmarshaling JSON string: %w", err)
	}

	delete(statusObject, "favicon")

	if forgeData, ok := statusObject["forgeData"]; ok {
		if forgeDataMap, ok := forgeData.(map[string]interface{}); ok {
			delete(forgeDataMap, "d")
		}
	}

	modifiedJSONBytes, err := json.Marshal(statusObject)
	if err != nil {
		if debugMode {
			log.Printf("DEBUG [%s]: JSON marshal error: %v", address, err)
		}
		return "", fmt.Errorf("error marshaling JSON again: %w", err)
	}

	if debugMode {
		log.Printf("DEBUG [%s]: SUCCESS! Processed and cleaned JSON.", address)
	}

	return string(modifiedJSONBytes), nil
}