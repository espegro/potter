# potter
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
    	listen address (default "127.0.0.1")
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
ssh-keygen -f ./host.key -N ''
```
$ ssh-keygen -f ./potter.key -N '' -t ed25519
```

Start the server listening on port 2000, ip 192.168.1.100, emulating Ubuntu 22.04, log to potter.json:

```
$ ./potter -p 2000 -l 192.168.1.100 -s "SSH-2.0-OpenSSH_8.9p1 Ubuntu-3" > potter.json
```

