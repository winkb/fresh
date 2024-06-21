package main

import (
	"context"
	"encoding/json"

	"github.com/winkb/fresh/runner"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cmdArr := []string{"cat", "abc/a.txt"}
	cmdStr, _ := json.Marshal(cmdArr)

	var sm = map[string]string{
		"build_commands": string(cmdStr),
		"valid_ext":      ".txt",
	}

	var s = runner.NewMySetting(sm)

	runner.Start(ctx, s)
}
