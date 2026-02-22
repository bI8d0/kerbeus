package stream

import (
	"bytes"

	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"

	"kerbeus/src/asrep"
	"kerbeus/src/asreq"
	"kerbeus/src/models"
)

type Factory struct {
	Filename    string
	ExitCh      chan struct{}
	Seen        map[string]bool
	PendingReqs map[string]*models.PendingASREQ
}

type Stream struct {
	net, transport gopacket.Flow
	filename       string
	buffer         bytes.Buffer
	exitCh         chan struct{}
	seen           map[string]bool
	pendingReqs    map[string]*models.PendingASREQ
	srcIP          string
	dstIP          string
}

func NewFactory(filename string, exitCh chan struct{}) *Factory {
	return &Factory{
		Filename:    filename,
		ExitCh:      exitCh,
		Seen:        make(map[string]bool),
		PendingReqs: make(map[string]*models.PendingASREQ),
	}
}

func (f *Factory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	return &Stream{
		net:         net,
		transport:   transport,
		filename:    f.Filename,
		exitCh:      f.ExitCh,
		seen:        f.Seen,
		pendingReqs: f.PendingReqs,
		srcIP:       net.Src().String(),
		dstIP:       net.Dst().String(),
	}
}

func (s *Stream) Reassembled(reassemblies []tcpassembly.Reassembly) {
	for _, reassembly := range reassemblies {
		s.buffer.Write(reassembly.Bytes)
	}
	s.processBuffer()
}

func (s *Stream) ReassemblyComplete() {
	s.processBuffer()
}

func (s *Stream) processBuffer() {
	for {
		if s.buffer.Len() < 4 {
			return
		}

		lenBytes := s.buffer.Bytes()[:4]
		msgLen := int(lenBytes[0])<<24 | int(lenBytes[1])<<16 | int(lenBytes[2])<<8 | int(lenBytes[3])

		if msgLen < 50 || msgLen > 65000 {
			s.buffer.Next(4)
			continue
		}

		if s.buffer.Len() < 4+msgLen {
			return
		}

		s.buffer.Next(4)
		buffer := make([]byte, msgLen)
		s.buffer.Read(buffer)

		if len(buffer) < 2 {
			continue
		}

		tag := buffer[0]

		switch tag {
		case 0x6A:
			asreq.Process(buffer, s.filename, s.seen, s.pendingReqs, s.srcIP, s.dstIP)
		case 0x6B:
			asrep.Process(buffer, s.filename, s.seen, s.pendingReqs, s.srcIP, s.dstIP)
		}
	}
}
