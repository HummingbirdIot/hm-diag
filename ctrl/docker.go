package ctrl

import (
	"bufio"
	"log"
	"os/exec"

	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/util"
)

var dockerResetScript = "./docker_reset.sh"

func DockerReset() {
	fn := func() error {
		log.Println("to reset docker")
		log.Println("exec cmd: bash", dockerResetScript)
		cmd := exec.Command("bash", dockerResetScript)
		cmd.Dir = config.Config().GitRepoDir
		p, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}
		err = cmd.Start()
		if err != nil {
			return err
		}

		in := bufio.NewScanner(p)
		for in.Scan() {
			log.Println("reset docker cmd output: ", in.Text())
		}
		return nil
	}
	util.Sgo(fn, "reset docker error")
}
