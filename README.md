# potter

![potter](potter.jpeg)

A simple SSH-honeypot written in Go


```
potter ssh honeypot.

Usage:
  -enc string
    	Custom enc. algs. (default "aes256-ctr,aes128-ctr")
  -i string
    	custom tag in JSON (default "ssh-pott")
  -k string
    	hostkey filename (default "potter.key")
  -kex string
    	Custom kex algs. (default "curve25519-sha256@libssh.org")
  -l string
    	listen address (default "0.0.0.0")
  -m string
    	Custom deny message (default "Access denied.")
  -mac string
    	Custom mac. algs. (default "hmac-sha2-256-etm@openssh.com,hmac-sha2-512-etm@openssh.com")
  -p string
    	ssh port (default "2222")
  -s string
    	ssh versionstring (default "OpenSSH_4.5")


Example:
potter -p 9000 -s "OpenSSH_8.9" -k ./potter.key -m "Go away" -i mysensor



```

Make a local ssh host key:
```
$ ssh-keygen -f ./potter.key -N '' -t ed25519
```

Or use the makefile

```
$ make key
```

Start the server listening on port 2000, ip 192.168.1.100, emulating Ubuntu 22.04, log to potter.json:

```
$ ./potter -p 2000 -l 192.168.1.100 -s "SSH-2.0-OpenSSH_8.9p1 Ubuntu-3" > potter.json
```


The log is in JSONL format
```json
{"timestamp": "2022-07-07T17:21:16+02:00", "id": "ssh-pott", "user": "user", "clientip": "167.99.214.128", "srcport": "47030", "password": "user", "clientversion": "SSH-2.0-Go" }
{"timestamp": "2022-07-07T17:21:37+02:00", "id": "ssh-pott", "user": "admin", "clientip": "167.99.214.128", "srcport": "47852", "password": "admin", "clientversion": "SSH-2.0-Go" }
{"timestamp": "2022-07-07T17:21:59+02:00", "id": "ssh-pott", "user": "steam", "clientip": "167.99.214.128", "srcport": "48692", "password": "steam", "clientversion": "SSH-2.0-Go" }
{"timestamp": "2022-07-07T17:22:22+02:00", "id": "ssh-pott", "user": "postgres", "clientip": "167.99.214.128", "srcport": "49518", "password": "postgres", "clientversion": "SSH-2.0-Go" }
{"timestamp": "2022-07-07T17:22:45+02:00", "id": "ssh-pott", "user": "oracle", "clientip": "167.99.214.128", "srcport": "50340", "password": "oracle", "clientversion": "SSH-2.0-Go" }
{"timestamp": "2022-07-07T17:39:16+02:00", "id": "ssh-pott", "user": "root", "clientip": "118.120.228.182", "srcport": "40042", "password": "root", "clientversion": "SSH-2.0-libssh_0.9.5" }
```

Public keys will also be logged.
```json
{"timestamp": "2022-07-07T18:59:01+02:00", "id": "ssh-pott", "user": "espegro", "clientip": "192.168.1.144", "srcport": "36478", "publickey": "sk-ssh-ed25519@openssh.com AAAAXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=", "clientversion": "SSH-2.0-OpenSSH_8.9p1 Ubuntu-3" }
```

It is also posible to change the crypto alg. to match a spesific server. They have to be supported in Go ssh lib.

