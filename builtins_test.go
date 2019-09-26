package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"testing"
)

var a = aurora.NewAurora(true)

/*TestCheckCommandAvailability tries every command */
func TestCheckCommandAvailability(t *testing.T) {
	ShellConfig.Error = false
	ShellConfig.InfosBuiltins = false

	var ts = TestScheme{}
	ts.TestResult = []bool{
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
	}
	ts.TestCase = [][]string{
		{"cd"},
		{"help"},
		{"reload"},
		{"unset"},
		{"exit"},
		{"config"},
		{"export"},
		{""},
		{"pokemon"},
	}

	fmt.Printf("%s\n    Entering %s\n%s\n", a.Bold(separatorU), a.Bold(a.Blue("CHECK_COMMAND_AVAILABILITY")), a.Bold(separatorD))

	for index, cmd := range ts.TestCase {
		switch CheckCommandAvailability(cmd[0]) == ts.TestResult[index] {
		case true:
			t.Logf("[%s (%d)]\n\t%s\n\n", a.Green(a.Bold("OK")), index, ts.TestCase[index][0])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Green(a.Bold("OK")), cmd[0], ts.TestResult[index])
		case false:
			t.Errorf("[%s (%d)]\n\t%s\n\n", a.Red(a.Bold("KO")), index, ts.TestCase[index][0])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Red(a.Bold("KO")), cmd[0], ts.TestResult[index])
		}
	}
}

/*TestChangeDirectory does...*/
func TestChangeDirectory(t *testing.T) {
	ShellConfig.Log = false
	ShellConfig.InfosBuiltins = false

	var ts = TestScheme{}
	ts.TestResult = []bool{
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}
	ts.TestCase = [][]string{
		/* Case True */
		{""},
		{" "},
		{"../"},
		{"\t"},
		/* Case False */
		{"`pwd`"},
		{"$PATH"},
		{"/dev/stdin"},
		{"`ls`"},
		{"/home/vdaviot/work/go/src/AWShell/awshrc.json"},
	}

	fmt.Printf("%s\n    Entering %s\n%s\n", a.Bold(separatorU), a.Bold(a.Blue("TEST_CHANGE_DIRECTORY")), a.Bold(separatorD))

	for index, cmd := range ts.TestCase {
		switch ChangeDirectory(cmd) == ts.TestResult[index] {
		case true:
			t.Logf("[%s (%d)]\n\t%s\n\n", a.Green(a.Bold("OK")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Green(a.Bold("OK")), cmd[0], ts.TestResult[index])
		case false:
			t.Errorf("[%s (%d)]\n\t%s\n\n", a.Red(a.Bold("KO")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Red(a.Bold("KO")), cmd[0], ts.TestResult[index])
		}
	}
}

/*TestSetEnvEntry does...*/
func TestSetEnvEntry(t *testing.T) {
	ShellConfig.Log = false
	ShellConfig.InfosBuiltins = false

	var ts = TestScheme{}
	ts.TestResult = []bool{
		true,
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		true,
	}
	ts.TestCase = [][]string{
		/* Case True */
		{"export", "TEST=b"},
		{"export", "CHEVRE", "=", "Meeeeeh"},
		{"export", "TEST=$TEST:$TEST"},
		{"export", "HOMEPATH=$HOME:$PATH"},
		{"export", "HARD_TEST=`pwd`:$HOMEPATH:fdp"}, // Not working as intended
		/* Case False */
		{"export", ""},
		{"export", "=$PATH"},
		{"export", "="},
		{"export", "b"},
		{"export", "TEST=$$$TEST"},
	}

	fmt.Printf("%s\n    Entering %s\n%s\n", a.Bold(separatorU), a.Bold(a.Blue("TEST_SET_ENV_ENTRY")), a.Bold(separatorD))

	for index, cmd := range ts.TestCase {
		switch SetEnvEntry(cmd) == ts.TestResult[index] {
		case true:
			t.Logf("[%s (%d)]\n\t%s\n\n", a.Green(a.Bold("OK")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s %s] -> %t\n", a.Green(a.Bold("OK")), cmd[0], cmd[1], ts.TestResult[index])
		case false:
			t.Errorf("[%s (%d)]\n\t%s\n\n", a.Red(a.Bold("KO")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s %s] -> %t\n", a.Red(a.Bold("KO")), cmd[0], cmd[1], ts.TestResult[index])
		}
	}
}

/*TestReloadConfig tries multiple configuration reload*/
func TestReloadConfig(t *testing.T) {
	ShellConfig.Log = false
	ShellConfig.InfosBuiltins = false

	var ts = TestScheme{}
	ts.TestResult = []bool{
		true,
		true,
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}
	ts.TestCase = [][]string{
		/* Case True */
		{""},
		{" "},
		{"               "},
		{"\t"},
		/* Case False */
		{"empty.json"},
		{"pokemon.yaml"},
		{"../"},
		{"$HOME"},
		{"AWShell"},
	}

	fmt.Printf("%s\n    Entering %s\n%s\n", a.Bold(separatorU), a.Bold(a.Blue("TEST_RELOAD_CONFIG")), a.Bold(separatorD))

	for index, cmd := range ts.TestCase {
		switch ReloadConfig(cmd) == ts.TestResult[index] {
		case true:
			t.Logf("[%s (%d)]\n\t%s\n\n", a.Green(a.Bold("OK")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Green(a.Bold("OK")), cmd[0], ts.TestResult[index])
		case false:
			t.Errorf("[%s (%d)]\n\t%s\n\n", a.Red(a.Bold("KO")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s] -> %t\n", a.Red(a.Bold("KO")), cmd[0], ts.TestResult[index])
		}
	}
}

/*TestUnsetEnvEntry does...*/
func TestUnsetEnvEntry(t *testing.T) {
	ShellConfig.Log = false
	ShellConfig.InfosBuiltins = false

	var ts = TestScheme{}
	ts.TestResult = []bool{
		true,
		true,
		false,
		false,
		false,
		false,
		false,
	}

	ts.TestCase = [][]string{
		{"SHLVL"},
		{"HOME", "PATH", "USER"},
		{"$HOME", "$PATH", "$USER"},
		{"$SHLVL"},
		{"pokemon"},
		{"a"},
		{""},
	}

	fmt.Printf("%s\n    Entering %s\n%s\n", a.Bold(separatorU), a.Bold(a.Blue("TEST_UNSET_ENV_ENTRY")), a.Bold(separatorD))

	for index, cmd := range ts.TestCase {
		switch UnsetEnvEntry(cmd) == ts.TestResult[index] {
		case true:
			// t.Logf("[%s (%d)]\n\t%s\n\n", a.Green(a.Bold("OK")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s %t]\n", a.Green(a.Bold("OK")), cmd[0], ts.TestResult[index])
		case false:
			t.Errorf("[%s (%d)]\n\t%s\n\n", a.Red(a.Bold("KO")), index, ts.TestCase[index])
			fmt.Fprintf(term, "\t%s: [%s %t]\n", a.Red(a.Bold("KO")), cmd[0], ts.TestResult[index])
		}
	}
}
