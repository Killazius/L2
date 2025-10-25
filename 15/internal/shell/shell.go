package shell

import (
	"15/internal/builtin"
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

const PS1 = "$ "
const commandExit = "exit"

type Shell struct {
	running bool
	built   *builtin.Builtin
}

func New() *Shell {
	return &Shell{
		running: true,
		built:   builtin.New(),
	}
}
func (s *Shell) RunInteractive() {
	fmt.Printf("%s Shell. Type 'exit' to quit.\n", runtime.GOOS)

	signal.Ignore(syscall.SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		for range sigChan {
			fmt.Printf("\n%s", PS1)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for s.running {
		fmt.Print(PS1)
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\nExiting shell.")
				break
			}
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == commandExit {
			break
		}

		s.executeCommand(input)
	}
}
func (s *Shell) executeCommand(input string) {
	parts := strings.Fields(input)
	cmd := parts[0]
	args := parts[1:]

	if s.built.IsBuiltin(cmd) {
		output, err := s.built.Execute(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else if output != "" {
			fmt.Println(output)
		}
	} else {
		cmd := exec.Command(cmd, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
		}
	}
}

func (s *Shell) RunBash(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening bash file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		s.executeCommand(line)
	}
}
