package shell

import (
	"15/internal/builtin"
	"15/internal/executor"
	"15/internal/parser"
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
const goodbye = "poka!"

type Shell struct {
	running  bool
	built    *builtin.Builtin
	executor *executor.Executor
}

func New() *Shell {
	built := builtin.New()
	return &Shell{
		running:  true,
		built:    built,
		executor: executor.New(built),
	}
}
func (s *Shell) RunInteractive() {
	fmt.Printf("%s Shell. Type 'exit' to quit.\n", runtime.GOOS)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	var currentCmd *exec.Cmd

	go func() {
		for range sigChan {
			if currentCmd != nil {
				currentCmd.Process.Signal(syscall.SIGINT)
				currentCmd = nil
			}
			fmt.Printf("\n%s", PS1)
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for s.running {
		fmt.Print(PS1)
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n" + goodbye)
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
		input = expandEnvVars(input)
		s.executeCommand(input)
		currentCmd = nil
	}
}
func (s *Shell) executeCommand(input string) {
	commands, err := parser.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return
	}
	if err := s.executor.Execute(commands); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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

func expandEnvVars(input string) string {
	return os.Expand(input, func(key string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return ""
	})
}
