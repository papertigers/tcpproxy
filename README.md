# tcpproxy

This was inspired by @bahamas10 who wrote this blog article about using
[SmartOS as a Home
Router](http://www.daveeddy.com/2015/05/31/smartos-as-a-home-router/). In order
to avoid hairpin NAT he wrote a TCP proxy in node.js that forwards the
connection to the appropriate service. I wanted to do the same in golang mostly
for fun.

## Usage
```
go get -u github.com/papertigers/tcpproxy/cmd/tcpproxy
```

```
$ tcpproxy -h
Usage of tcpproxy:
  -l string
        local address (default ":8080")
  -r string
        remote address (default "localhost:80")
```

### Example

Proxy side

```
[link@typhoon ~]$ tcpproxy -r "192.168.1.10:22"
2017/08/30 20:02:10 New connection from: 10.0.1.28:54090
2017/08/30 20:02:24 Connection closed (10.0.1.28:54090 -> 192.168.1.10:22) sent: 3.5 kB received: 4.2 kB
```

Client side

```
nebula  link  ~  ssh root@10.0.1.30 -p 8080
Welcome to Ubuntu 16.04.3 LTS (GNU/Linux 4.3.0 x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/advantage
   __        .                   .
 _|  |_      | .-. .  . .-. :--. |-
|_    _|     ;|   ||  |(.-' |  | |
  |__|   `--'  `-' `;-| `-' '  ' `-'
                   /  ;  Instance (Ubuntu 16.04 20170403)
                   `-'   https://docs.joyent.com/images/container-native-linux

Last login: Wed Aug 30 20:02:01 2017 from 10.0.1.30
root@plex:~# uptime
 20:02:15 up 12 days, 19:20,  1 user,  load average: 0.00, 0.03, 0.02
root@plex:~# logout
Connection to 10.0.1.30 closed.
```
