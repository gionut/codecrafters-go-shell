package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/chzyer/readline"
	"slices"
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
		
		command, args := tokens[0], tokens[1:]
		
		// Add to history
		s._updateHistory(command, args)

		// Check for input redirection
		pos := slices.IndexFunc(args, func(s string) bool {
			return s == ">" || s == "1>"
		})
		if pos != -1 {
			if len(args) < pos + 1 {
				fmt.Printf("%s: Error: redirection requires a file argument\n", s.name)
				continue
			}
			file := args[pos + 1]
			args = args[:pos]
			absPath, nil := s.toAbs(file)
			if err != nil {
				fmt.Printf("%s: Error: %s\n", s.name, err)
        		continue
    		}	
	
			// Create the if it does not exist
			fd, err := os.Create(absPath)
			if err != nil {
				fmt.Printf("%s: %s\n", s.name, err)
			}
			defer fd.Close()
			s.stdin = fd
		}

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
	
	s := &Shell{
		cwd: cwd,
		path: path,
        loop:     true,
        builtins: make(map[string]func([]string)),
		defaultHistoryLimit: 16,
		name: "myshell",
		stdin: os.Stdin,
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
