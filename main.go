package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

const (
	major = 0
	minor = 1
	patch = 0
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

func createProject(pName, mName string) {
	performStepWithLogging("creating project directory", "project directory created", func() error {
		return os.Mkdir(pName, fs.ModePerm)
	})

	performStepWithLogging("changing directory", "directory changed", func() error {
		return os.Chdir(pName)
	})

	performStepWithLogging("creating go module", "go module created", func() error {
		return exec.Command("go", "mod", "init", mName).Run()
	})

	performStepWithLogging("creating main.go", "main.go created", func() error {
		f, err := os.Create("main.go")
		if err != nil {
			return err
		}
		defer f.Close()
		mainGoContent := "package main\n"
		mainGoContent += "\n"
		mainGoContent += "import (\n"
		mainGoContent += fmt.Sprintf("\t\"%s/pkg/templates\"\n", mName)
		mainGoContent += "\n"
		mainGoContent += "\t\"github.com/labstack/echo/v4\"\n"
		mainGoContent += ")\n"
		mainGoContent += "\n"
		mainGoContent += "func main() {\n"
		mainGoContent += "\te := echo.New()\n"
		mainGoContent += "\n"
		mainGoContent += "\te.Static(\"/\", \"public\")\n"
		mainGoContent += "\n"
		mainGoContent += "\te.GET(\"/\", func(c echo.Context) error {\n"
		mainGoContent += "\t\treturn templates.Hello(\"world 🐐\").Render(c.Request().Context(), c.Response().Writer)\n"
		mainGoContent += "\t})\n"
		mainGoContent += "\n"
		mainGoContent += "\te.Logger.Fatal(e.Start(\":8080\"))\n"
		mainGoContent += "}\n"
		_, err = f.WriteString(mainGoContent)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("creating pkg directory", "pkg directory created", func() error {
		err := os.Mkdir("pkg", fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("creating templates directory", "templates directory created", func() error {
		err := os.Mkdir("pkg/templates", fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("creating index.templ", "index.templ created", func() error {
		f, err := os.Create("pkg/templates/index.templ")
		if err != nil {
			return err
		}
		defer f.Close()
		indexTemplContent := "package templates\n"
		indexTemplContent += "\n"
		indexTemplContent += "templ Hello(name string) {\n"
		indexTemplContent += "\t<!DOCTYPE html>\n"
		indexTemplContent += "\t<html lang=\"en\">\n"
		indexTemplContent += "\t\t<head>\n"
		indexTemplContent += "\t\t\t<meta charset=\"UTF-8\"/>\n"
		indexTemplContent += "\t\t\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"/>\n"
		indexTemplContent += "\t\t\t<title>thegoat 🐐</title>\n"
		indexTemplContent += "\t\t\t<link rel=\"stylesheet\" href=\"/tailwind.css\"/>\n"
		indexTemplContent += "\t\t</head>\n"
		indexTemplContent += "\t\t<body>\n"
		indexTemplContent += "\t\t\t<div class=\"flex flex-col items-center justify-center h-screen\">\n"
		indexTemplContent += "\t\t\t\t<h1 class=\"text-4xl font-bold text-gray-800\">Hello, {name}</h1>\n"
		indexTemplContent += "\t\t\t</div>\n"
		indexTemplContent += "\t\t</body>\n"
		indexTemplContent += "\t</html>\n"
		indexTemplContent += "}\n"
		_, err = f.WriteString(indexTemplContent)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("generating templates", "templates generated", func() error {
		return exec.Command("templ", "generate").Run()
	})

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

	performStepWithLogging("initializing air", "air initialized", func() error {
		return exec.Command("air", "init").Run()
	})

	performStepWithLogging("configuring air", "air configured", func() error {
		f, err := os.ReadFile(".air.toml")
		if err != nil {
			return err
		}
		fStr := string(f)
		// find exclude_dir list and append node_modules
		excludeDirIdx := strings.Index(fStr, "exclude_dir")
		insertIdx := strings.Index(fStr[excludeDirIdx:], "]")
		fStr = fStr[:excludeDirIdx+insertIdx] + ", \"node_modules\"" + fStr[excludeDirIdx+insertIdx:]
		// find exclude_regex list and append _templ.go
		excludeRegexIdx := strings.Index(fStr, "exclude_regex")
		insertIdx = strings.Index(fStr[excludeRegexIdx:], "]")
		fStr = fStr[:excludeRegexIdx+insertIdx] + ", \"_templ.go\"" + fStr[excludeRegexIdx+insertIdx:]
		// find include_ext list and append .templ
		includeExtIdx := strings.Index(fStr, "include_ext")
		insertIdx = strings.Index(fStr[includeExtIdx:], "]")
		fStr = fStr[:includeExtIdx+insertIdx] + ", \".templ\"" + fStr[includeExtIdx+insertIdx:]
		// write changes
		err = os.WriteFile(".air.toml", []byte(fStr), fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

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

	performStepWithLogging("tidying go modules", "go modules tidied", func() error {
		return exec.Command("go", "mod", "tidy").Run()
	})

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
		fmt.Println("usage: thegoat <command>")
		fmt.Println("commands:")
		fmt.Println("  new\t\tcreate a new project")
		fmt.Println("  help\t\tshow this help message")
		fmt.Println("  version\tshow version information")
	case "version":
		fmt.Printf("v%d.%d.%d\n", major, minor, patch)
	default:
		logErrorAndExit("invalid command")
	}
}
