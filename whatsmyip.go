package main

import (
	"context"
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

// Timeout for DNS lookup
const TimeoutSeconds = 3

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

func newResolver(lm lookupMethod) *net.Resolver {
	return &net.Resolver{
		// https://github.com/golang/go/issues/19268
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", lm.nameServer+":53")
		},
	}
}

func lookupRecord(lm lookupMethod, resolver *net.Resolver) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutSeconds*time.Second)
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
		resolver := newResolver(lm)
		res, err := lookupRecord(lm, resolver)
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
