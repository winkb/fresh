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

	cmdArr := []string{}
	if s.getBuildCommand != nil {
		cmdArr = s.getBuildCommand(s.settings, arguments)
	} else {
		json.Unmarshal([]byte(s.buildCommand()), &cmdArr)
	}

	for k, v := range cmdArr {
		cmdArr[k] = assignArguments(v, arguments)
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)

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
