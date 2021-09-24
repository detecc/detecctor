# Subscribe and unsubscribe feature

## Concept

Subscriptions on **Detecctor** are similar to the pub-sub architecture, well, sort of. The basic idea is to subscribe to
messages, sent by the client without prior request (scheduled monitoring events).

## Command structure

The first command is the `/sub` or `/subscribe` command. The subscribe command is structured as:

```text
/sub node=exampleNode1 command=/auth notifyInterval=1
```

The order of the arguments is insignificant, since they are **key-value** arguments. The `commands` key specifies the
command, which will be monitored on the specified node. You can monitor multiple commands by passing the value as a
comma-separated list. If the command key is missing, all the commands from the node will be monitored.

The `nodes` key specifies which node to listen to. You can specify multiple nodes by passing the value as comma-separated
list. If the node key is not specified, a command or list of commands will be listened to on all the nodes.

The `notifyInterval` key is used to periodically notify the user about the last known result of the command that was
sent to the server. The interval is in minutes. If the notifyInterval key is not specified, the server will
automatically notify the user about the command immediately after receiving it.
