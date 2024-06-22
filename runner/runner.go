package runner

import (
	"io"
	"os/exec"
)

func (l *Starter) run(s *mySetting) bool {
	runnerLog("Running...")

	cmd := exec.Command(s.buildPath())

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

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-l.stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
