package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/miekg/dns"
)

var (
	waitTime    = flag.Int("wait", 10, "Duration in [s] to run discovery.")
	serviceName = flag.String("service", "_workstation._tcp.local.", "Service name")
)

// // A func to check a message related to interested service name will be nice
// func isRelated(msg *dns.Msg) bool {
// 	return true
// }

func main() {
	// Create a new DNS message
	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = false

	m.Question = []dns.Question{
		{
			Name:   *serviceName,
			Qtype:  dns.TypeSRV,
			Qclass: dns.ClassINET,
		},
	}

	// Add custom TXT record to the Additional section
	txtRR := &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   *serviceName,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    255,
		},
		Txt: []string{"a=3", "b=5"},
	}
	m.Extra = append(m.Extra, txtRR)

	// Query
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*waitTime))
	defer cancel()

	// Listen response
	msgCh := make(chan *dns.Msg)
	go func(results <-chan *dns.Msg) {
		for msg := range results {
			// Filter non-related
			if !msg.MsgHdr.Response {
				continue
			}

			log.Println("ANSWER:")
			for _, ans := range msg.Answer {
				if ans.Header().Name != *serviceName {
					continue
				}
				log.Println(ans)
			}
			log.Print("EXTRA:")
			for _, extra := range msg.Extra {
				if extra.Header().Name != *serviceName {
					continue
				}
				log.Println(extra)
			}
		}
		log.Println("No more entries.")
	}(msgCh)

	err = resolver.Query(ctx, m, msgCh)

	if err != nil {
		log.Fatalln("Failed to query:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}
