package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func verifyDomain(domain string) bool {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error looking up MX records for domain %s: %v", domain, err)
	}

	if err == nil && len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error looking up TXT records for domain %s: %v", domain, err)
	}

	for _, txt := range txtRecords {
		if strings.HasPrefix(txt, "v=spf1") {
			hasSPF = true
			spfRecord =
			break
		}
	}
	dmarcRecord, err = net.LookupTXT("_dmarc." + domain)

	if err != nil {
		log.Printf("Error looking up DMARC records for domain %s: %v", domain, err)
	}

	for _, txt := range dmarcRecord {
		if strings.HasPrefix(txt, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = txt
			break
		}
	}

	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("  MX Records: %v\n", hasMX)
	fmt.Printf("  SPF Record: %v\n", hasSPF)
}



func main() {
	scannner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter email addresses (one per line). Press Ctrl+D (Unix) or Ctrl+Z (Windows) to end input:")
	for scannner.Scan() {
		verifyDomain(scannner.Text())
	}
	if err := scannner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}