package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type Shell struct {
	cwd string
	builtins map[string]func([]string)
	path string
	loop bool
	reader *readline.Instance
	history []string
	defaultHistoryLimit int
}

func (s *Shell) _exit(args []string) {
	s.loop = false
}

func (s *Shell) _history(args []string) {
	limit := s.defaultHistoryLimit
	if len(args) > 0 {
		if val, err := strconv.Atoi(args[0]); err == nil {
        	limit = val
    	}
	}
	var history strings.Builder
	for i := max(0, len(s.history)-limit); i < len(s.history); i++ {
		fmt.Fprintf(&history, "%d %s\n", i, s.history[i])
	}
	fmt.Printf("%s", history.String())
}

func (s Shell) _type(args []string) {
	if len(args) == 0 {
		return
	}
	
	name := args[0]
	// Search builtins
	if _, ok := s.builtins[name]; ok {
		fmt.Printf("%s is a shell builtin\n", name)
		return
	}
	
	// Search PATH
	path, err := exec.LookPath(name)
	if err == nil {
		fmt.Printf("%s is %s\n", name, path)
		return
	}

	// Not Found
	fmt.Printf("%s: not found\n", name)		
}

func (s Shell) _echo(args []string) {
	if len(args) == 0 {
		return
	}

	fmt.Println(strings.Join(args, " "))
}

func (s* Shell) _pwd(args []string) {
	fmt.Println(s.cwd)
}

func (s* Shell) abs(path string) (string, error){
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

func (s* Shell) _cd(args []string) {
	if len(args) != 1 {
		return
	}

	input := args[0]
	absPath, err := s.abs(input)
	if err != nil {
		fmt.Printf("cd: %s: %v\n", input, err)
        return
    }
	
	info, err := os.Stat(absPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
 			fmt.Printf("cd: %s: No such file or directory\n", absPath)
		} else {
			fmt.Printf("cd: %s: %v\n", absPath, err)
		}
		return
	}

	if !info.IsDir() {
		fmt.Printf("cd: %s: Not a directory\n", absPath)
        return
    }

	s.cwd = absPath
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

func (s *Shell) Loop() {
	defer s.reader.Close()
	
	for s.loop {
		line, err := s.reader.Readline()
		if err != nil { // io.EOF
			break
		}

		tokens := strings.Fields(line)
        if len(tokens) == 0 {
            continue
        }
		
		command := tokens[0]
		args := tokens[1:]
		
		// Add to history
		s._updateHistory(command, args)

		// Try builtin
		if cmd, ok := s.builtins[command]; ok {
			cmd(args)
			continue
		} 

		// Search PATH
		_, err = exec.LookPath(command)
		if err == nil {
			s.executePathCommand(command, args)
			continue
		} 
		
		// Command not found
		fmt.Println(command + ": command not found")
		
	}
}

func NewShell() *Shell {
	path := os.Getenv("PATH")
	cwd, err := os.Getwd()
	if err != nil {
    	fmt.Println(err)
	}
	
	rl, err := readline.New("$ ")
	if err != nil {
		panic(err)
	}
    
	s := &Shell{
		cwd: cwd,
		path: path,
        loop:     true,
        builtins: make(map[string]func([]string)),
		reader: rl,
		defaultHistoryLimit: 16,
    }
    
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
