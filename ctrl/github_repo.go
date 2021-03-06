package ctrl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/kpango/glg"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
)

const checkGitRepoIsUpdateCmd = config.MAIN_SCRIPT + " toUpdate"

func GitRepoReset() error {
	newGitDir := path.Join(os.TempDir(), "hnt-"+uuid.NewString())
	defer func() {
		if _, err := os.Stat(newGitDir); err == nil {
			glg.Info("clean git repo:", newGitDir)
			os.RemoveAll(newGitDir)
		}
	}()

	err := hntIotClone(newGitDir)
	if err != nil {
		return err
	}

	conf := config.Config()
	gitDir := conf.GitRepoDir
	if _, err = os.Stat(gitDir); err == nil {
		err = os.RemoveAll(gitDir)
		if err != nil {
			return errors.WithStack(errors.WithMessage(err, "delete old git repo error"))
		}
	} else {
		glg.Warn("git repo dir not exist:", gitDir)
	}

	err = os.Rename(newGitDir, gitDir)
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "rename new git repo error"))
	}

	err = copyGitRepoProxy()
	if err != nil {
		glg.Error(err)
	}
	err = copyGitReleaseProxy()
	if err != nil {
		glg.Error(err)
	}

	return nil
}

func getRepoBranch() (string, error) {
	cmdArr := strings.Split("git rev-parse --abbrev-ref HEAD", " ")
	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	cmd.Dir = config.Config().GitRepoDir
	buf, err := cmd.Output()
	if err != nil {
		return "", err
	}
	res := strings.TrimSpace(string(buf))
	return res, nil
}

func hntIotClone(dir string) error {
	gitRepoUrl := config.Config().GitRepoUrl
	proxyUrl, err := RepoMirrorUrl(config.Config().GitRepoDir)
	if err == nil && proxyUrl != "" {
		glg.Infof("use proxy %s to clone git repo %s\n", proxyUrl, gitRepoUrl)
		gitRepoUrl = strings.ReplaceAll(gitRepoUrl, config.GITHUB_URL, proxyUrl)
	}
	branch, err := getRepoBranch()
	if err != nil {
		return err
	}
	cmdStr := fmt.Sprintf(" git clone -b %s --depth=1 %s %s",
		branch, gitRepoUrl, dir)
	glg.Debug("exec cmd:", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	p, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "start git repo clone cmd pipe error"))
	}
	r := bufio.NewReader(p)

	err = cmd.Start()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "start git repo clone cmd error"))
	}
	for err == nil {
		ln, _, errIn := r.ReadLine()
		err = errIn
		if err == nil {
			s := string(ln)
			glg.Info("git clone cmd output:", s)
		} else if err == io.EOF {
			break
		} else {
			glg.Error("read git repo clone ouput error:", err.Error())
		}
	}

	err = cmd.Wait()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "git clone exit error"))
	}

	setOrigin := "git remote set-url origin " + config.Config().GitRepoUrl
	cmd = exec.Command("bash", "-c", setOrigin)
	cmd.Dir = dir
	buf, err := cmd.Output()
	if err != nil {
		glg.Error("git remote set-url exit error:", err.Error())
		return err
	} else {
		glg.Info("git remote set-url output:", string(buf))
	}

	return nil
}

func IsGitRepoToUpdate() (bool, error) {
	result := false
	resPrefix := ">>>state:"
	glg.Debug("exec cmd: bash -c ", checkGitRepoIsUpdateCmd)
	cmd := exec.Command("bash", "-c", checkGitRepoIsUpdateCmd)
	cmd.Dir = config.Config().GitRepoDir
	p, err := cmd.StdoutPipe()
	if err != nil {
		return false, errors.WithStack(errors.WithMessage(err, "start check git update cmd pipe error"))
	}
	r := bufio.NewReader(p)

	err = cmd.Start()
	if err != nil {
		return false, errors.WithStack(errors.WithMessage(err, "start check git update cmd error"))
	}
	for err == nil {
		ln, _, errIn := r.ReadLine()
		err = errIn
		if err == nil {
			s := string(ln)
			glg.Info("check git update cmd output:", s)
			if strings.HasPrefix(s, resPrefix) {
				res := strings.TrimPrefix(s, resPrefix)
				glg.Info("check git update cmd result:", res)
				if res != "yes" && res != "no" {
					return false, fmt.Errorf("check git update cmd, result is invalid: %s", s)
				}
				result = res == "yes"
			}
		} else if err == io.EOF {
			break
		} else {
			glg.Error("read check git update cmd ouput error:", err.Error())
		}
	}

	err = cmd.Wait()
	if err != nil {
		return false, errors.WithStack(errors.WithMessage(err, "check git update cmd exit error"))
	}
	return result, nil
}

func ExecMainUpdate() {
	go func() {
		glg.Debug("exec cmd: bash ", config.MAIN_SCRIPT)
		cmd := exec.Command("bash", config.MAIN_SCRIPT)
		cmd.Dir = config.Config().GitRepoDir
		err := cmd.Run()
		if err != nil {
			glg.Error("run main script error:", err)
		} else {
			glg.Info("run update cmd exit success")
		}
	}()
}
