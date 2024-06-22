package runner

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"strings"
)

func assignArguments(str string, arguments map[string]string) string {
	for key, value := range arguments {
		str = strings.ReplaceAll(str, "{"+key+"}", value)
	}

	return str
}

func build(s *mySetting, arguments map[string]string) (string, bool) {
	buildLog("Building...")

	cmds := []string{}
	json.Unmarshal([]byte(s.buildCommand()), &cmds)

	for k, v := range cmds {
		cmds[k] = assignArguments(v, arguments)
	}

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
