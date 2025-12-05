# Subject: Email Verifier (Assignment + Solution)

Assignment (Problem)
- Build a CLI that checks a domain’s email readiness:
  - MX: domain has at least one MX record
  - SPF: domain has a TXT record starting with v=spf1
  - DMARC: _dmarc.domain has a TXT record starting with v=DMARC1
- Input: one item per line from stdin (email address or domain)
- Output: print a summary per input
- Use Go’s net package (net.LookupMX, net.LookupTXT)

Solution (This Code)
- Reads lines from stdin and sends each line directly to a verifier
- Verifier performs:
  - MX lookup via net.LookupMX
  - SPF detection by scanning TXT records for v=spf1
  - DMARC TXT lookup at _dmarc.<domain>
- Prints a simple summary:
  - Domain: <input>
  - MX Records: true|false
  - SPF Record: true|false
- Logs DNS lookup errors but continues processing

Run
1) cd "Email Verifier"
2) go run main.go
3) Type inputs (Ctrl+D to end on Linux/macOS)

Try
- Domain input:
  - echo "example.com" | go run main.go
- Multiple:
  - printf "example.com\ngmail.com\n" | go run main.go
- Note: Input line is used as the domain as-is.

Files
- main.go