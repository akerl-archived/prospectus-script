package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/akerl/prospectus/v3/plugin"
)

type config struct {
	Command string            `json:"command"`
	Env     map[string]string `json:"env"`
}

func run(cfg config) error {
	if cfg.Command == "" {
		return fmt.Errorf("required arg not provided: command")
	}

	cmd := exec.Command(cfg.Command)

	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Environ(), fmt.Sprintf("%s=%s", k, v))
	}

	var stdoutBytes bytes.Buffer
	var stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderrBytes.String())
	}

	fmt.Printf(stdoutBytes.String())
	return nil
}

func main() {
	cfg := config{}
	err := plugin.ParseConfig(&cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s", err)
		os.Exit(1)
	}
	err = run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
