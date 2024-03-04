# hello

Super simple http server which responds with just 'hello' message.

Built for health checking or something like that.

## usage

```bash
# print this help message
$ hello -h
$ hello --help

# run http server on default port: 9999
$ hello

# run http server on port number: PORT_NUMBER
$ hello PORT_NUMBER
```

## systemd

Put following lines in `/etc/systemd/system/hello.service`:

```
[Unit]
Description=Hello
After=syslog.target
After=network.target

[Service]
Type=simple
WorkingDirectory=/dir/to/hello
ExecStart=/path/to/bin/hello 9999
Restart=always
RestartSec=5
DynamicUser=yes
ReadOnlyPaths=/
MemoryLimit=10M
NoExecPaths=/bin /sbin /usr/bin /usr/sbin /usr/local/bin /usr/local/sbin

[Install]
WantedBy=multi-user.target
```

enable with:

```bash
$ sudo systemctl enable hello.service
```

and start with:

```bash
$ sudo systemctl start hello.service
```

## license

MIT

