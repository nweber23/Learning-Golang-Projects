package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func verifyDomain(domain string) {
	// Allow passing an email; extract domain if needed
	if strings.Contains(domain, "@") {
		parts := strings.Split(domain, "@")
		if len(parts) == 2 {
			domain = parts[1]
		}
	}

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
			spfRecord = txt
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error looking up DMARC records for domain %s: %v", domain, err)
	}
	for _, txt := range dmarcRecords {
		if strings.HasPrefix(txt, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = txt
			break
		}
	}

	fmt.Printf("Domain: %s\n", domain)
	fmt.Printf("  MX Records: %v\n", hasMX)
	fmt.Printf("  SPF Record: %v\n", hasSPF)
	if hasSPF {
		fmt.Printf("    SPF: %s\n", spfRecord)
	}
	fmt.Printf("  DMARC Record: %v\n", hasDMARC)
	if hasDMARC {
		fmt.Printf("    DMARC: %s\n", dmarcRecord)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter email addresses or domains (one per line). Press Ctrl+D to end input:")
	for scanner.Scan() {
		verifyDomain(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}