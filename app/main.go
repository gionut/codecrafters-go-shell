package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"os/exec"
)


type Shell struct {
	builtins map[string]func([]string)
	path string
	loop bool
	reader *bufio.Scanner
}

func (s *Shell) _exit(args []string) {
	s.loop = false
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

func (s *Shell) Loop() {
	for s.loop {
		fmt.Print("$ ")

		if !s.reader.Scan() {
            break
        }

		tokens := strings.Fields(s.reader.Text())
        if len(tokens) == 0 {
            continue
        }
		
		command := tokens[0]
		args := tokens[1:]
		
		if cmd, ok := s.builtins[command]; ok {
			cmd(args)
		} else { 
			fmt.Println(command + ": command not found")
		}
	}
}

func NewShell() *Shell {
	path := os.Getenv("PATH")
	scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer(make([]byte, 0, 1024), 1024)
    
	s := &Shell{
		path: path,
        loop:     true,
        builtins: make(map[string]func([]string)),
		reader: scanner,
    }
    
    // Register commands here
    s.builtins["exit"] = s._exit
    s.builtins["type"] = s._type
    s.builtins["echo"] = s._echo
    
    return s
}

func main() {
	shell := NewShell()
	shell.Loop()
}
