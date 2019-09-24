# What's my IP

A single binary to look up the machine's external IP address via a DNS query.  
There are multiple servers, if one fails, it tries the next.


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


## Example
```bash
$ ./whatsmyip
1.2.3.4
```
