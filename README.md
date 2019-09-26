#							AWShell

## Introduction

	AWShell is a SHELL made in go, in order to push my Golang's skills further.
	Feel free to contact me* if you find any bugs or have suggestions.
	This program was only tested on Ubuntu 16.04LTS

## Features

	- Configuration file (JSON), add -c YOURFILE.json to load yours
	- Multi commands interpretation
	- Redirections / Pipes 
	- Built-in's (cd / echo / export / unset / reload / config). Type help <yourcommand> for more infos
	- History (arrow up & down)
	- Copy & Paste
	- Clear Screen & Clear screen buffer
	- Errors and logs in $GOPATH/src/AWShell/.AWShell.err   .AWShell.log
	- Unit Tests, use go test to run it
	- Excution time for each method, if AWShell is invoked with --bench
	- Colors if AWShell is invoked with --color
	- Simple lexer (can use backquotes & $VAR)
	- Git integration in prompt

## Planned

	- Better lexer
	- Hash Table/Alias
	- Autocompletion
	
## Installation

		AWShell REQUIRE Golang with $GOPATH set.

		git clone  https://git.rnd.alterway.fr/valentin.daviot/AWShell.git
		cd AWShell
		go get -u github.com/golang/dep/cmd/dep
		dep ensure
		go build && go test && go install

		OR

		go get git.rnd.alterway.fr/valentin.daviot/AWShell.git
		dep ensure


## Configuration
	
	To properly use AWShell, you need a json configuration file.
	The basic configuration is described in .awshrc.json.

## Usage

	Basic:
		AWShell

	Recommanded:
		
		AWShell --color
	
	Benchmark

		AWShell --bench

### Contact

	If you need any informations, feel free to contact me at valentin.daviot@alterway.fr

	Thanks for using AWShell!

###           Made with <3 by vdaviot@AlterWay
