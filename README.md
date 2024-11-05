# Http Analyzer
Simple and easy tool to analyze http clients, get the http headers ordered, the client's ip and the ja3 fingerprint. All of this **without third-party dependencies or external libraries**.

## Prerequisites
Before you can build and run this project, ensure you have the following installed on your machine:
- Go (version 1.23.1+)
- openssl for generating auto-signed certificates
- golangci-lint for linting

## Installation
### Clone the repository:
```bash
git clone git@github.com:senonide/http-analyzer.git
cd http-analyzer
```

### Install dependencies:
```bash
go mod tidy
```

## Usage

> [!NOTE]
> Try it in less than 1 minute:

First of all, run the following commands:
```bash
make prepare
make run
```

We are ready to go. Send an HTTPS request:
```bash
curl "https://localhost:8443" --insecure | jq
```
The response will be a json containing all information:
```json
{
  "client_ip": "[::1]:33712",
  "ja3": "90e8176230f553294d993752a60f78bb",
  "http_method": "GET",
  "headers": [
    "Host: localhost:8443",
    "User-Agent: curl/8.9.1",
    "Accept: */*"
  ]
}
```
