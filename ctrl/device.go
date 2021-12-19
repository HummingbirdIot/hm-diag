// conrol device or miner
package ctrl

import (
	"log"
	"os/exec"
)

func RebootDevice() error {
	log.Println("to reboot device")
	go func() { exec.Command("reboot").Run() }()
	log.Println("sent reboot device cmd")
	return nil
}
