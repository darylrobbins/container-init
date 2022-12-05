package main

import (
	"errors"
	"github.com/darylrobbins/container-init/internal/lang"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	Service ServiceConfig `hcl:"service,block"`
	Env     []EnvConfig   `hcl:"env,block"`
}

type EnvConfig struct {
	Name  string `hcl:"name,label"`
	Value string `hcl:"value"`
}

type ServiceConfig struct {
	Processes []ProcessConfig `hcl:"process,block"`
}

type ProcessConfig struct {
	Command   []string `hcl:"command"`
	Directory string   `hcl:"directory,optional"`
}

func main() {
	scope := lang.Scope{}
	ctx := scope.EvalContext()

	var config Config
	err := hclsimple.DecodeFile("config.hcl", ctx, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	log.Printf("Configuration is %#v", config)

	done := make(chan struct{})
	go run(done)
	<-done // wait for background goroutine to finish
}

func run(done chan<- struct{}) {
	cmd := exec.Command("env")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	done <- struct{}{} // signal the main goroutine
}
