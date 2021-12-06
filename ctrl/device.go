// conrol device or miner
package ctrl

import (
	"log"
	"syscall"
)

func RebootDevice() error {
	log.Println("to reboot device")
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		log.Println("[error] reboot device error:", err)
		return err
	}
	log.Println("sent reboot device cmd")
	return nil
}
