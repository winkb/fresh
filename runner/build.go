package runner

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

func build() (string, bool) {
	buildLog("Building...")

	cmds := []string{}
	json.Unmarshal([]byte(buildCommand()), &cmds)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	io.Copy(os.Stdout, stdout)
	errBuf, _ := io.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
