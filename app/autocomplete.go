package main

import (
	"strings"
	"os"
	"slices"
	"fmt"
	"github.com/samber/lo"
)

func findLcp(s1, s2 string) string {
    // Find the shorter length to avoid index out of bounds
    minLen := len(s1)
    if len(s2) < minLen {
        minLen = len(s2)
    }

    for i := 0; i < minLen; i++ {
        // If characters at index i don't match, 
        // return the string up to that index
        if s1[i] != s2[i] {
            return s1[:i]
        }
    }

    // If the loop finishes, one string is a prefix of the other
    return s1[:minLen]
}

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

	// Prepare the file names from cwd
	files = lo.Filter(files, func(f os.DirEntry, _ int) bool {
    	return strings.HasPrefix(f.Name(), prefix)
	})
	names := lo.Map(files, func(f os.DirEntry, _ int) string {
    		if f.IsDir() {
				return f.Name() + "/"
			}
			return f.Name()
		})
	slices.Sort(names)
	
	if len(names) > 1 {
		// Complete until LCP
		lcp := findLcp(names[0], names[1])
		if len(lcp) > len(prefix) {
			suffix := lcp[len(prefix):]
			newLine := append(line, []rune(suffix)...)
			newPos := pos + len(suffix)

			return newLine, newPos, true
		}

		if !bellRang {
			return append(line, '\x07'), pos + 1, true
		}

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