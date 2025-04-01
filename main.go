package main

import (
	"context"
	"fmt"
	"github.com/turmantis/terawatt/terraform"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()
	versionsWant, err := terraform.DesiredVersion(ctx)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed trying to determine the appropriate version", err)
		os.Exit(1)
	}
	bin, err := terraform.BinaryHostPath(ctx, versionsWant)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to find host binary.", err)
		os.Exit(1)
	}
	argv := os.Args[1:]
	cmd := exec.CommandContext(ctx, bin, argv...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
	}
}
