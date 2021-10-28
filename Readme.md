# Detecctor

Detecctor is a ⚡fast, fully customizable 🖥️ monitoring platform. It uses various 🤖 chatbots as a 📲 notification
service. It is designed for use with 🔌 plugins, which enable total control over the functionality of both server and
the client. All you do is issue a command and let the 🔌 plugin deal with the rest. You can include provided plugins,
write your own or include plugins from the community.

## 🔧 Configuration

Before running the server, check out the [configuration guide](/docs/configuration.md).

## 🤖 Supported bots

| Chat service | Supported     |
|    :----:   |    :----:     |
| Telegram       | ✔️   |
| Slack        | Planned      |
| Discord       | Planned      |

## 🏃 Running the server

### Using 🐳 Docker or docker-compose

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