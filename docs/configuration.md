# Detecctor configuration

## Configuration file

Before the Detecctor is run, it _must_ contain a configuration file. The attributes, required to successfully run the
server, are:

1. server host
2. authPassword
3. pluginDir
4. botToken
5. database username
6. database password

The configuration file formats supported are **YAML, JSON and TOML**.

Example `config` file in `yaml` format:

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

## Flags

To change the location of the configuration files, to enable persistence, both the configuration file location and
format can be specified with flags.

```bash
./main --help 
    --config-format # Format of the configuration files (yaml, json or toml)
    --config-file # Path of the configuration file (default: working directory)
```
