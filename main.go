//bin/cp .awshrc.json ~
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type TestScheme struct {
	TestResult []bool
	TestCase   [][]string
}

type EnvironmentScheme struct {
	Path string `json:"path,omitempty"`
	Home string `json:"home,omitempty"`
	User string `json:"user,omitempty"`
}

type ConfigScheme struct {
	Prompt        string
	Log           bool   `json:"log,omitempty"`
	LogFile       string `json:"logfile,omitempty"`
	Error         bool   `json:"error,omitempty"`
	ErrorFile     string `json:"errorfile,omitempty"`
	Bench         bool   `json:"bench,omitempty"`
	Color         bool   `json:"color,omitempty"`
	InfosBuiltins bool   `json:"infosbuiltins,omitempty"`
	EnvVars       EnvironmentScheme
}

/* Constant values */
const (
	separatorU        = "\n/-----------------------------------------\\"
	separatorD        = "\\-----------------------------------------/\n"
	banner     string = `
		  ___   _    _  _____  _            _  _ 
		 / _ \ | |  | |/  ___|| |          | || |
		/ /_\ \| |  | |\ '--. | |__    ___ | || |
		|  _  || |/\| | '--. \| '_ \  / _ \| || |
		| | | |\  /\  //\__/ /| | | ||  __/| || |
		\_| |_/ \/  \/ \____/ |_| |_| \___||_||_|
`
)

/* Global values */
var (
	au           = aurora.NewAurora(false)
	ShellConfig  ConfigScheme
	baseConfig   EnvironmentScheme
	oldState     *terminal.State
	signature    string
	helpMessage  = "Usage: ./AWShell [options]\n\toptions:\n\t  -b   \t\tDisable configuration file.\n\n\t  -c,   \t\tSpecify configuration file.\n\nAvailable Built-in's:\n\t- help\n\t- config\n"
	configFile   = ".awshrc.json"
	cmdBuffer    string
	term         = terminal.NewTerminal(screen, "awsh>")
	g_stdErrPipe = make(chan ErrorHandler)
	g_stdOutPipe = make(chan []byte, 1)
	screen       = struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
)

func writeToLogFile(args []string) {
	defer timeTrack(time.Now(), "WRITE_TO_LOGFILE")

	if ShellConfig.Log {
		out, err := os.OpenFile(ShellConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))
		defer out.Close()
		if err != nil {
			out = os.Stdout
			go generateError(err)
		}
		actualTime := time.Now()
		logMessage := fmt.Sprintf("[%d/%02d/%02d - %02d:%02d:%02d]-> %s\n",
			actualTime.Day(), actualTime.Month(), actualTime.Year(),
			actualTime.Hour(), actualTime.Minute(), actualTime.Second(), strings.Join(args, " "))
		out.WriteString(logMessage)
	}
}

func minimalConfig(scf *ConfigScheme) {
	defer timeTrack(time.Now(), "MINIMAL_CONFIGURATION")

	path, home, user := os.Getenv("PATH"), os.Getenv("HOME"), os.Getenv("USER")
	if len(path) == 0 || len(home) == 0 || len(user) == 0 {
		os.Setenv("PATH", "/bin:/usr/bin/:/sbin")
		os.Setenv("HOME", "/home")
		os.Setenv("USER", "Minimal")
	}
	*scf = ConfigScheme{
		Log:           false,
		LogFile:       "",
		Error:         true,
		ErrorFile:     ".AWShell.err",
		Bench:         false,
		Color:         false,
		InfosBuiltins: true,
		EnvVars: EnvironmentScheme{
			Path: os.Getenv("PATH"),
			Home: os.Getenv("HOME"),
			User: os.Getenv("USER"),
		},
	}
	printBuiltins("Minimal configuration loaded.\n")
	scf.Prompt = generatePrompt()
	term.SetPrompt(scf.Prompt)
}

func ReloadConfig(args []string) bool {
	defer timeTrack(time.Now(), "RELOAD_CONFIGURATION")

	if len(args) == 0 || len(args) == 1 && strings.TrimSpace(args[0]) == "" {
		minimalConfig(&ShellConfig)
		return true
	} else if len(args) > 0 && strings.HasSuffix(args[0], ".json") {

		file, err := os.OpenFile(args[0], os.O_RDONLY, os.FileMode(0666))
		defer file.Close()
		if err != nil {
			go generateErrorString(fmt.Sprintf("%s does not exist.\n", args[0]))
			return false
		}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&ShellConfig)
		if err != nil {
			go generateError(err)
			minimalConfig(&ShellConfig)
			return false
		}
		os.Setenv("PATH", ShellConfig.EnvVars.Path)
		os.Setenv("HOME", ShellConfig.EnvVars.Home)
		os.Setenv("USER", ShellConfig.EnvVars.User)
		ShellConfig.Prompt = generatePrompt()
		term.SetPrompt(ShellConfig.Prompt)
		printBuiltins(fmt.Sprintf("%s given, reloaded configuration.\n", args[0]))
		configFile = fmt.Sprintf("%s", args[0])
	} else {
		printBuiltins("No JSON specified.\n")
		minimalConfig(&ShellConfig)
		return false
	}
	return true
}

func fillConfig(color, blank, bench bool, scf *ConfigScheme, personnalConfigFile string) {
	defer timeTrack(time.Now(), "FILL_CONFIGURATION")

	au = aurora.NewAurora(color)
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		go generateErrorString("$GOPATH not set. Please install and setup Golang before using AWShell.\n")
		exitWrapper(oldState)
	}

	var configMessage string

	baseConfig = EnvironmentScheme{
		Path: os.Getenv("PATH"),
		Home: os.Getenv("HOME"),
		User: os.Getenv("USER"),
	}

	if blank == false { // Config file is activated
		file, err := os.Open(personnalConfigFile)
		defer file.Close()
		if err != nil {
			go generateErrorString(fmt.Sprintf("Could not open configuration file (%s does not exist).\n", personnalConfigFile))
			minimalConfig(scf)
			return
		}

		decoder := json.NewDecoder(file)
		err = decoder.Decode(scf)
		if err != nil {
			go generateErrorString(fmt.Sprintf("Could not load %s's content.\n", personnalConfigFile))
			minimalConfig(scf)
			return
		}
		os.Setenv("PATH", scf.EnvVars.Path)
		os.Setenv("HOME", scf.EnvVars.Home)
		os.Setenv("USER", scf.EnvVars.User)
		scf.LogFile = fmt.Sprintf("%s/src/AWShell/%s", os.Getenv("GOPATH"), scf.LogFile)
		scf.ErrorFile = fmt.Sprintf("%s/src/AWShell/%s", os.Getenv("GOPATH"), scf.ErrorFile)
		configMessage = "Configuration file loaded successfully.\n"
	} else {
		user := os.Getenv("USER")
		*scf = ConfigScheme{
			Log:           false,
			LogFile:       "",
			Error:         true,
			ErrorFile:     "",
			Bench:         false,
			Color:         true,
			InfosBuiltins: true,
			EnvVars: EnvironmentScheme{
				Path: fmt.Sprintf("/home/%s/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", user),
				Home: fmt.Sprintf("/home/%s", user),
				User: user,
			},
		}
		configMessage = "Blank mode activated. Base configuration loaded.\n"
	}
	scf.Color = color
	scf.Bench = bench
	configFile = fmt.Sprintf("%s%s", os.Getenv("HOME"), configFile[1:])
	fmt.Fprint(term, configMessage)
	scf.Prompt = generatePrompt()
	signature = fmt.Sprintf("\t\tMade with %s by %s\n\n", au.Bold(au.Red("<3")), au.Bold(au.Green("valentin.daviot@alterway.fr")))
	term.SetPrompt(scf.Prompt)
}

// Example:   defer timeTrack(time.Now(), "Minimal config")
func timeTrack(start time.Time, name string) {
	if ShellConfig.Bench {
		fmt.Fprintf(term, "%s Method %s took %s %s\n", au.Bold(au.Blue(string("Δ"))), au.Bold(au.Blue(name)), time.Since(start), au.Bold(au.Blue(string("Δ"))))
	}
}

func exitWrapper(oldState *terminal.State) {
	terminal.Restore(0, oldState)
	exitRoutine()
}

func signalHandler(sigs chan os.Signal) {
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGINT:
			fmt.Fprintf(term, "Received: %s (%d)\n", au.Brown(syscall.SIGINT), au.Brown(syscall.SIGINT))
		case syscall.SIGTERM:
			fmt.Fprintf(term, "Received: %s (%d)\n", au.Brown(syscall.SIGTERM), au.Brown(syscall.SIGTERM))
			exitWrapper(oldState)
		case syscall.SIGQUIT:
			fmt.Fprintf(term, "Received: %s (%d)\n", au.Brown(syscall.SIGQUIT), au.Brown(syscall.SIGQUIT))
			exitWrapper(oldState)
		default:
			fmt.Fprintf(term, "Received: %s\n", sig)
		}
	}
}

func main() {
	// Parsing
	var configMode = flag.Bool("wo", false, "Blank mode (without configuration file)")
	var colorMode = flag.Bool("color", false, "Disable color mode.")
	var benchMode = flag.Bool("bench", false, "Enable Benchmark mode")
	var personnalConfigFile = flag.String("config", configFile, "Specify configuration file")
	flag.Parse()

	// Configuration
	fillConfig(*colorMode, *configMode, *benchMode, &ShellConfig, *personnalConfigFile)
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// Terminal Initialisation
	term = terminal.NewTerminal(screen, ShellConfig.Prompt)
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)
	if err != nil {
		go generateErrorString("Terminal is not initialized.\nExiting...\n")
		time.Sleep(1 * time.Second)
		return
	}

	// Welcome message
	fmt.Fprint(term, au.Bold(au.Gray(banner)), signature)

	// Error/Stdout/Signal pipes
	go g_errorPipe()
	go g_outputPipe()
	go signalHandler(sigs) // Only catch signal sent from the outside

	// Main loop
	for {

		cmdBuffer, err = term.ReadLine()
		if err == io.EOF {
			exitWrapper(oldState)
		}

		if len(cmdBuffer) <= 0 {
			continue
		} else {

			// Formatting
			var cArg = strings.Split(cmdBuffer, ";")
			for _, cmd := range cArg {
				if len(cmd) <= 0 {
					continue
				}
				cArg = strings.Fields(cmd)

				// Execution
				if cArg[0] == "exit" {
					exitWrapper(oldState)
				}

				found := CheckCommandAvailability(cArg[0])
				lexedArgs := lexer(cArg)
				if found {
					handleBuiltins(cArg[0], lexedArgs)
					continue
				} else {

					path, err := exec.LookPath(cArg[0])
					if err != nil {
						go generateError(err)
						continue
					}
					handleCommand(path, lexedArgs)
				}
			}
			ShellConfig.Prompt = generatePrompt()
			term.SetPrompt(ShellConfig.Prompt)
		}
	}
}
