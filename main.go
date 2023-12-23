package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
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

	// setup tailwind
	performStepWithLogging("installing tailwind", "tailwind installed", func() error {
		return exec.Command("npm", "install", "-D", "tailwindcss").Run()
	})

	performStepWithLogging("initializing tailwind", "tailwind initialized", func() error {
		return exec.Command("npx", "tailwindcss", "init").Run()
	})

	performStepWithLogging("configure tailwind", "tailwind configured", func() error {
		f, err := os.ReadFile("tailwind.config.js")
		if err != nil {
			return err
		}
		fStr := string(f)
		fStr = strings.Replace(fStr, "content: []", "content: [\"./pkg/templates/**/*.templ\"]", 1)
		err = os.WriteFile("tailwind.config.js", []byte(fStr), fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("creating input css", "input css created", func() error {
		f, err := os.Create("input.css")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n")
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("generating tailwind css", "tailwind css generated", func() error {
		return exec.Command("npx", "tailwindcss", "-i", "./input.css", "-o", "./public/tailwind.css", "--minify").Run()
	})

	// setup air
	performStepWithLogging("initializing air", "air initialized", func() error {
		return exec.Command("air", "init").Run()
	})

	// TODO: configure air

	// create .gitignore
	performStepWithLogging("creating gitignore", "gitignore created", func() error {
		f, err := os.Create(".gitignore")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString("node_modules/\ntmp/\n*_templ.go\n")
		if err != nil {
			return err
		}
		return nil
	})

	// TODO: run go mod tidy

	// success message
	fmt.Println()
	fmt.Println("🎉 project created successfully")
	fmt.Println("👉 cd " + pName)
	fmt.Println("👉 air")
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
		fmt.Print("project name (default: my-project): ")
		var pName string
		fmt.Scanln(&pName)
		if pName == "" {
			pName = "my-project"
		}
		fmt.Print("module name (default my-module): ")
		var mName string
		fmt.Scanln(&mName)
		fmt.Println()
		if mName == "" {
			mName = "my-module"
		}
		// check prerequisites
		fmt.Println("checking prerequisites")
		checkPreReq("go")
		checkPreReq("air")
		checkPreReq("npm")
		checkPreReq("npx")
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
