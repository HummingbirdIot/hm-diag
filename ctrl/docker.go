package ctrl

import (
	"bufio"
	"os/exec"

	"github.com/kpango/glg"
	"xdt.com/hm-diag/config"
	"xdt.com/hm-diag/util"
)

var dockerResetScript = "./docker_reset.sh"

func DockerReset() {
	fn := func() error {
		glg.Info("to reset docker")
		glg.Debug("exec cmd: bash", dockerResetScript)
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
			glg.Info("reset docker cmd output: ", in.Text())
		}
		return nil
	}
	util.Sgo(fn, "reset docker error")
}
