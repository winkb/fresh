package main

import (
	"context"
	"encoding/json"
	"fresh/runner"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cmdArr := []string{"ls", "."}
	cmdStr, _ := json.Marshal(cmdArr)

	runner.Settings["build_commands"] = string(cmdStr)
	runner.Settings["valid_ext"] = ".txt"

	runner.Start(ctx, "./abc")
}
