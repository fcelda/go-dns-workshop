package main

import (
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/miekg/dns"
)

// mockWriter records all messages written with WriteMsg.
type mockWriter struct {
	dns.ResponseWriter // embed to satisfy interface
	msgs               []*dns.Msg
}

func (w *mockWriter) WriteMsg(msg *dns.Msg) error {
	w.msgs = append(w.msgs, msg)
	return nil
}

func mockServer() *MyServer {
	l := log.NewNopLogger()
	return &MyServer{
		logger: l,
	}
}

func TestNoQuestion(t *testing.T) {
	srv := mockServer()
	writer := &mockWriter{}

	query := new(dns.Msg)
	srv.ServeDNS(writer, query)

	if len(writer.msgs) != 1 {
		t.Fatalf("expected 1 message")
	}

	if writer.msgs[0].Rcode != dns.RcodeFormatError {
		t.Fatalf("expected format error")
	}
}

func TestHello(t *testing.T) {
	srv := mockServer()
	writer := &mockWriter{}

	query := new(dns.Msg).SetQuestion("hello.test.", dns.TypeTXT)
	srv.ServeDNS(writer, query)

	if len(writer.msgs) != 1 {
		t.Fatalf("expected 1 message")
	}

	resp := writer.msgs[0]

	if resp.Rcode != dns.RcodeSuccess {
		t.Fatalf("expected no error")
	}

	if len(resp.Answer) != 1 {
		t.Fatalf("expected 1 record in answer section")
	}

	rr := resp.Answer[0]
	if rr.Header().Rrtype != dns.TypeTXT {
		t.Fatalf("expected TXT record in the answer")
	}

	txt := rr.(*dns.TXT)
	if len(txt.Txt) != 1 || !strings.HasPrefix(txt.Txt[0], "Hi!") {
		t.Fatalf("expected TXT record data containing Hi!")
	}
}
