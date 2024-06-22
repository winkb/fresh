package main

import (
	"context"
	"encoding/json"

	"github.com/winkb/fresh/runner"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cmdArr := []string{"cat", "{filename}"}
	cmdStr, _ := json.Marshal(cmdArr)

	var sm = map[string]string{
		"build_commands": string(cmdStr),
		"valid_ext":      ".txt",
		"root":           "./abc",
	}

	var s = runner.NewMySetting(sm)

	runner.Start(ctx, s)
}
