package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	commands := map[string]bool {
		"exit": true,
		"echo": true,
		"type": true,
	}

	var loop = true
	for loop {
		fmt.Print("$ ")
		
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command = strings.Trim(command, "\n")
		
		if command == "exit" {
			loop = false
			continue
		}

		var tokens = strings.SplitN(command, " ", 2)
		if len(tokens) == 0 {
			continue
		}
		command = tokens[0]
		if command == "echo" {
			if len(tokens) < 2 {
				continue
			}
			var args = tokens[1]
			fmt.Println(args)
			continue
		}
		if command == "type" {
			if len(tokens) < 2 {
				continue
			}
			var args string = tokens[1]
			_, ok := commands[args]
			if ok {
				fmt.Println(args + " is a shell builtin")
			} else {
				fmt.Println(args + ": not found")
			}
			continue
		}

		fmt.Println(command + ": command not found")
	}
}
