package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	major = 0
	minor = 1
	patch = 0
)

func spinner(loadingMsg, completedMsg string, done chan bool) {
	for {
		select {
		case <-done:
			// clear line
			fmt.Print("\r\033[K")
			fmt.Println("✅ " + completedMsg)
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c  %s", r, loadingMsg)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func performStepWithLogging(loadingMsg, completeMsg string, step func() error) {
	done := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		spinner(loadingMsg, completeMsg, done)
	}()
	err := step()
	done <- true
	wg.Wait()
	if err != nil {
		logErrorAndExit(err.Error())
	}
}

func checkPreReq(cmd string) {
	performStepWithLogging("Checking if "+cmd+" is installed", cmd+" is installed", func() error {
		_, err := exec.LookPath(cmd)
		if err != nil {
			return err
		}
		return nil
	})
}

func logErrorAndExit(err string) {
	fmt.Println("❌ " + err)
	fmt.Println("👉 goat help")
	os.Exit(1)
}

func createProject(pName, mName string, usingVscode bool) {
	performStepWithLogging("Creating project directories", "Project directories created", func() error {
		err := os.Mkdir(pName, fs.ModePerm)
		if err != nil {
			return err
		}
		err = os.Chdir(pName)
		if err != nil {
			return err
		}
		err = os.Mkdir("pkg", fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("Initializing Go module", "Go module initialized", func() error {
		return exec.Command("go", "mod", "init", mName).Run()
	})

	performStepWithLogging("Setting up Git", "Git set up", func() error {
		err := exec.Command("git", "init").Run()
		if err != nil {
			return err
		}
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

	performStepWithLogging("Creating templates package", "templates package created", func() error {
		err := os.Mkdir("pkg/templates", fs.ModePerm)
		if err != nil {
			return err
		}

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
		indexTemplContent += "\t\t\t<title>The Goat Stack</title>\n"
		indexTemplContent += "\t\t\t<link rel=\"stylesheet\" href=\"/tailwind.css\"/>\n"
		indexTemplContent += "\t\t\t<script src=\"https://unpkg.com/htmx.org@1.9.10\" integrity=\"sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC\" crossorigin=\"anonymous\"></script>\n"
		indexTemplContent += "\t\t</head>\n"
		indexTemplContent += "\t\t<body>\n"
		indexTemplContent += "\t\t\t<div class=\"flex flex-col items-center justify-center h-screen\">\n"
		indexTemplContent += "\t\t\t\t<h1 class=\"text-4xl font-bold text-gray-800\">Hello, { name }</h1>\n"
		indexTemplContent += "\t\t\t</div>\n"
		indexTemplContent += "\t\t</body>\n"
		indexTemplContent += "\t</html>\n"
		indexTemplContent += "}\n"
		_, err = f.WriteString(indexTemplContent)
		if err != nil {
			return err
		}
		return exec.Command("templ", "generate").Run()
	})

	performStepWithLogging("Creating handlers package", "handlers package created", func() error {
		err := os.Mkdir("pkg/handlers", fs.ModePerm)
		if err != nil {
			return err
		}
		f, err := os.Create("pkg/handlers/root.go")
		if err != nil {
			return err
		}
		defer f.Close()
		rootGoContent := "package handlers\n"
		rootGoContent += "\n"
		rootGoContent += "import (\n"
		rootGoContent += fmt.Sprintf("\t\"%s/pkg/templates\"\n", mName)
		rootGoContent += "\n"
		rootGoContent += "\t\"github.com/labstack/echo/v4\"\n"
		rootGoContent += ")\n"
		rootGoContent += "\n"
		rootGoContent += "func Root(c echo.Context) error {\n"
		rootGoContent += "\treturn templates.Hello(\"world 🐐\").Render(c.Request().Context(), c.Response().Writer)\n"
		rootGoContent += "}\n"
		_, err = f.WriteString(rootGoContent)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("Creating services package", "services package created", func() error {
		err := os.Mkdir("pkg/services", fs.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create("pkg/services/.gitkeep")
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("Creating db package", "db package created", func() error {
		err := os.Mkdir("pkg/db", fs.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create("pkg/db/.gitkeep")
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("Creating main.go", "main.go created", func() error {
		f, err := os.Create("main.go")
		if err != nil {
			return err
		}
		defer f.Close()
		mainGoContent := "package main\n"
		mainGoContent += "\n"
		mainGoContent += "import (\n"
		mainGoContent += fmt.Sprintf("\t\"%s/pkg/handlers\"\n", mName)
		mainGoContent += "\n"
		mainGoContent += "\t\"github.com/labstack/echo/v4\"\n"
		mainGoContent += ")\n"
		mainGoContent += "\n"
		mainGoContent += "func main() {\n"
		mainGoContent += "\te := echo.New()\n"
		mainGoContent += "\n"
		mainGoContent += "\te.Static(\"/\", \"public\")\n"
		mainGoContent += "\n"
		mainGoContent += "\te.GET(\"/\", handlers.Root)\n"
		mainGoContent += "\n"
		mainGoContent += "\te.Logger.Fatal(e.Start(\":8080\"))\n"
		mainGoContent += "}\n"
		_, err = f.WriteString(mainGoContent)
		if err != nil {
			return err
		}
		return nil
	})

	performStepWithLogging("Tidying Go module", "Go module tidied", func() error {
		return exec.Command("go", "mod", "tidy").Run()
	})

	performStepWithLogging("Setting up Tailwind", "Tailwind set up", func() error {
		err := exec.Command("npm", "install", "-D", "tailwindcss").Run()
		if err != nil {
			return err
		}
		err = exec.Command("npx", "tailwindcss", "init").Run()
		if err != nil {
			return err
		}
		// update tailwind.config.js
		fConfig, err := os.ReadFile("tailwind.config.js")
		if err != nil {
			return err
		}
		fConfigStr := string(fConfig)
		fConfigStr = strings.Replace(fConfigStr, "content: []", "content: [\"./pkg/templates/**/*.templ\"]", 1)
		err = os.WriteFile("tailwind.config.js", []byte(fConfigStr), fs.ModePerm)
		if err != nil {
			return err
		}
		// create input.css
		fInput, err := os.Create("input.css")
		if err != nil {
			return err
		}
		defer fInput.Close()
		_, err = fInput.WriteString("@tailwind base;\n@tailwind components;\n@tailwind utilities;\n")
		if err != nil {
			return err
		}
		// generate tailwind css
		return exec.Command("npx", "tailwindcss", "-i", "./input.css", "-o", "./public/tailwind.css", "--minify").Run()
	})

	performStepWithLogging("Setting up Air", "Air set up", func() error {
		err := exec.Command("air", "init").Run()
		if err != nil {
			return err
		}
		f, err := os.ReadFile(".air.toml")
		if err != nil {
			return err
		}
		fStr := string(f)
		// find cmd and add tailwind and templ steps
		cmdIdx := strings.Index(fStr, "cmd")
		cmdEndIdx := strings.Index(fStr[cmdIdx:], "\n")
		fStr = fStr[:cmdIdx] + "cmd = \"npx tailwindcss -i ./input.css -o ./public/tailwind.css --minify && templ generate && go build -o ./tmp/main .\"" + fStr[cmdIdx+cmdEndIdx:]
		// find exclude_dir list and append node_modules
		excludeDirIdx := strings.Index(fStr, "exclude_dir")
		insertIdx := strings.Index(fStr[excludeDirIdx:], "]")
		fStr = fStr[:excludeDirIdx+insertIdx] + ", \"node_modules\"" + fStr[excludeDirIdx+insertIdx:]
		// find exclude_regex list and append _templ.go
		excludeRegexIdx := strings.Index(fStr, "exclude_regex")
		insertIdx = strings.Index(fStr[excludeRegexIdx:], "]")
		fStr = fStr[:excludeRegexIdx+insertIdx] + ", \"_templ.go\"" + fStr[excludeRegexIdx+insertIdx:]
		// find include_ext list and append templ
		includeExtIdx := strings.Index(fStr, "include_ext")
		insertIdx = strings.Index(fStr[includeExtIdx:], "]")
		fStr = fStr[:includeExtIdx+insertIdx] + ", \"templ\"" + fStr[includeExtIdx+insertIdx:]
		// write changes
		err = os.WriteFile(".air.toml", []byte(fStr), fs.ModePerm)
		if err != nil {
			return err
		}
		return nil
	})

	if usingVscode {
		performStepWithLogging("Setting up VS Code", "VS Code set up", func() error {
			err := os.Mkdir(".vscode", fs.ModePerm)
			if err != nil {
				return err
			}
			fSettings, err := os.Create(".vscode/settings.json")
			if err != nil {
				return err
			}
			defer fSettings.Close()
			settingsJSONContent := "{\n"
			settingsJSONContent += "\t\"editor.formatOnSave\": true,\n"
			settingsJSONContent += "\t\"[templ]\": {\n"
			settingsJSONContent += "\t\t\"editor.defaultFormatter\": \"a-h.templ\"\n"
			settingsJSONContent += "\t},\n"
			settingsJSONContent += "\t\"tailwindCSS.includeLanguages\": {\n"
			settingsJSONContent += "\t\t\"templ\": \"html\"\n"
			settingsJSONContent += "\t}\n"
			settingsJSONContent += "}\n"
			_, err = fSettings.WriteString(settingsJSONContent)
			if err != nil {
				return err
			}
			fExt, err := os.Create(".vscode/extensions.json")
			if err != nil {
				return err
			}
			defer fExt.Close()
			extensionsJSONContent := "{\n"
			extensionsJSONContent += "\t\"recommendations\": [\n"
			extensionsJSONContent += "\t\t\"golang.go\",\n"
			extensionsJSONContent += "\t\t\"a-h.templ\",\n"
			extensionsJSONContent += "\t\t\"bradlc.vscode-tailwindcss\"\n"
			extensionsJSONContent += "\t]\n"
			extensionsJSONContent += "}\n"
			_, err = fExt.WriteString(extensionsJSONContent)
			if err != nil {
				return err
			}
			return nil
		})
	}

	fmt.Println()
	fmt.Println("🎉 Project created successfully")
	fmt.Println("👉 cd " + pName)
	fmt.Println("👉 air")
}

func main() {
	fmt.Print("The Goat Stack 🐐\n\n")

	if len(os.Args) > 2 {
		logErrorAndExit("Too many arguments")
	} else if len(os.Args) < 2 {
		logErrorAndExit("Not enough arguments")
	}

	switch os.Args[1] {
	case "new":
		// get project name and module name
		fmt.Print("Project name (default: my-project): ")
		var pName string
		fmt.Scanln(&pName)
		if pName == "" {
			pName = "my-project"
		}
		fmt.Print("Module name (default: my-module): ")
		var mName string
		fmt.Scanln(&mName)
		if mName == "" {
			mName = "my-module"
		}
		var vscodeRes string
		var usingVscode bool
		fmt.Print("Are you using VS Code? (y/n): ")
		fmt.Scanln(&vscodeRes)
		fmt.Println()
		usingVscode = vscodeRes == "y" || vscodeRes == "Y"
		// check prerequisites
		fmt.Println("Checking prerequisites:")
		checkPreReq("go")
		checkPreReq("air")
		checkPreReq("npm")
		checkPreReq("npx")
		checkPreReq("templ")
		checkPreReq("git")
		fmt.Println()
		// create project
		fmt.Println("Creating project:")
		createProject(pName, mName, usingVscode)
	case "help":
		fmt.Println("Usage: goat <command>")
		fmt.Println("Commands:")
		fmt.Println("  new\t\tcreate a new project")
		fmt.Println("  help\t\tshow this help message")
		fmt.Println("  version\tshow version information")
	case "version":
		fmt.Printf("v%d.%d.%d\n", major, minor, patch)
	default:
		logErrorAndExit("Invalid command")
	}
}
