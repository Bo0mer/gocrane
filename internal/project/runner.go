package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func NewRunner() *Runner {
	return &Runner{}
}

type Runner struct {
}

type logWriter struct {
	logger *log.Logger
}

func (w logWriter) Write(data []byte) (int, error) {
	w.logger.Print(string(data))
	return len(data), nil
}

func (r *Runner) Run(ctx context.Context, path string) (*Process, error) {
	output := logWriter{
		logger: log.New(log.Writer(), "program ", log.Ltime|log.Lmsgprefix),
	}
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdout = output
	cmd.Stderr = output
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start program: %w", err)
	}
	return &Process{
		process: cmd.Process,
	}, nil
}

type Process struct {
	process *os.Process
}

func (p *Process) Stop() error {
	if err := p.process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send sigterm signal to program: %w", err)
	}
	state, err := p.process.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for program to stop: %w", err)
	}
	if code := state.ExitCode(); code != 0 {
		return fmt.Errorf("program exited with non-zero exit code: %d", code)
	}
	return nil
}