# Kerbeus - Kerberos AS-REQ & AS-REP Sniffer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Linux-green.svg)](https://www.kernel.org/)

Kerbeus is a specialized packet sniffer designed to capture and analyze Kerberos authentication traffic for security research and penetration testing. It captures AS-REQ and AS-REP messages to extract Pre-Authentication hashes that can be cracked offline with tools like John the Ripper.

<img width="1171" height="821" alt="image" src="https://github.com/user-attachments/assets/c3186a21-6276-462c-bf17-25063ee68678" />

## 🎯 Purpose

Kerbeus facilitates the security testing of Active Directory environments by:
- Capturing Kerberos Pre-Authentication hashes without user interaction
- Supporting multiple encryption types (AES128, AES256, RC4-HMAC)
- Automatically extracting salt values from AS-REP messages
- Removing duplicate hashes automatically
- Generating John the Ripper compatible hash formats

## ⚠️ Legal Disclaimer

**This tool is designed for authorized security audits and penetration testing only.** Unauthorized access to computer systems is illegal. Always obtain explicit written permission before testing any network or system.

## 🚀 Features

- ✅ **Dual Protocol Support**: Captures both TCP/UDP Kerberos traffic on port 88
- ✅ **Multiple Encryption Types**: AES128 (etype 17), AES256 (etype 18), RC4-HMAC (etype 23)
- ✅ **Automatic Salt Extraction**: Parses PA-ETYPE-INFO and PA-ETYPE-INFO2 structures
- ✅ **Deduplication**: Prevents saving duplicate hashes to output file
- ✅ **John the Ripper Format**: Hash output in `$krb5pa$` format for direct cracking
- ✅ **MITM ARP Spoofing**: Optional ettercap integration for network interception
- ✅ **Interactive Capture**: Real-time color-coded output with statistics
- ✅ **Automatic Correlation**: Matches AS-REQ and AS-REP packets by user
- ✅ **Linux Native**: Uses AF_PACKET for efficient packet capture

## 📋 Requirements

### System Requirements
- **OS**: Linux (Kali Linux recommended)
- **Kernel**: AF_PACKET support (standard in modern Linux)
- **Privileges**: Root/sudo (required for packet capture)
- **Go Version**: 1.21 or higher

### Dependencies
- `gopacket` - Packet capture and analysis
- `gokrb5` - Kerberos protocol library
- `golang.org/x/term` - Terminal handling
- `ettercap` (optional) - For MITM ARP spoofing

### Installation on Kali Linux
```bash
sudo apt-get update
sudo apt-get install -y golang-go ettercap
```

## 📦 Installation

### From Source

1. **Clone the repository**
   ```bash
   git clone https://github.com/bI8d0/kerbeus.git
   cd kerbeus
   ```

2. **Build the binary**
   ```bash
   go run build.go
   # or manually
   go build -o kerbeus main.go
   ```

3. **Run Kerbeus**
   ```bash
   sudo ./kerbeus
   ```

### Using the built binary
```bash
cd build/
sudo ./kerbeus -i eth0
```

## 🎮 Usage

### Basic Usage (Auto-detect interface)
```bash
sudo ./kerbeus
```

### Specify Network Interface
```bash
sudo ./kerbeus -i eth0
```

### Run in Background with Log
```bash
sudo ./kerbeus -i ens33 > kerbeus.log 2>&1 &
```

### Stop Capture
Press `q` + Enter or `Ctrl+C`

## 📖 How It Works

### Capture Scenarios

1. **Local Network**
   - Be on the same network segment as the Domain Controller
   - Passive capture of authentication traffic

2. **MITM - ARP Spoofing** (Recommended)
   ```bash
   sudo ./kerbeus -i eth0
   # Automatically starts ettercap for ARP spoofing
   ```

3. **Port Mirroring**
   - Configure SPAN/mirror on your network switch
   - Mirror traffic to the interface where Kerbeus runs

4. **Gateway Position**
   - Run Kerbeus on the network router/gateway
   - Captures all authentication traffic from the network

### Packet Flow

```
┌─────────────────────────────────────────────────────┐
│ 1. Client sends AS-REQ (Pre-Auth Request)           │
│    - Username, Realm, Encryption Type               │
│    - Encrypted Timestamp (PA-ENC-TIMESTAMP)         │
└──────────────────┬──────────────────────────────────┘
                   │
                   ↓
         ┌─────────────────────┐
         │ Kerbeus captures it │
         │ Stores in memory    │
         └──────────┬──────────┘
                    │
                    ↓
┌─────────────────────────────────────────────────────┐
│ 2. Server sends AS-REP (Pre-Auth Reply)             │
│    - Client name, Realm                             │
│    - PA-ETYPE-INFO with salt                        │
│    - Encrypted response                             │
└──────────────────┬──────────────────────────────────┘
                   │
                   ↓
         ┌─────────────────────┐
         │ Kerbeus captures it │
         │ Matches with AS-REQ │
         │ Extracts hash & salt│
         │ Deduplicates        │
         │ Saves to file       │
         └─────────────────────┘
```

## 📄 Output Format

### Generated Files
```
hash_YYYY-MM-DD_HH-MM-SS.txt
```

### Hash Format
```
username:$krb5pa$ETYPE$USERNAME$REALM$SALT$CIPHER
```

### Example
```
administrator:$krb5pa$18$administrator$SQL.LOCAL$SQL.LOCALadministrator$a1b2c3d4e5f6...
```

### Field Explanation
| Field | Description | Example |
|-------|-------------|---------|
| ETYPE | Encryption type | 17=AES128, 18=AES256, 23=RC4-HMAC |
| USERNAME | User name (no domain) | administrator |
| REALM | Domain/Realm name | SQL.LOCAL |
| SALT | Extracted salt value | SQL.LOCALadministrator |
| CIPHER | Encrypted hash (hex) | a1b2c3d4e5f6... |

## 🔐 Cracking Hashes with John the Ripper

### Basic Cracking (Auto-detect format) - RECOMMENDED
```bash
john hash.txt --wordlist=rockyou.txt
```

### Force Format Manually
```bash
john --format=krb5pa-sha1 hash.txt --wordlist=rockyou.txt
```

### Show Cracked Passwords
```bash
john --show --format=krb5pa-sha1 hash.txt
```

### With Mutation Rules
```bash
john hash.txt --wordlist=rockyou.txt --rules=best64
```

### Incremental Attack (No Wordlist)
```bash
john --incremental hash.txt
```

### Multiple Dictionaries
```bash
john hash.txt --wordlist=rockyou.txt,passwords.txt
```

### GPU Acceleration (if available)
```bash
john --format=krb5pa-sha1 --device=1 hash.txt --wordlist=rockyou.txt
```

## 🛠️ Troubleshooting

### Problem: "afpacket failed: operation not permitted"
**Solution**: Run with sudo/root privileges
```bash
sudo ./kerbeus
```

### Problem: No packets captured
**Verify interface is UP:**
```bash
ip link set eth0 up
```

**Check for Kerberos traffic:**
```bash
sudo tcpdump -i eth0 port 88 -c 5
```

### Problem: "AS-REP arrived BEFORE AS-REQ"
**Normal in networks with high latency.** The hash is automatically discarded and will be retried on the next login attempt. This is expected behavior.

### Problem: Captures hashes but doesn't crack
1. Verify hash format is correct
2. Try with larger wordlists
3. The user might have a very complex password
4. Consider using rules or incremental mode

### Problem: Ettercap fails to start
**Verify ettercap is installed:**
```bash
which ettercap
sudo apt-get install ettercap-graphical
```

## 🔧 Technical Details

### Protocol Information
- **Protocol**: Kerberos v5 (RFC 4120)
- **Port**: 88/TCP and 88/UDP
- **Capture Types**: AS-REQ (0x6A) + AS-REP (0x6B)
- **PA-DATA Types**: 
  - Type 2: PA-ENC-TIMESTAMP (Pre-Authentication)
  - Type 11: PA-ETYPE-INFO (Encryption Type Info)
  - Type 19: PA-ETYPE-INFO2 (Encryption Type Info v2)

### Processing Pipeline
1. Packet capture at Layer 2 (Ethernet)
2. Extraction of Layer 3 (IPv4/IPv6) and Layer 4 (TCP/UDP)
3. AS-REQ/AS-REP message parsing
4. Pre-Authentication data extraction
5. Salt value extraction from PA-ETYPE-INFO structures
6. Hash formatting for John the Ripper
7. Deduplication by username and salt
8. File output with append mode

### Deduplication Strategy
- Tracks unique hashes by `username:salt:cipher` combination
- Prevents duplicate writes to output file
- Efficient memory usage with map-based tracking
- Updates hash if newer encryption type is found

## 📊 Performance Characteristics

| Aspect | Details |
|--------|---------|
| Capture Rate | Limited by network interface (typically 1-10k packets/sec) |
| Memory Usage | ~50MB baseline + map for tracking hashes |
| CPU Usage | Minimal (mostly waiting for packets) |
| Hash Extraction | Real-time (milliseconds per packet) |
| File I/O | Append operations (efficient) |

## 🔗 Dependencies

```go
require (
	github.com/google/gopacket v1.1.19
	github.com/jcmturner/gokrb5/v8 v8.4.4
	golang.org/x/term v0.38.0
)
```

## 📝 Compilation

### On Linux
```bash
go run build.go
```

### Manual Build
```bash
go build -ldflags="-w -s" -o kerbeus main.go
```

### Build for Debian Package

#### Prerequisites:
```bash
sudo apt-get update
sudo apt-get install build-essential devscripts debhelper dh-golang golang-go
```

#### Build the .deb package:
```bash
git clone https://github.com/bI8d0/kerbeus.git
cd kerbeus
dpkg-buildpackage -us -uc
```

The `.deb` package will be created in the parent directory.

#### Install the package:
```bash
sudo dpkg -i kerbeus_1.2_amd64.deb
```

#### Verify installation:
```bash
sudo ./kerbeus -h
```

## 🐛 Known Limitations

- **Linux Only**: Requires AF_PACKET support (Linux specific)
- **Root Privileges**: Packet capture requires elevated privileges
- **Network Position**: Needs access to Kerberos traffic (same segment or MITM)
- **Hash Format**: Only supports krb5pa format (Pre-Authentication)
- **Single Machine**: No multi-host clustering support

## 🤝 Contributing

Contributions are welcome! Please feel free to submit:
- Bug reports
- Feature requests
- Pull requests
- Documentation improvements

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ✍️ Author

**bI8d0** - [bI8d0@protonmail.com](mailto:bI8d0@protonmail.com)

## 🔗 Links

- **GitHub**: https://github.com/bI8d0/kerbeus
- **Issues**: https://github.com/bI8d0/kerbeus/issues
- **Wiki**: https://github.com/bI8d0/kerbeus/wiki

## 📚 References

- [RFC 4120 - Kerberos V5](https://tools.ietf.org/html/rfc4120)
- [John the Ripper - Kerberos Cracking](https://www.openwall.com/john/)
- [Kali Linux Tools](https://www.kali.org/tools/)
- [Active Directory Security](https://adsecurity.org/)

## ⚡ Quick Start Cheatsheet

```bash
# Clone and build
git clone https://github.com/bI8d0/kerbeus.git && cd kerbeus && go run build.go

# Run with auto-detect
sudo ./build/kerbeus

# Run with specific interface
sudo ./build/kerbeus -i eth0

# Check for hashes
ls -lh hash_*.txt

# Crack with John
john hash_*.txt --wordlist=rockyou.txt

# Show results
john --show hash_*.txt
```

---

⚠️ **Remember**: This tool is for authorized security testing only. Unauthorized use is illegal.
