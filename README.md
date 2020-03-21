# ip_diff
[![CircleCI](https://circleci.com/gh/x-way/ip_diff/tree/master.svg?style=svg)](https://circleci.com/gh/x-way/ip_diff/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/x-way/ip_diff)](https://goreportcard.com/report/github.com/x-way/ip_diff)

Compare two lists of IP prefixes (added/removed subnets).

## Installation

```
# go get github.com/x-way/ip_diff
```

## Usage

```
# cat a.txt
192.168.0.0/16
10.0.0.0/9
2001:db8::/64
10.128.0.0/9

# cat b.txt
10.0.0.0/8
192.168.0.0/17
192.168.128.0/24
192.168.129.0/24
192.168.130.0/24
192.168.132.0/22
192.168.136.0/21
192.168.144.0/20
192.168.160.0/19
192.168.192.0/18
2001:db8::/48

# ip_diff a.txt b.txt
ip_diff a.txt b.txt
--- a.txt
+++ b.txt
+2001:db8:0:1::/64
+2001:db8:0:2::/63
+2001:db8:0:4::/62
+2001:db8:0:8::/61
+2001:db8:0:10::/60
+2001:db8:0:20::/59
+2001:db8:0:40::/58
+2001:db8:0:80::/57
+2001:db8:0:100::/56
+2001:db8:0:200::/55
+2001:db8:0:400::/54
+2001:db8:0:800::/53
+2001:db8:0:1000::/52
+2001:db8:0:2000::/51
+2001:db8:0:4000::/50
+2001:db8:0:8000::/49
-192.168.131.0/24
```
