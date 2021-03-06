package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// https://code.blogs.iiidefix.net/posts/get-public-ip-using-dns/
var lookupMethods = []lookupMethod{
	{"ns1.google.com", "o-o.myaddr.l.google.com", txt},
	{"resolver1.opendns.com", "myip.opendns.com", aRecord},
	{"ns1-1.akamaitech.net", "whoami.akamai.net", aRecord},
}

var version string
var timeout time.Duration

const usage = `WhatsmyIP is a small IP address lookup utility.
It works by sending a DNS lookup. There are multiple
services set, if one fails, it tries the next.

Usage: %s [options]

Options:
`

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, usage, os.Args[0])
		flag.PrintDefaults()
	}

	ip4 := flag.Bool("4", false, "Lookup IPv4 address (default)")
	ip6 := flag.Bool("6", false, "Lookup IPv6 address")
	timeoutSeconds := flag.Int("t", 3, "Timeout for DNS lookup in seconds")
	flag.Parse()

	if *ip4 && *ip6 {
		fmt.Fprintln(os.Stderr, "Only one of -4 or -6 can be specified!")
		os.Exit(1)
	}

	if *ip4 || (!*ip4 && !*ip6) {
		version = "4"
	} else {
		version = "6"
	}

	timeout = time.Duration(*timeoutSeconds) * time.Second
}

type lookupType int

const (
	txt = iota
	aRecord
)

type lookupMethod struct {
	nameServer string
	lookupAddr string
	lookupType lookupType
}

func newResolver(lm lookupMethod, version string) *net.Resolver {
	return &net.Resolver{
		// https://github.com/golang/go/issues/19268
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp"+version, lm.nameServer+":53")
		},
	}
}

func lookupRecord(lm lookupMethod, resolver *net.Resolver, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var res string
	var err error
	switch lm.lookupType {
	case txt:
		var record []string
		record, err = resolver.LookupTXT(ctx, lm.lookupAddr)
		if err == nil {
			res = record[0]
		}
	case aRecord:
		var addrs []net.IPAddr
		addrs, err = resolver.LookupIPAddr(ctx, lm.lookupAddr)
		if err == nil {
			res = addrs[0].String()
		}
	}

	return res, err
}

func main() {
	success := false
	for _, lm := range lookupMethods {
		resolver := newResolver(lm, version)
		res, err := lookupRecord(lm, resolver, timeout)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		fmt.Println(res)
		success = true
		break
	}

	if !success {
		os.Exit(1)
	}
}
