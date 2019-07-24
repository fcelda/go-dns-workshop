package main

import (
	"github.com/go-kit/kit/log"
	"github.com/miekg/dns"
)

// MyServer is a custom DNS server which implements dns.Handler interface.
type MyServer struct {
	logger log.Logger
}

func (srv *MyServer) ServeDNS(w dns.ResponseWriter, query *dns.Msg) {
	// Respond with "Format Error" in case of invalid number of questions.
	if len(query.Question) != 1 {
		srv.writeError(w, query, dns.RcodeFormatError)
		return
	}

	// Respond with "Not Implemented" error if QCLASS != IN.
	question := query.Question[0]
	if question.Qclass != dns.ClassINET {
		srv.writeError(w, query, dns.RcodeNotImplemented)
		return
	}

	srv.logger.Log("msg", "got query", "qname", question.Name, "qtype", question.Qtype)

	// Respond with "Refused" for other than "test." subdomains.
	if !dns.IsSubDomain("test.", question.Name) {
		srv.writeError(w, query, dns.RcodeRefused)
		return
	}

	// Magic happens here.
	if question.Name == "hello.test." && question.Qtype == dns.TypeTXT {
		txt := dns.TXT{}

		// Build record header based on the question.
		txt.Hdr.Name = question.Name
		txt.Hdr.Rrtype = question.Qtype
		txt.Hdr.Class = question.Qclass
		txt.Hdr.Ttl = 60

		// Record data.
		txt.Txt = []string{"Hi! I'm a DNS server in Go."}

		// Send response.
		resp := new(dns.Msg).SetReply(query)
		resp.Answer = make([]dns.RR, 1)
		resp.Answer[0] = &txt
		w.WriteMsg(resp)
		return
	}

	// Respond with "Server Failure" if we got this far.
	srv.writeError(w, query, dns.RcodeServerFailure)
}

// writeError constructs and sends error message with given RCODE.
func (srv *MyServer) writeError(w dns.ResponseWriter, query *dns.Msg, rcode int) {
	resp := new(dns.Msg).SetRcode(query, rcode)
	w.WriteMsg(resp)
}
