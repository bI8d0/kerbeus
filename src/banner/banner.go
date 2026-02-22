package banner

import (
	"fmt"
	"os"
	"time"
)

func Print(dev string) {
	clearScreen()
	now := time.Now().Format(time.RFC3339)
	fmt.Println()
	fmt.Println("  ╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("  ║                                                               ║")
	fmt.Println("  ║   ██╗  ██╗███████╗██████╗ ██████╗ ███████╗██╗   ██╗███████╗   ║")
	fmt.Println("  ║   ██║ ██╔╝██╔════╝██╔══██╗██╔══██╗██╔════╝██║   ██║██╔════╝   ║")
	fmt.Println("  ║   █████╔╝ █████╗  ██████╔╝██████╔╝█████╗  ██║   ██║███████╗   ║")
	fmt.Println("  ║   ██╔═██╗ ██╔══╝  ██╔══██╗██╔══██╗██╔══╝  ██║   ██║╚════██║   ║")
	fmt.Println("  ║   ██║  ██╗███████╗██║  ██║██████╔╝███████╗╚██████╔╝███████║   ║")
	fmt.Println("  ║   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝   ║")
	fmt.Println("  ║                                                               ║")
	fmt.Println("  ║             Kerberos AS-REQ & AS-REP Sniffer v1.2             ║")
	fmt.Println("  ║                           By bI8d0                            ║")
	fmt.Println("  ╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  ▶ Interface: %s\n", dev)
	fmt.Printf("  ▶ Start:     %s\n", now)
	fmt.Printf("  ▶ Capture:   AS-REQ + AS-REP\n")
	fmt.Printf("  ▶ Exit: Press 'q' or Ctrl+C to quit\n")
	fmt.Println()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func ShowHelp() {
	fmt.Fprintf(os.Stderr, `
╔═══════════════════════════════════════════════════════════════════════════╗
║                   Kerberos AS-REQ & AS-REP Sniffer v1.2                   ║
║                                                                           ║
║           Kerberos packet capture tool for security analysis              ║
║          Extracts Pre-Authentication hashes without duplicates            ║
╚═══════════════════════════════════════════════════════════════════════════╝

USAGE:
  sudo ./kerbeus [options]

OPTIONS:
  -i string
      Network interface to monitor (e.g: eth0, wlan0, ens33)
      If empty, uses the first active non-loopback interface
      Example: -i eth0

  -h
      Show this help message

DESCRIPTION:
  Kerbeus captures and correlates AS-REQ and AS-REP packets from the
  Kerberos protocol (port 88/TCP and 88/UDP) to extract authentication
  hashes that can be cracked with John the Ripper.

  ✓ Supports ciphers: AES128, AES256, RC4-HMAC (etype 17, 18, 23)
  ✓ Automatically extracts salt from AS-REP (PA-ETYPE-INFO and PA-ETYPE-INFO2)
  ✓ Automatically removes duplicates (saves each hash only once)
  ✓ Format compatible with John the Ripper

REQUIREMENTS:
  • Root/sudo permissions (for packet capture)
  • Linux with AF_PACKET support (Kali Linux recommended)
  • Ettercap installed (optional, for MITM ARP spoofing)
  • Access to Kerberos traffic on the network (port 88)

  CAPTURE SCENARIOS:
    1. Local network: Be on the same segment as the Domain Controller
    2. MITM: Use arp-spoofing (ettercap, bettercap, arpspoof)
    3. Port mirroring: Configure SPAN/mirror on switch
    4. Gateway: Run on the network router/gateway

GENERATED FILES:
  hash_YYYY-MM-DD_HH-MM-SS.txt
    └─> Hashes in John the Ripper format (no duplicates)

EXAMPLES:

  # Automatic capture (detects interface)
  sudo ./kerbeus

  # Specify interface manually
  sudo ./kerbeus -i eth0

  # Run in background with log
  sudo ./kerbeus -i ens33 > kerbeus.log 2>&1 &

  # Stop capture
  Press 'q' + Enter

CAPTURED HASH FORMAT:

  AS-REQ Pre-Authentication:
    $krb5pa$ETYPE$USERNAME$DOMAIN$SALT$CIPHER

    Example AES256:
    $krb5pa$18$administrator$SQL.LOCAL$SQL.LOCALadministrator$abc123...

  Fields:
    • ETYPE:   17=AES128, 18=AES256, 23=RC4-HMAC
    • USERNAME: User name (without domain)
    • DOMAIN:   Realm/domain (e.g. CONTOSO.LOCAL)
    • SALT:     Automatically extracted from AS-REP
    • CIPHER:   Encrypted hash in hexadecimal

CRACKING HASHES WITH JOHN THE RIPPER:

  ┌─────────────────────────────────────────────────────────────────────────┐
  │ Basic cracking (auto-detects format) (RECOMMENDED)                      │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john hash.txt --wordlist=rockyou.txt                                  │
  ├─────────────────────────────────────────────────────────────────────────┤
  │ Force format manually                                                   │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john --format=krb5pa-sha1 hash.txt --wordlist=rockyou.txt             │
  ├─────────────────────────────────────────────────────────────────────────┤
  │ Show cracked passwords                                                  │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john --show --format=krb5pa-sha1 hash.txt                             │
  ├─────────────────────────────────────────────────────────────────────────┤
  │ With mutation rules                                                     │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john hash.txt --wordlist=rockyou.txt --rules=best64                   │
  ├─────────────────────────────────────────────────────────────────────────┤
  │ Incremental attack (no wordlist)                                        │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john --incremental hash.txt                                           │
  ├─────────────────────────────────────────────────────────────────────────┤
  │ Attack with multiple dictionaries                                       │
  ├─────────────────────────────────────────────────────────────────────────┤
  │   john hash.txt --wordlist=rockyou.txt,passwords.txt                    │
  └─────────────────────────────────────────────────────────────────────────┘

SECURITY TIPS:

  ⚠️  LEGAL WARNING:
      This tool is ONLY for authorized security audits
      Unauthorized use is illegal and may result in criminal penalties

  ✓ Only use in test environments or with explicit authorization
  ✓ Captured hashes contain sensitive information
  ✓ Delete hash.txt files after analysis
  ✓ Do not share hashes on insecure channels

TROUBLESHOOTING:

  Issue: "afpacket failed: operation not permitted"
  Solution: Run with sudo/root

  Issue: No packets captured
  Solution: Verify interface is UP (ip link set eth0 up)
            Verify there is Kerberos traffic (tcpdump -i eth0 port 88)

  Issue: "AS-REP arrived BEFORE AS-REQ"
  Solution: Normal in networks with latency. Hash is automatically
            discarded and will be retried on the next login

  Issue: Captures hashes but doesn't crack
  Solution: Verify hash format
            Try with larger dictionaries
            The user might have a very complex password

TECHNICAL INFORMATION:

  • Protocol: Kerberos v5 (RFC 4120)
  • Port: 88/TCP and 88/UDP
  • Capture: AS-REQ (0x6A) + AS-REP (0x6B)
  • PA-DATA: Type 2 (PA-ENC-TIMESTAMP), Type 11/19 (ETYPE-INFO)
  • Processing: AS-REQ → AS-REP correlation by user
  • Deduplication: Avoids saving the same hash twice

AUTHOR:
  Developed for Active Directory pentesting and audits

`)
}
