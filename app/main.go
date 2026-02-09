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
	var run = true
	for run {
		fmt.Print("$ ")
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command = strings.Trim(command, "\n")
		if command == "exit" {
			run = false
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
		fmt.Println(command + ": command not found")
	}
}
