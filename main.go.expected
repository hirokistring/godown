package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	VERSION           = "0.0.2"
	ENV_GO_COMMAND    = "GODOWN_GO_COMMAND"
	ENV_GOFMT_COMMAND = "GODOWN_GOFMT_COMMAND"
	GO_COMMAND        = "go"
	GOFMT_COMMAND     = "gofmt"
)

var (
	verbose   bool
	gopath    string
	gofmtpath string
)

var rootCmd = &cobra.Command{
	Use:     "godown <command> [arguments] ..",
	Version: VERSION,
	Short:   "'godown' is a 'go' command wrapper for Markdown files.",
	Long: `'godown' is a 'go' command wrapper for Markdown files.
You can generate '*.go' files from '*.go.md', etc.

For more information,
  https://github.com/hirokistring/godown

Example:
  $ cat hello.go.md

  # hello.go
` + "    ```" + `
    package main
    import "fmt"
    const WHO = "godown"
` + "    ```" + `
  # func main
` + "    ```" + `
    fmt.Printf("Hello, %s!\n", WHO)
` + "    ```" + `

  $ godown run hello.go.md
  Hello, godown!
`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		godownCommand(cmd, args)
	},
}

func godownCommand(cmd *cobra.Command, args []string) {
	// setup logging configurations
	setupLogging()

	// Find the input markdown files.
	markdown_files := findInputMarkdownFiles()
	for _, input_filename := range markdown_files {
		log.Printf("The input markdown file is %q\n", input_filename)

		// Read the input file.
		var md []byte
		var err error

		md, err = ioutil.ReadFile(input_filename)
		if err != nil {
			err = fmt.Errorf("failed to read bytes from %q: %v", input_filename, err)
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}

		// Generate the go source contents from the input markdown file.
		out, _, err := GenerateGoSourceFromMarkdown(input_filename, md)
		if err != nil {
			err = fmt.Errorf("failed to extract go code snippets from %q: %v", input_filename, err)
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}

		// Write the go source to the output file.
		output_filename := input_filename[:len(input_filename)-3] // trim '.md'
		log.Printf("generating go source file in %q\n", output_filename)
		err = ioutil.WriteFile(output_filename, out, 0644)
		if err != nil {
			err = fmt.Errorf("failed to write go source file in %q: %v", output_filename, err)
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}

		// Execute gofmt
		gofmtCommandPath := determinGoFmtCommandPath()

		// Execute gofmt command
		if exists(output_filename) {
			log.Printf("executing gofmt command %q for %q\n", gofmtCommandPath, output_filename)

			cmd := exec.Command(gofmtCommandPath, "-w", output_filename)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				err = fmt.Errorf("gofmt: gofmt command error: %v\n", err)
				fmt.Fprintln(os.Stderr, err)
				log.Fatal(err)
			}
		}
	}

	// Determine make command path
	goCommandPath := determinGoCommandPath()

	// Execute go command
	if len(args) > 0 {
		// Trim .md from the arguments
		log.Printf("trimming .md from the arguments: %q\n", args)
		var trimedArgs []string
		for i, v := range args {
			trimedArg := args[i]
			if strings.HasSuffix(trimedArg, ".md") {
				trimedArg = strings.TrimSuffix(v, ".md")
			}
			trimedArgs = append(trimedArgs, trimedArg)
		}
		log.Printf("trimmed arguments are %q\n", trimedArgs)

		// Pass the rest of command line arguments to the 'go' command.
		log.Printf("executing go command %q with args: %q\n", goCommandPath, trimedArgs)
		cmd := exec.Command(goCommandPath, trimedArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			err = fmt.Errorf("godown: go command error: %v\n", err)
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}
	}
}

func setupLogging() {
	log.SetPrefix("[godown] ")

	// Enable debug logs or not
	if verbose {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func askOverwrite(output_filename string) bool {
	entered := false
	overwrite := false
	var answer string

	for i := 0; !entered && i < 3; i++ {
		fmt.Printf("File %q already exists. Overwrite? (y/n): ", output_filename)
		_, err := fmt.Scanln(&answer)
		if err == nil {
			switch answer {
			case "Y", "y", "yes", "Yes", "YES":
				overwrite = true
				entered = true
			case "N", "n", "no", "No", "NO":
				overwrite = false
				entered = true
			}
		}
	}

	log.Printf("File %q already exists. Overwrite? (y/n): %s\n", output_filename, strconv.FormatBool(overwrite))

	return overwrite
}

func determineCommandPath(command string) string {
	command_path, err := exec.LookPath(command)
	if err != nil {
		err = fmt.Errorf("command not found in the PATH: %q", command)
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}
	return command_path // resolved command path
}

func determinGoCommandPath() string {
	return determineCommandPathBy(gopath, ENV_GO_COMMAND, GO_COMMAND)
}

func determinGoFmtCommandPath() string {
	return determineCommandPathBy(gofmtpath, ENV_GOFMT_COMMAND, GOFMT_COMMAND)
}

func determineCommandPathBy(argpath string, envvar string, defpath string) string {
	if gopath != "" {
		return determineCommandPath(argpath)
	}

	// Otherwise, check if the environment variable is set.
	env_value, env_set := os.LookupEnv(envvar)
	if env_set {
		return determineCommandPath(env_value)
	}

	// default go command
	return determineCommandPath(defpath)
}

func findInputMarkdownFiles() []string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		err := fmt.Errorf("getting the current working directory: %q", err)
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}

	// List files
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		err := fmt.Errorf("listing the current working directory: %q", err)
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}

	// find *.go.md files
	markdown_files := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go.md") {
			markdown_files = append(markdown_files, file.Name())
		}
	}

	return markdown_files
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "prints verbose messages")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}
}

// This file is generated from "main.go.md" by godown.
// https://github.com/hirokistring/godown
