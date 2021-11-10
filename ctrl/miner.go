package ctrl

import (
	"log"
	"os/exec"
)

func ResyncMiner() error {
	log.Println("to resync miner")
	cmd := exec.Command("bash", "/home/pi/hnt_iot/trim_miner.sh")
	data, err := cmd.Output()
	if err != nil {
		log.Println("[error] resync miner error:", string(data))
		return err
	}
	log.Println("resync miner output:", string(data))
	return nil
}
