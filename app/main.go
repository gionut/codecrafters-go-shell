package main

import (
	"fmt"
    "bufio"
    "os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	var run = true
	for run {
		fmt.Print("$ ")
		command, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		command = command[:len(command)-1]
		if command == "exit" {
			run = false
			continue
		}
		fmt.Println(command + ": command not found")
	}
}
