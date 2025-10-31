package executor

import (
	"15/internal/builtin"
	"15/internal/parser"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type Executor struct {
	builtin *builtin.Builtin
}

func New(builtin *builtin.Builtin) *Executor {
	return &Executor{builtin: builtin}
}

func (e *Executor) Execute(commands []*parser.Command) error {
	if len(commands) == 0 {
		return nil
	}

	if len(commands) == 1 {
		return e.executeSingleCommand(commands[0])
	}
	return e.executePipeline(commands)
}

func (e *Executor) executeSingleCommand(cmd *parser.Command) error {
	if e.builtin.IsBuiltin(cmd.Name) {
		output, err := e.builtin.Execute(cmd.Name, cmd.Args)
		if err != nil {
			return err
		}
		if output != "" {
			fmt.Println(output)
		}
		return nil
	}

	return e.executeExternalCommand(cmd, os.Stdin, os.Stdout)
}

func (e *Executor) executePipeline(commands []*parser.Command) error {
	var cmds []*exec.Cmd
	var pipes []*io.PipeReader
	var writers []*io.PipeWriter

	for i, cmdDef := range commands {
		cmd := exec.Command(cmdDef.Name, cmdDef.Args...)

		if i > 0 {
			cmd.Stdin = pipes[i-1]
		} else {
			cmd.Stdin = os.Stdin
		}

		if i < len(commands)-1 {
			reader, writer := io.Pipe()
			cmd.Stdout = writer
			pipes = append(pipes, reader)
			writers = append(writers, writer)
		} else {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr
		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			for _, w := range writers {
				_ = w.Close()
			}
			return fmt.Errorf("cannot start command %s: %v", cmd.Path, err)
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr error

	for i, cmd := range cmds {
		wg.Go(func() {
			func(i int, cmd *exec.Cmd) {
				if err := cmd.Wait(); err != nil {
					mu.Lock()
					lastErr = err
					mu.Unlock()
				}
				if i < len(cmds)-1 {
					_ = writers[i].Close()
				}
			}(i, cmd)
		})
	}

	wg.Wait()
	return lastErr
}
func (e *Executor) executeExternalCommand(cmd *parser.Command, stdin io.Reader, stdout io.Writer) error {
	command := exec.Command(cmd.Name, cmd.Args...)
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = os.Stderr

	return command.Run()
}
