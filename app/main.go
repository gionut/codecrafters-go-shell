package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


type Shell struct {
	commands map[string]func([]string)
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
	if _, ok := s.commands[name]; ok {
		fmt.Printf("%s is a shell builtin\n", name)
	} else {
		fmt.Printf("%s: not found\n", name)
	}
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
		
		if cmd, ok := s.commands[command]; ok {
			cmd(args)
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}

func NewShell() *Shell {
	scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer(make([]byte, 0, 1024), 1024)
    
	s := &Shell{
        loop:     true,
        commands: make(map[string]func([]string)),
		reader: scanner,
    }
    
    // Register commands here
    s.commands["exit"] = s._exit
    s.commands["type"] = s._type
    s.commands["echo"] = s._echo
    
    return s
}

func main() {
	shell := NewShell()
	shell.Loop()
}
