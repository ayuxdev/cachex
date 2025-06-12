<h1 align="center">
  <img src="images/cachex-logo.png" alt="cachex" width="100px">
  <br>
</h1>
<h3 align="center">Tool to detect cache poisoning vulnerabilities in Web API Endpoints</h3>

<p align="center">
  <img src="https://img.shields.io/badge/cacheX-blueviolet?style=flat-square">
  <img src="https://img.shields.io/github/go-mod/go-version/ayuxdev/cachex?style=flat-square">
  <img src="https://img.shields.io/github/license/ayuxdev/cachex?style=flat-square">
</p>

![demo](images/cachex-demo.gif)

## Features

- High-speed multi-threaded scanning
- Accurate detection of response manipulation and cache persistence
- Single and multi-header scan modes
- YAML-based payload configuration
- JSON or pretty output logging
- Optional file-based result export
- Tentative vs confirmed vulnerability tagging

## Installation

```bash
go install github.com/ayuxdev/cachex/cmd/cachex@latest
````

Or build manually:

```bash
git clone https://github.com/ayuxdev/cachex
cd cachex
make build
./cachex -h
```

## Usage

### Scan a single URL

```bash
cachex -u https://example.com
```

### Scan multiple targets

```bash
cachex -l urls.txt
```

### Scan URLs piped from another command

```bash
echo "https://example.com" | cachex
```

or

```bash
cat urls.txt | cachex
```

### All CLI Flags

| Category          | Flag              | Description                              |
| ----------------- | ----------------- | ---------------------------------------- |
| Input             | `-u, --url`       | URL to scan                              |
|                   | `-l, --list`      | File with list of URLs                   |
| Concurrency       | `-t, --threads`   | Number of threads to use                 |
|                   | `-m, --scan-mode` | `single` or `multi` header scan mode     |
| HTTP Client       | `--timeout`       | Total request timeout                    |
|                   | `--proxy`         | Proxy URL to use                         |
| Persistence Check | `--no-chk-prst`   | Disable persistence check                |
|                   | `--prst-requests` | Number of poisoning requests to send     |
|                   | `--prst-threads`  | Threads to use for persistence poisoning |
| Output            | `-o, --output`    | Output file (stdout if not specified)    |
|                   | `-j, --json`      | Enable JSON output                       |
| Payloads          | `--pcf`           | Path to custom payload config YAML       |

## Example

```bash
cachex -l targets.txt -t 50 --pcf payloads.yaml --json -o results.json
```

## Configuration

By default, these files are created:

* `~/.config/cachex/config.yaml`
* `~/.config/cachex/payloads.yaml`

You can configure:

* Custom payload headers
* Request headers
* Logging formats
* Concurrency and timeout
* Proxy and output preferences

## Output Formats

### Pretty Format

```
[vuln] [https://target.com] [Location Poisoning] [header: X-Forwarded-Host: evil.com] [poc: https://target.com?cache=XYZ]
```

### JSON Format

```json
{
  "URL": "https://target.com/",
  "IsVulnerable": true,
  "IsResponseManipulable": true,
  "ManipulationType": "ChangedBody",
  "RequestHeaders": {
    "Accept": "*/*",
    "User-Agent": "Mozilla/5.0"
  },
  "PayloadHeaders": {
    "X-Forwarded-Host": "evil.com"
  },
  "OriginalResponse": {
    "StatusCode": 200,
    "Headers": {
      "...": "..."
    },
    "Body": "...",
    "Location": ""
  },
  "ModifiedResponse": {
    "StatusCode": 200,
    "Headers": {
      "...": "..."
    },
    "Body": "...",
    "Location": ""
  },
  "PersistenceCheckResult": {
    "IsPersistent": true,
    "PoCLink": "https://target.example.com/?cache=XYZ",
    "FinalResponse": {
      "StatusCode": 200,
      "Headers": {
        "...": "..."
      },
      "Body": "...",
      "Location": ""
    }
  }
}
```

## Scan Modes

* `single`: scans each payload header independently (more precise)
* `multi`: scans all payload headers together (faster, less precise)

## Payload Headers

Defined in `~/.config/cachex/payloads.yaml`. Includes:

```yaml
payload_headers:
    Forwarded: for=127.0.0.1;host=evil.com;proto=https
    X-Client-IP: 127.0.0.1
    X-Custom-IP-Authorization: 127.0.0.1
    X-Forwarded-For: 127.0.0.1
    X-Forwarded-Host: evil.com
    X-Forwarded-Host-Override: evil.com
    X-Forwarded-Port: "443"
    X-Forwarded-Proto: https
    X-Forwarded-Scheme: https
    X-Forwarded-Server: evil.com
    X-HTTP-Method-Override: POST
    X-Host: evil.com
    X-Host-Override: evil.com
    X-Original-Host: evil.com
    X-Original-URL: /evilpath
    X-ProxyUser-Ip: 127.0.0.1
    X-Requested-With: XMLHttpRequest
    X-Rewrite-URL: /evilpath
    X-Url-Scheme: https
```

## How It Works

1. Sends a baseline request
2. Injects headers and observes differences
3. Confirms cache persistence via repeat requests
4. Logs the vulnerability with optional PoC link

## Contribute

Sure, PR's are welcome!

## License

MIT © [@ayuxdev](https://github.com/ayuxdev)
