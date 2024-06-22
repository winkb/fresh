package main

import (
	"context"

	"github.com/winkb/fresh/runner"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var sm = map[string]string{
		"build_commands": "",
		"valid_ext":      ".txt",
		"root":           "./abc",
	}

	var s = runner.NewMySetting(sm, func(st, arguments map[string]string) []string {
		return []string{
			"bash", "-c", "ls {filename} && echo '" + st["root"] + "'",
		}
	})

	runner.Start(ctx, s)
}
