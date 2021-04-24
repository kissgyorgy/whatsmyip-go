# What's my IP

A single binary to look up the machine's external IP address via a DNS query.  
There are multiple servers, if one fails, it tries the next.

It supports querying both IPv4 and IPv6 addresses.

## How it works

It's basically a Go implementation of this:

```bash
$ dig o-o.myaddr.l.google.com txt @ns1.google.com +short || \
dig myip.opendns.com @resolver1.opendns.com +short || \
dig whoami.akamai.net. @ns1-1.akamaitech.net. +short
```

but you don't need `dig` or anything else.

The idea and the addresses are from:
https://code.blogs.iiidefix.net/posts/get-public-ip-using-dns/

## Installation

Just grab the latest binary from the Releases page  
OR you can install it with Go:

```
$ go get -u github.com/kissgyorgy/whatsmyip-go
```

## Command line interface

```
$ ./whatsmyip -h
WhatsmyIP is a small IP address lookup utility.
It works by sending a DNS lookup. There are multiple
services set, if one fails, it tries the next.

Usage: whatsmyip [options]

Options:
  -4    Lookup IPv4 address (default)
  -6    Lookup IPv6 address
  -t int
        Timeout for DNS lookup in seconds (default 3)
```

## Example

```bash
$ ./whatsmyip
1.2.3.4
```
