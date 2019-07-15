package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func main() {
	var (
		tcp     bool
		addr    string
		qname   string
		qtype   uint16 = dns.TypeA
		timeout time.Duration
	)

	flag.BoolVar(&tcp, "tcp", false, "Use TCP instead of UDP.")
	flag.StringVar(&addr, "addr", "127.0.0.1:5300", "Server address.")
	flag.DurationVar(&timeout, "timeout", 3*time.Second, "Query timeout.")
	flag.Var((*fqdnVar)(&qname), "name", "Query name.")
	flag.Var((*dnsType)(&qtype), "type", "Query type.")

	flag.Parse()

	if qname == "" {
		fmt.Fprintln(os.Stderr, "No query name provided.")
		os.Exit(1)
	}

	// create a new DNS client:

	client := new(dns.Client)
	if tcp {
		client.Net = "tcp"
	}

	// create the query message:

	query := new(dns.Msg).SetQuestion(qname, qtype)

	// send the message and wait for the response:

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	resp, rtt, err := client.ExchangeContext(ctx, query, addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send query: %s\n", err)
		os.Exit(1)
	}

	// display the result:

	fmt.Printf("%s\n", resp.String())
	fmt.Printf("; server = %s\n", addr)
	fmt.Printf("; time = %s\n", rtt)
}

// fqdnVar is a DNS name and implements flag.Value interface.
type fqdnVar string

func (v *fqdnVar) Set(s string) error {
	*(*string)(v) = dns.Fqdn(strings.ToLower(s))
	return nil
}

func (v *fqdnVar) String() string {
	return *(*string)(v)
}

// dnsType is a DNS record type and implements flag.Value interface.
type dnsType uint16

func (v *dnsType) Set(s string) error {
	value, ok := dns.StringToType[s]
	if !ok {
		return fmt.Errorf("unknown record type")
	}

	*(*uint16)(v) = value
	return nil
}

func (v *dnsType) String() string {
	str, ok := dns.TypeToString[*(*uint16)(v)]
	if !ok {
		return fmt.Sprintf("TYPE%d", *v)
	}
	return str
}
