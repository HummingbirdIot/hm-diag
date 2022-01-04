package ctrl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"xdt.com/hm-diag/config"
)

func GitRepoReset() error {
	newGitDir := path.Join(os.TempDir(), "hnt-"+uuid.NewString())
	defer func() {
		if _, err := os.Stat(newGitDir); err == nil {
			log.Println("clean git repo:", newGitDir)
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
		log.Println("git repo dir not exist:", gitDir)
	}

	err = os.Rename(newGitDir, gitDir)
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "rename new git repo error"))
	}

	err = copyGitRepoProxy()
	if err != nil {
		log.Println(err)
	}
	err = copyGitReleaseProxy()
	if err != nil {
		log.Println(err)
	}

	return nil
}

func hntIotClone(dir string) error {
	gitRepoUrl := config.Config().GitRepoUrl
	proxyUrl, err := RepoMirrorUrl(config.Config().GitRepoDir)
	if err == nil && proxyUrl != "" {
		log.Printf("use proxy %s to clone git repo %s\n", proxyUrl, gitRepoUrl)
		gitRepoUrl = strings.ReplaceAll(gitRepoUrl, config.GITHUB_URL, proxyUrl)
	}
	cmdStr := fmt.Sprintf(" git clone -b main --depth=1 %s %s", gitRepoUrl, dir)
	log.Println("exec cmd:", cmdStr)
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
			log.Println("git clone cmd output:", s)
		} else if err == io.EOF {
			break
		} else {
			log.Println("read git repo clone ouput error:", err.Error())
		}
	}

	err = cmd.Wait()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "git clone exit error"))
	}

	rmOriginCmd := "git remote remove origin"
	addOriginCmd := "git remote add origin " + config.Config().GitRepoUrl
	exec.Command("bash", "-c", rmOriginCmd).Run()
	exec.Command("bash", "-c", addOriginCmd).Run()

	return nil
}
