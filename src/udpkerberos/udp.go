package udpkerberos

import (
	"kerbeus/src/asrep"
	"kerbeus/src/asreq"
	"kerbeus/src/stream"
)

func Process(payload []byte, factory *stream.Factory, srcIP, dstIP string) {
	if len(payload) < 10 {
		return
	}

	tag := payload[0]
	if tag != 0x6A && tag != 0x6B {
		return
	}

	switch tag {
	case 0x6A:
		asreq.Process(payload, factory.Filename, factory.Seen, factory.PendingReqs, srcIP, dstIP)
	case 0x6B:
		asrep.Process(payload, factory.Filename, factory.Seen, factory.PendingReqs, srcIP, dstIP)
	}
}
