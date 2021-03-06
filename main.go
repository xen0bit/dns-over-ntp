package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/beevik/ntp"
	"github.com/miekg/dns"
)

var records = map[string]string{
	"test.service.": "192.168.0.2",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		// case dns.TypeA:
		// 	log.Printf("Query for %s\n", q.Name)
		// 	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Println(time)
		// 	//ip := records[q.Name]
		// 	ip := "123.123.123.123"
		// 	if ip != "" {
		// 		rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
		// 		if err == nil {
		// 			m.Answer = append(m.Answer, rr)
		// 		}
		// 	}
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
			if err != nil {
				panic(err)
			}
			rr, err := dns.NewRR(fmt.Sprintf("%s TXT %s", q.Name, time))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("service.", handleDnsRequest)

	// start server
	port := 5354
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
