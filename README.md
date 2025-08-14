<div align="center">

# mcpscan

**M**ine**c**raft server **p**ort **s**canner.

<a href="https://github.com/tmquan2508/mcpscan/releases/latest">
    <img alt="Latest Release" src="https://img.shields.io/github/v/release/BẠN/mcpscan?style=for-the-badge&logo=github">
</a>

<br/>
<img alt="mcpscan icon" src="icon.png" height="200" width="200" style="border-radius: 50% 20% / 10% 40%;"/>
<br/>

*A fast, concurrent Minecraft server port scanner written in Go.*

</div>

## 🚀 Getting Started

The easiest way to use `mcpscan` is to download the latest pre-compiled binary for your operating system from the **[Releases Page](https://github.com/tmquan2508/mcpscan/releases/latest)**.

1.  Go to the [Releases page](https://github.com/BẠN/mcpscan/releases/latest).
2.  Find the asset that matches your OS and architecture (e.g., `mcpscan-windows-amd64.exe`).
3.  Download the file, and you're ready to run it from your terminal.

## Usage

The tool is controlled via command-line flags. The `--host` (or `-h`) flag is required.

### Basic Syntax

```bash
# On Linux/macOS
./mcpscan -h < host | hosts.txt > [flags]

# On Windows
.\mcpscan.exe -h < host | hosts.txt > [flags]
```

### Flags

| Short | Long         | Description                                                          | Default                  |
| :---- | :----------- | :------------------------------------------------------------------- | :----------------------- |
| `-h`  | `--host`     | **(Required)** A single domain or a path to a `.txt` file with hosts. | ` `                        |
| `-o`  | `--output`   | Path to save the results.                                            | `<host_input>_results.txt` |
| `-s`  | `--start-port` | Port to start scanning from.                                         | `25000`                  |
| `-e`  | `--end-port`   | Port to end scanning at.                                             | `30000`                  |
| `-w`  | `--workers`  | Number of concurrent scan threads.                                   | `150`                    |
| `-r`  | `--rate`     | Max scans per second.                                                | `200`                    |
| `-t`  | `--timeout`  | Connection timeout in seconds.                                       | `5`                      |
| `-d`  | `--debug`    | Enable detailed debug logging.                                       | `false`                  |
| `-?`  | `--help`     | Show the help message.                                               | `false`                  |

### Examples

**1. Scan a single server on the default port:**
```bash
./mcpscan -h hypixel.net
```

**2. Scan a range of ports on a single server:**
```bash
./mcpscan -h my-server.com -s 25500 -e 25600 -o found.txt
```

**3. Scan multiple hosts from a file with high performance:**
*(Assumes you have a `hosts.txt` file in the same directory)*
```bash
./mcpscan -h hosts.txt -s 20000 -e 30000 -w 250 -r 500
```

## 🛠️ Building From Source

If you prefer to build the project yourself:

1. Install [golang](https://go.dev/dl/)

2.  Clone the repository:
    ```bash
    git clone https://github.com/tmquan2508/mcpscan.git
    cd mcpscan
    ```
3.  Build the binary:
    ```bash
    go build -o mcpscan .
    ```