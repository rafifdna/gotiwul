# Gotiwul

A high-performance TLS proxy server written in Go.

## Features

- TLS/SSL support
- Load balancing
- Health checking
- Docker support
- YAML configuration
- Logging middleware

## Setup

1. Generate SSL certificates:
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/server.key -out certs/server.crt \
  -subj "/CN=yourdomain.com"