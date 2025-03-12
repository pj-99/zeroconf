package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/miekg/dns"
)


type SumHandler struct{}


func parseTxtToInt(t string, prefix string, dest *int) error {
	valStr := strings.TrimPrefix(t, prefix)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return err
	}
	*dest = val
	return nil
}

func (h *SumHandler) CheckIsMatch(query *dns.Msg) bool {
	for _, rr := range query.Extra {
		if txt, ok := rr.(*dns.TXT); ok {
			for _, t := range txt.Txt {
				// Some logic to check if the query is matched
				if strings.HasPrefix(t, "needToResponse=true") {
					return true
				}
			}
		}
	}
	return false;
}

func(h *SumHandler) Handle(query *dns.Msg, resp *dns.Msg) error {
	var a, b int
	var hasA, hasB bool

	// Parse TXT (maybe can use regex)
	for _, rr := range query.Extra {
		if txt, ok := rr.(*dns.TXT); ok {
			for _, t := range txt.Txt {
				fmt.Println("T:", t)
				if strings.HasPrefix(t, "a=") {
					if err := parseTxtToInt(t, "a=", &a); err == nil {
						hasA = true
					}
				}
				if strings.HasPrefix(t, "b=") {
					if err := parseTxtToInt(t, "b=", &b); err == nil {
						hasB = true
					}
				}
			}
		}
	}

	if !hasA || !hasB {
		fmt.Println("No matching")
		return nil
	}

	// Put the response RR in the additional
	// Add a txt record with ans=a+b
	txt := &dns.TXT{
		Hdr: dns.RR_Header{
			// Assume only one question
			Name:   query.Question[0].Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    0,
		},
		Txt: []string{fmt.Sprintf("ans=%d", a+b)},
	}
	fmt.Println("Result", txt)

	resp.Extra = append(resp.Extra, txt)
	return nil
}

func main() {
	server, err := zeroconf.Register("GoZeroconf", "_workstation._tcp", "local.", 42424, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		panic(err)
	}

	handler := &SumHandler{}
	server.SetAdditionalHandler(handler)
	
	defer server.Shutdown()

	// Clean exit.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		// Exit by user
	case <-time.After(time.Second * 120):
		// Exit by timeout
	}

	log.Println("Shutting down.")
}