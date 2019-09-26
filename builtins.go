package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func printBuiltins(msg string) {
	if ShellConfig.InfosBuiltins {
		fmt.Fprint(term, msg)
	}
}

func ChangeDirectory(where []string) bool {
	defer timeTrack(time.Now(), "CHANGE_DIRECTORY")
	var oldPwd, _ = os.Getwd()
	var err error

	cmdLen := len(where)
	if cmdLen == 1 && strings.TrimSpace(where[0]) == "" {
		if oldPwd == os.Getenv("HOME") {
			return true
		}
		err = os.Chdir(os.Getenv("HOME"))
	} else if strings.Compare(where[0], "-") == 0 {
		err = os.Chdir(os.Getenv("OLDPWD"))
	} else if _, err = os.Stat(where[0]); err == nil {
		err = os.Chdir(where[0])
	}

	if err == nil {
		pwd, _ := os.Getwd()
		os.Setenv("OLDPWD", oldPwd)
		os.Setenv("PWD", pwd)
		printBuiltins(fmt.Sprintf("Moved from %s to %s\n", au.Blue(oldPwd), au.Blue(pwd)))
		if ShellConfig.Log {
			writeToLogFile(where)
		}
	} else {
		go generateErrorString(fmt.Sprintf("Could not move to %s\n", au.Red(where[1:])))
		return false
	}
	return true
}

func exitRoutine() {
	fmt.Fprint(term, "\n\t\tThanks for using AWShell.\r\n")
	os.Setenv("PATH", baseConfig.Path)
	os.Setenv("HOME", baseConfig.Home)
	os.Setenv("USER", baseConfig.User)
	os.Exit(0)
}

// Rajouter un beautifulPrint pour la config en memoire et pas en fichier
func showConfig() {
	defer timeTrack(time.Now(), "SHOW_CONFIG")
	file, err := os.OpenFile(configFile, os.O_RDONLY, os.FileMode(0666))
	defer file.Close()
	if err != nil {
		go generateError(err)
		return
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintf(term, "\n%s\n", scanner.Text())
	}
}

func SetEnvEntry(args []string) bool {
	defer timeTrack(time.Now(), "SET_ENV_ENTRY")

	argLen := len(args)
	if argLen <= 1 {
		go generateErrorString("Wrong export format. Expected \"export valueA=valueB\".\n")
		return false
	} else {
		args[1] = strings.TrimSpace(args[1])
		if args[1] == "" {
			go generateErrorString("Cannot export nothing.\n")
			return false
		} else if strings.Index(args[1], "=") == 0 {
			go generateErrorString("Wrong export format, =value is not the expected format.\n")
			return false
		} else if len(args[1]) < 2 {
			go generateErrorString("Wrong export format, please specify value.\n")
			return false
		}

	}
	// Capitalized value added
	newVar := strings.ToUpper(args[0])
	values := strings.Join(args[1:], ":")
	err := os.Setenv(newVar, values)
	if err != nil {
		go generateError(err)
		return false
	}
	printBuiltins(fmt.Sprintf("Added %s:[%s] to environment.\n", au.Blue(au.Bold(strings.ToUpper(newVar))), au.Blue(values)))
	return true
}

func UnsetEnvEntry(args []string) bool {
	defer timeTrack(time.Now(), "UNSET_ENV_ENTRY")

	if len(args) >= 1 {
		if args[0] == "" {
			go generateErrorString("Cannot unset nothing.\n")
			return false
		}
		for _, cmd := range args {
			envVar := os.Getenv(cmd)
			if envVar == "" {
				go generateErrorString(fmt.Sprintf("Variable %s doesn't exist.\n", au.Blue(cmd)))
				return false
			}
			if strings.Contains(cmd, "$") {
				go generateErrorString("Wrong unset usage, please remove $ (Should appear only in tests).\n")
				return false
			}
			err := os.Unsetenv(cmd)
			if err != nil {
				go generateError(err)
				return false
			}
			printBuiltins(fmt.Sprintf("Removed %s:[%s] from environment.\n", au.Bold(au.Blue(cmd)), au.Blue(envVar)))
		}
	} else {
		go generateErrorString("Wrong unset format. Usage: unset TO_UNSET\n")
		return false
	}
	return true
}

func help(cmd []string) {
	defer timeTrack(time.Now(), "HELP")
	if len(cmd) < 1 {
		fmt.Fprint(term, "Specify a command from cd / config / reload / export / unset / help to show more infos.\nPrompt symbols:\n\t|_ [<branch>]: Actual git branch\n\t|_ (↑): Changes to be committed\n\t|_ (↓): Changed not staged for commit\n\t|_ (⚙): Unmerged path\n")
		return
	}

	switch cmd[0] {
	case "cd":
		fmt.Fprint(term, "cd, Change Directory.\nUsage:\n\t|_ cd <folder>\n\t|_ cd - (goes to last folder, $OLDPWD)\n\t|_ cd $VAR\n\t|_ cd `command`\n\n")
	case "config":
		fmt.Fprint(term, "config, Show Configuration.\nUsage:\n\t|_ config (show your personnal configuration)\n\n")
	case "reload":
		fmt.Fprint(term, "reload, Reload Configuration.\nUsage:\n\t|_ reload <file.json> (Check .awshrc.json)\n\t|_ reload (Minimal config)\n\n")
	case "export":
		fmt.Fprint(term, "export, Export Environment Variable.\nUsage:\n\t|_ export VAR=VALUE\n\t|_ export VAR=$ENVVARIABLE\n\t|_ export VAR=$HOME:$PATH:`pwd` separated by ':'\n\n")
	case "unset":
		fmt.Fprint(term, "unset, Unset Environment Variable.\nUsage:\n\t|_ unset VALUE1 VALUE2 VALUE3\n\t|_ unset will stop at the first empty value\n\n")
	default:
		fmt.Fprintf(term, "\nReally this helpless ? Feel free to contact me at %s (%s)\n\n", au.Green(au.Bold("valentin.daviot@alterway.fr")), au.Red(au.Bold("<3")))
	}
}

func handleBuiltins(command string, args []string) {
	defer timeTrack(time.Now(), "HANDLE_BUILTIN")
	switch command {
	case "cd":
		ChangeDirectory(args)
	case "help":
		help(args)
	case "config":
		showConfig()
	case "reload":
		ReloadConfig(args)
	case "export":
		SetEnvEntry(args)
	case "unset":
		UnsetEnvEntry(args)
	}
}
