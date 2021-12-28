package ctrl

import (
	"bufio"
	"log"
	"os/exec"

	"xdt.com/hm-diag/util"
)

// todo: 优化 mainWorkDir 参数配置方式
const dockerResetScript = mainWorkDir + "docker_reset.sh"

func DockerReset() {
	fn := func() error {
		log.Println("to reset docker")
		log.Println("exec cmd: bash", dockerResetScript)
		cmd := exec.Command("bash", dockerResetScript)
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
