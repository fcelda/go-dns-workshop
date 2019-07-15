package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"
	"github.com/miekg/dns"
	"github.com/oklog/run"
)

func main() {
	var (
		addr string
	)

	flag.StringVar(&addr, "addr", ":5300", "Server interface to listen on.")
	flag.Parse()

	l := newLogger()

	// DNS server handler
	handler := MyServer{
		logger: l,
	}

	g := run.Group{}

	// network interfaces
	for _, net := range []string{"udp", "tcp"} {
		net := net

		srv := dns.Server{
			Addr:    addr,
			Net:     net,
			Handler: &handler,
		}
		g.Add(func() error {
			l.Log("msg", "starting server", "addr", addr, "net", net)
			err := srv.ListenAndServe()
			if err != nil {
				l.Log("msg", "failed to listen", "addr", addr, "net", net, "err", err)
			}
			return err
		}, func(error) {
			srv.Shutdown()
		})
	}

	// termination handler
	{
		term := make(chan os.Signal, 1)
		g.Add(func() error {
			signal.Notify(term, os.Interrupt)
			if _, ok := <-term; ok {
				l.Log("msg", "requested to stop")
			}
			return nil
		}, func(error) {
			signal.Stop(term)
			close(term)
		})
	}

	err := g.Run()
	if err != nil {
		l.Log("msg", "terminated with error", "err", err)
		os.Exit(1)
	}
}

func newLogger() log.Logger {
	w := log.NewSyncWriter(os.Stdout)
	l := log.NewLogfmtLogger(w)

	return log.With(
		l,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
}
