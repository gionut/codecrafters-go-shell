package main 

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"fmt"
	"os/exec"
	"os"
)

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
	os.Chdir(absPath)
}