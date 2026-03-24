package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/chzyer/readline"
	"slices"
	"syscall"
	"path/filepath"
)

type Shell struct {
	cwd string
	builtins map[string]func([]string)
	path string
	loop bool
	reader *readline.Instance
	history []string
	defaultHistoryLimit int
	name string
	stdin *os.File
}

// Convert a path argument to its absolute path
func (s* Shell) toAbs(path string) (string, error){
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}

	elems := strings.SplitN(path, "/", 2)
	if elems[0] == "~" {
		home := os.Getenv("HOME")
		return filepath.Join(home, strings.Join(elems[1:], "")), nil
	}

	return filepath.Join(s.cwd, path), nil
}

func (s Shell) executePathCommand(command string, args []string) {
	cmd := exec.Command(command, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Command failed with error: %v\n", err)
	}

	fmt.Printf("%s", string(output))
}

func (s* Shell) _updateHistory(command string, args []string) {
	entry := command + " " + strings.Join(args, " ")
    
    s.history = append(s.history, entry)
}

func parseRedirect(args []string) (file string, cleaned []string, err error) {
	pos := slices.IndexFunc(args, func(s string) bool {
		return s == ">" || s == "1>"
	})

	if pos == -1 {
		return "", args, nil
	}

	if pos+1 >= len(args) {
		return "", nil, fmt.Errorf("redirection requires a file argument")
	}

	return args[pos+1], args[:pos], nil
}

func (s *Shell) withRedirection(file string, fn func()) error {
	if file == "" {
		fn()
		return nil
	}

	absPath, err := s.toAbs(file)
    if err != nil {
        return err
    }

    fd, err := os.Create(absPath)
    if err != nil {
        return err
    }
    defer fd.Close()

    savedStdout, err := syscall.Dup(int(os.Stdout.Fd()))
    if err != nil {
        return fmt.Errorf("failed to save stdout: %w", err)
    }
    defer func() {
        syscall.Dup2(savedStdout, int(os.Stdout.Fd()))
        syscall.Close(savedStdout)
    }()

    if err := syscall.Dup2(int(fd.Fd()), int(os.Stdout.Fd())); err != nil {
        return err
    }

    fn()
    return nil
}

func (s *Shell) dispatch(command string, args []string) {
	if cmd, ok := s.builtins[command]; ok {
		cmd(args)
		return
	}

	if _, err := exec.LookPath(command); err == nil {
		s.executePathCommand(command, args)
		return
	} 
	fmt.Fprintf(os.Stderr, "%s: command not found", command)
}

func (s *Shell) Loop() {
	defer s.reader.Close()
	
	for s.loop {
		line, err := s.reader.Readline()
		if err != nil {
			break
		}

		tokens := strings.Fields(line)
        if len(tokens) == 0 {
            continue
        }
		
		command, args := tokens[0], tokens[1:]
		s._updateHistory(command, args)

		outFile, args, err := parseRedirect(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", s.name, err)
			continue
		}

		if err := s.withRedirection(outFile, func() {
			s.dispatch(command, args)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", s.name, err)
		}
	}
}

func NewShell() *Shell {
	path := os.Getenv("PATH")
	cwd, err := os.Getwd()
	if err != nil {
    	fmt.Println(err)
	}
	
	s := &Shell{
		cwd: cwd,
		path: path,
        loop:     true,
        builtins: make(map[string]func([]string)),
		defaultHistoryLimit: 16,
		name: "myshell",
    }

	config := &readline.Config{
		Prompt: "$ ",
		Listener: s, 
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		panic(err)
	}
    
	s.reader = rl
    // Register commands here
    s.builtins["exit"] = s._exit
    s.builtins["type"] = s._type
    s.builtins["echo"] = s._echo
	s.builtins["history"] = s._history
	s.builtins["pwd"] = s._pwd
	s.builtins["cd"] = s._cd
    
    return s
}

func main() {
	shell := NewShell()
	shell.Loop()
}
