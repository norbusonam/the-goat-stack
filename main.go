package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func performStepWithLogging(loadingMsg, completeMsg string, step func() error) {
	fmt.Println("⌛️ " + loadingMsg)
	err := step()
	if err != nil {
		logErrorAndExit(err.Error())
	}
	fmt.Printf("\033[1A\033[K")
	fmt.Println("✅ " + completeMsg)
}

func checkPreReq(cmd string) {
	performStepWithLogging("checking if "+cmd+" is installed", cmd+" is installed", func() error {
		_, err := exec.LookPath(cmd)
		if err != nil {
			return err
		}
		return nil
	})
}

func logErrorAndExit(err string) {
	fmt.Println("❌ " + err)
	os.Exit(1)
}

func logHelp() {
	fmt.Println("usage: thegoat <command>")
	fmt.Println("commands:")
	fmt.Println("  new\t\tCreate a new project")
	fmt.Println("  help\t\tShow this help message")
}

func createProject(pName, mName string) {
	// make project directory
	performStepWithLogging("creating project directory", "project directory created", func() error {
		return os.Mkdir(pName, fs.ModePerm)
	})
	// cd into project directory
	performStepWithLogging("changing directory", "directory changed", func() error {
		return os.Chdir(pName)
	})
	// create go module
	performStepWithLogging("creating go module", "go module created", func() error {
		return exec.Command("go", "mod", "init", mName).Run()
	})

	// TODO: create main.go that sets up echo server and serves static files

	// TODO: create pkg/templates w/ index.templ (incl tailwind and htmx)

	// TODO: create pkg/handlers w/ index.go

	// TODO: run templ generate

	// TODO: setup tailwind

	// TODO: setup air

	// TODO: create .gitignore

	// TODO: run go mod tidy
}

func main() {
	fmt.Print("thegoat 🐐\n\n")

	if len(os.Args) > 2 {
		logErrorAndExit("too many arguments")
	} else if len(os.Args) < 2 {
		logErrorAndExit("not enough arguments")
	}

	switch os.Args[1] {
	case "new":
		// get project name and module name
		fmt.Print("project name: ")
		var pName string
		fmt.Scanln(&pName)
		fmt.Print("module name: ")
		var mName string
		fmt.Scanln(&mName)
		fmt.Println()
		// check prerequisites
		fmt.Println("checking prerequisites")
		checkPreReq("go")
		checkPreReq("air")
		checkPreReq("npm")
		checkPreReq("templ")
		fmt.Println()
		// create project
		fmt.Println("creating project")
		createProject(pName, mName)
	case "help":
		logHelp()
	default:
		logErrorAndExit("invalid command")
	}
}
