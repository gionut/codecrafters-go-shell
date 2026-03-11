package main

import (
	"strings"
	"os"
	"slices"
	"fmt"
	"github.com/samber/lo"
)

func (s *Shell) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	if key != '\t' {
		return nil, 0, false
	}
	
	// Remove tab
	line = line[:pos-1]
	pos--
	
	currentText := string(line)
	i := strings.LastIndex(currentText, " ")
	prefix := currentText[i+1:]
	searchDir := s.cwd

	bellRang := false
	if strings.Contains(prefix, string('\x07')) {
		bellRang = true
		prefix = prefix[:len(prefix)-1]
		line = line[:pos-1]
		pos--
	}

	// Nested file completion
	if strings.Contains(prefix, "/") {
		i := strings.LastIndex(prefix, "/")
		searchDir += "/" + prefix[:i]
		prefix = prefix[i+1:]
	}

	files, err := os.ReadDir(searchDir)
	if err != nil {
		return line, pos, true
	}

	matchesCnt := lo.CountBy(files, func(f os.DirEntry) bool {
    	return strings.HasPrefix(f.Name(), prefix)
	})

	if matchesCnt > 1 {
		if !bellRang {
			return append(line, '\x07'), pos + 1, true
		}
		names := lo.Map(files, func(f os.DirEntry, _ int) string {
    		if f.IsDir() {
				return f.Name() + "/"
			}
			return f.Name()
		})
		slices.Sort(names)
		fmt.Printf("\n%s\n", strings.Join(names, "  "))
		return line, pos, true
	}

	for _, file := range files {
		name := file.Name()
		if strings.HasPrefix(name, prefix) {
			suffix := name[len(prefix):] + " "
			if file.IsDir() {
				suffix = name[len(prefix):] + "/"
			}

			newLine := append(line, []rune(suffix)...)
			newPos := pos + len(suffix)

			return newLine, newPos, true
		}
	}
	
	// No match return the bell '\x07' charactrer
	return append(line, '\x07'), pos + 1, true
}