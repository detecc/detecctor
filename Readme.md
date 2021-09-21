# detecctor

## About

Sn-bot is a fast, fully customizable monitoring system for **Oxen** Service Nodes. It uses **Telegram** as a
notification service. It consists of a server, MongoDB and multiple clients, installed on the servers running the
service nodes. Both the server and the client support use of custom plugins for issuing commands to service nodes.

## Configuration

Before the server is run, it should contain a configuration file. The configuration file formats supported are
**YAML, JSON and TOML**. Configuration file location and format can be specified with flags.

Example `config.yaml`:

```yaml
server:
  host: 192.168.1.2
  port: 7777
  authPassword: "yourPassword"
  pluginDir: "/usr/detecctor/plugins"
  plugins:
    - "examplePlugin1"
    - "examplePlugin2"
telegram:
  botToken: "yourTelegramApiToken"
mongodb:
  host: localhost
  port: 27017
  database: "example"
  username: "admin"
  password: "admin"
```

## Running the server

### Using docker-compose

The _docker-compose_ file will run both the MongoDB and the server.

```bash
docker-compose up -d
``` 

### Deploying on Kubernetes

ToDo

### Standalone

```bash
go build main.go 
./main # use --help to get the flags 
```
