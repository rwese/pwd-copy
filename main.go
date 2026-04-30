package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func detectClipboard() string {
	if cmd := os.Getenv("PWD_COPY_CLIPBOARD_COMMAND"); cmd != "" {
		return cmd
	}
	for _, cmd := range []string{"pbcopy", "wl-copy", "xclip", "xsel"} {
		if _, err := exec.LookPath(cmd); err == nil {
			return cmd
		}
	}
	fmt.Fprintln(os.Stderr, "ERROR: No clipboard command found. Set PWD_COPY_CLIPBOARD_COMMAND")
	os.Exit(1)
	return ""
}

func runClipboard(clipboard, content string) {
	var cmd *exec.Cmd
	
	// Replace placeholders
	clipboard = strings.ReplaceAll(clipboard, "%p", fmt.Sprintf("%q", content))
	clipboard = strings.ReplaceAll(clipboard, "%raw_p", content)

	// If using a known clipboard tool directly, use stdin
	switch clipboard {
	case "pbcopy":
		cmd = exec.Command("pbcopy")
	case "wl-copy":
		cmd = exec.Command("wl-copy")
	case "xsel":
		cmd = exec.Command("xsel", "--clipboard", "--input")
	case "xclip":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		// Custom command - parse and run
		args := parseShellLike(clipboard)
		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "ERROR: Invalid clipboard command")
			os.Exit(1)
		}
		cmd = exec.Command(args[0], args[1:]...)
	}


	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	io.WriteString(stdin, content+"\n")
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: clipboard failed: %s\n", err)
		os.Exit(1)
	}
}

func parseShellLike(s string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := '"'

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"', '\'':
			if !inQuote {
				inQuote = true
				quoteChar = rune(c)
			} else if rune(c) == quoteChar {
				inQuote = false
			} else {
				current.WriteByte(c)
			}
		case ' ':
			if inQuote {
				current.WriteByte(c)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}

func main() {
	var path string

	switch {
	case len(os.Args) == 1:
		// Current directory full path
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case os.Args[1] == "-r":
		// Relative path from cwd
		target := "."
		if len(os.Args) > 2 {
			target = os.Args[2]
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		targetAbs, err := filepath.Abs(target)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		path, err = filepath.Rel(cwd, targetAbs)
		if err != nil {
			path = targetAbs
		}
	case os.Args[1] == "-h" || os.Args[1] == "--help":
		fmt.Println("Usage: pwd-copy [-r] [target]")
		fmt.Println("  (no args)  Copy current directory full path")
		fmt.Println("  target     Copy target directory full path")
		fmt.Println("  -r target  Copy relative path from cwd to target")
		fmt.Println()
		fmt.Println("Environment:")
		fmt.Println("  PWD_COPY_CLIPBOARD_COMMAND  Custom clipboard command")
		fmt.Println("    %p    Replaced with path wrapped in quotes")
		fmt.Println("    %raw_p Replaced with raw path (may break with spaces)")
		os.Exit(0)
	default:
		// Target directory full path
		var err error
		path, err = filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}

	clipboard := detectClipboard()
	runClipboard(clipboard, path)
}
