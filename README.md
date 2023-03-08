# logtimer
[![Actions Status](https://github.com/Eun/logtimer/workflows/push/badge.svg)](https://github.com/Eun/logtimer/actions)
[![Coverage Status](https://coveralls.io/repos/github/Eun/logtimer/badge.svg?branch=master)](https://coveralls.io/github/Eun/logtimer?branch=master)
[![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/logtimer)
[![go-report](https://goreportcard.com/badge/github.com/Eun/logtimer)](https://goreportcard.com/report/github.com/Eun/logtimer)
---
Enhance your output with a timer

```
$ ping 8.8.8.8 | logtimer
[11:26:45] PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
[11:26:45] 64 bytes from 8.8.8.8: icmp_seq=1 ttl=123 time=16.7 ms
[11:26:46] 64 bytes from 8.8.8.8: icmp_seq=2 ttl=123 time=16.5 ms
[11:26:47] 64 bytes from 8.8.8.8: icmp_seq=3 ttl=123 time=16.1 ms
[11:26:48] 64 bytes from 8.8.8.8: icmp_seq=4 ttl=123 time=18.3 ms
```

# Custom Format
```
$ ping 8.8.8.8 | logtimer --format="[%a, %d %b %Y %02H:%02M:%02S %Z] "
[Thu, 7 Feb 2019 11:27:18 CET] PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
[Thu, 7 Feb 2019 11:27:18 CET] 64 bytes from 8.8.8.8: icmp_seq=1 ttl=123 time=16.8 ms
[Thu, 7 Feb 2019 11:27:19 CET] 64 bytes from 8.8.8.8: icmp_seq=2 ttl=123 time=17.4 ms
[Thu, 7 Feb 2019 11:27:20 CET] 64 bytes from 8.8.8.8: icmp_seq=3 ttl=123 time=15.7 ms
[Thu, 7 Feb 2019 11:27:21 CET] 64 bytes from 8.8.8.8: icmp_seq=4 ttl=123 time=16.0 ms
```

# Relative to the start
```
$ ping 8.8.8.8 | logtimer --relative
[00:00:00] 64 bytes from 8.8.8.8: icmp_seq=21 ttl=123 time=18.3 ms
[00:00:01] 64 bytes from 8.8.8.8: icmp_seq=21 ttl=123 time=18.3 ms
[00:00:02] 64 bytes from 8.8.8.8: icmp_seq=21 ttl=123 time=18.3 ms
...
[85:30:04] 64 bytes from 8.8.8.8: icmp_seq=22 ttl=123 time=18.5 ms
```

## Build History
[![Build history](https://buildstats.info/github/chart/Eun/logtimer?branch=master)](https://github.com/Eun/logtimer/actions)%  