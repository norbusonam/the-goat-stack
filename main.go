package main

import (
	"fmt"
	"os"
	"os/exec"
)

func performStepWithLogging(loadingMsg string, completeMsg string, step func() error) {
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

func main() {
	fmt.Print("thegoat 🐐\n\n")

	if len(os.Args) > 2 {
		logErrorAndExit("too many arguments")
	} else if len(os.Args) < 2 {
		logErrorAndExit("not enough arguments")
	}

	switch os.Args[1] {
	case "new":
		// check prerequisites
		fmt.Println("checking prerequisites")
		checkPreReq("go")
		checkPreReq("air")
		checkPreReq("npm")
		checkPreReq("templ")
		fmt.Println()

	case "help":
		logHelp()
	default:
		logErrorAndExit("invalid command")
	}

	os.Exit(0)
}
