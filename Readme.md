# Detecctor

Detecctor is a :zap: fast, fully customizable :heartpulse: monitoring platform. It uses **Telegram** as a notification
service :inbox_tray:. It is designed for use with plugins :electric_plug:, which enable total control over the
functionality of both server and the client. All you do is issue a command and let the :electric_plug: plugin deal with
the rest. You can include provided plugins, write your own or include plugins from the community.

## Configuration

Before running the server, check out the [configuration guide](/docs/configuration.md).

## Running the server

### Using Docker or docker-compose

The provided _docker-compose_ file will run both the MongoDB and the server.

```bash
docker-compose up -d
``` 

### Deploying on Kubernetes

```bash
todo
```

### Standalone

```bash
go build main.go 
./main # use --help to get the flags 
```

### Note:

The project is still under development. More features coming soon.