package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func CheckCommandAvailability(command string) bool {
	defer timeTrack(time.Now(), "COMMAND_AVAILABILITY")
	if strings.Compare(command, "cd") == 0 || strings.Compare(command, "exit") == 0 ||
		strings.Compare(command, "help") == 0 || strings.Compare(command, "config") == 0 ||
		strings.Compare(command, "reload") == 0 || strings.Compare(command, "export") == 0 ||
		strings.Compare(command, "unset") == 0 {
		return true
	}
	return false
}

// Goroutine: Catch output
func g_outputPipe() {
	for elem := range g_stdOutPipe {
		toPrint := strings.Replace(string(elem), "\n", "\b\n", -1)
		fmt.Fprint(term, string(toPrint))
	}
}

func lexer(args []string) []string {
	defer timeTrack(time.Now(), "LEXER")
	var lexed []string

	for _, cmd := range args[1:] {
		cmd = strings.TrimSpace(cmd)
		if strings.HasPrefix(cmd, "--") && strings.Contains(cmd, "=") && len(cmd) > 3 {
			lexed = append(lexed, cmd)
			continue
		}
		for _, param := range strings.Split(cmd, "=") {
			for _, value := range strings.Split(param, ":") {
				if strings.HasPrefix(value, "$") && strings.Count(value, "$") == 1 {
					val := os.Getenv(value[1:])
					if val == "" {
						val = value
					}
					lexed = append(lexed, val)
				} else if strings.Contains(value, "`") && (strings.Count(value, "`")%2) == 0 {
					lexed = append(lexed, getCommandOutput(value[1:len(value)-1]))
				} else if strings.Contains(value, "~") {
					index := strings.Index(value, "~")
					lexed = append(lexed, value[:index]+os.Getenv("HOME")+value[index+1:])
				} else {
					lexed = append(lexed, value)
				}
			}
		}
	}
	return lexed
}

func handleCommand(command string, args []string) {
	if ShellConfig.Color == true && command == "/bin/ls" {
		command = "/bin/ls --color=always"
	}
	defer timeTrack(time.Now(), "HANDLE_COMMAND")
	var stdOutBuf, stdErrBuf bytes.Buffer

	command = fmt.Sprintf("%s %s", command, strings.Join(args, " "))

	// Preparing command
	cmd := exec.Command("sh", "-c", command)
	cmdReaderIn, _ := cmd.StdoutPipe()
	cmd.Stdin = os.Stdin
	cmd.Stderr = &stdErrBuf
	if strings.Contains(command, "vim") || strings.Contains(command, "top") {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdOutBuf
	}

	// Redirect outputs
	stdout := io.MultiWriter(os.Stdout, &stdOutBuf)

	// Goroutine forward stdout
	go func() {
		_, _ = io.Copy(stdout, cmdReaderIn)
	}()

	// Run
	err := cmd.Start()
	if err != nil {
		go generateErrorString(stdErrBuf.String())
	}
	err = cmd.Wait()
	if err != nil {
		go generateErrorString(stdErrBuf.String())
	}

	// Get output
	g_stdOutPipe <- stdOutBuf.Bytes()

	// Log
	if ShellConfig.Log {
		writeToLogFile(cmd.Args[2:])
	}
}
