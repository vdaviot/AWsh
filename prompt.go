package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type GitPrompt struct {
	Branch      string
	ToCommit    string
	NotCommited string
	Untracked   string
}

type Prompt struct {
	User    string
	Symbols []rune
	Git     GitPrompt
}

// "bash -c git rev-parse --abbrev-ref HEAD"
func getCommandOutput(command string) string {
	go timeTrack(time.Now(), "GET_COMMAND_OUTPUT")

	toExec := strings.Split(command, " ")

	cmd := exec.Command(toExec[0], toExec[1:]...)
	out, err := cmd.Output()
	if err != nil {
		go generateError(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

/*    ↑        ↓         ⚙            */
func formatGitPrompt(branch, status string, symbol rune) string {
	go timeTrack(time.Now(), "FORMAT_GIT_PROMPT")
	if strings.Contains(branch, "fatal") {
		return ""
	}

	var git GitPrompt
	git.Branch = branch
	if strings.Contains(status, "Modifications qui seront valid") || strings.Contains(status, "Changes to be committed:") {
		git.ToCommit = fmt.Sprintf(" ↑")
	} else {
		git.ToCommit = "-"
	}
	if strings.Contains(status, "Modifications qui ne seront pas valid") || strings.Contains(status, "Changes not staged for commit") {
		git.NotCommited = fmt.Sprintf(" ↓")
	} else {
		git.NotCommited = "-"
	}
	if strings.Contains(status, "Fichiers non suivis") || strings.Contains(status, "Unmerged path") {
		git.Untracked = fmt.Sprintf(" ⚙")
	} else {
		git.Untracked = "-"
	}
	return fmt.Sprintf("[%s %s %s %s ]", au.Gray(au.Bold(git.Branch)), git.ToCommit, git.NotCommited, git.Untracked)
}

// ⚙λ∞Δ
func generatePrompt() string {
	go timeTrack(time.Now(), "GENERATE_PROMPT")
	var user, git string
	var symbol = '∞'

	if ShellConfig.EnvVars.User != "" {
		user = ShellConfig.EnvVars.User
	} else {
		user = os.Getenv("USER")
		if user == "" {
			user = "Anonymous"
		}
	}
	pwd := os.Getenv("PWD")
	if _, err := os.Stat(fmt.Sprintf("%s/.git", pwd)); os.IsNotExist(err) || os.Getenv("PATH") == "" {
		return fmt.Sprintf("%s %s %s %s ", au.Bold(au.Green(user)), au.Brown(string(symbol)), au.Bold(au.Magenta("AlterWay™")), au.Bold("→"))
	}
	git = formatGitPrompt(getCommandOutput("git rev-parse --abbrev-ref HEAD"), getCommandOutput("git status"), symbol)
	return fmt.Sprintf("%s %s %s %s ", git, au.Bold(au.Green(user)), au.Brown(string(symbol)), au.Bold(au.Magenta("AlterWay™")))
}
