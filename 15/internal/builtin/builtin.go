package builtin

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Command func(args []string) (string, error)
type Builtin struct {
	cmds map[string]Command
}

func New() *Builtin {
	b := &Builtin{
		cmds: make(map[string]Command),
	}
	b.register()
	return b
}

func (b *Builtin) register() {
	b.cmds["cd"] = b.cd
	b.cmds["pwd"] = b.pwd
	b.cmds["echo"] = b.echo
	b.cmds["kill"] = b.kill
	b.cmds["ps"] = b.ps
}
func (b *Builtin) IsBuiltin(cmd string) bool {
	_, exists := b.cmds[cmd]
	return exists
}

func (b *Builtin) Execute(cmd string, args []string) (string, error) {
	if fn, exists := b.cmds[cmd]; exists {
		return fn(args)
	}
	return "", fmt.Errorf("command not found: %s", cmd)
}

func (b *Builtin) cd(args []string) (string, error) {
	var dir string
	switch len(args) {
	case 0:
		dir = os.Getenv("HOME")
	case 1:
		dir = args[0]
	default:
		return "", fmt.Errorf("cd: too many arguments")
	}
	if err := os.Chdir(dir); err != nil {
		return "", fmt.Errorf("cd: %v", err)
	}
	return "", nil
}

func (b *Builtin) pwd(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("pwd: too many arguments")
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("pwd: %v", err)
	}
	return dir, nil
}
func (b *Builtin) echo(args []string) (string, error) {
	return strings.Join(args, " "), nil
}
func (b *Builtin) kill(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("kill: usage: kill <pid>")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("kill: invalid pid: %v", err)
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return "", fmt.Errorf("kill: %v", err)
	}
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return "", fmt.Errorf("kill: %v", err)
	}
	return "", nil
}
func (b *Builtin) ps(args []string) (string, error) {
	cmd := exec.Command("ps", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return "", cmd.Run()
}
