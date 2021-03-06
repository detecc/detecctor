# Detecctor

Detecctor is a β‘fast, fully customizable π₯οΈ monitoring platform. It uses various π€ chatbots as a π² notification
service. It is designed for use with π plugins, which enable total control over the functionality of both server and
the client. All you do is issue a command and let the π plugin deal with the rest. You can include provided plugins,
write your own or include plugins from the community.

## π§ Configuration

Before running the server, check out the [configuration guide](/docs/configuration.md).

## π€ Supported bots

| Chat service | Supported     |
|    :----:   |    :----:     |
| Telegram       | βοΈ   |
| Slack        | Planned      |
| Discord       | Planned      |

## π Running the server

### Using π³ Docker or docker-compose

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