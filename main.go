package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func execCommand(inputArr []string) error {
	bin := inputArr[0]
	if bin == "vi" || bin == "vim" {
		bin = "nvim"
	}
	path, err := exec.LookPath(bin)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if len(inputArr) == 0 {
		cmd = exec.Command(path)
	} else {
		rest := inputArr[1:]
		cmd = exec.Command(path, rest...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func parseEnv(inputArr []string) []string {
	for s := range inputArr {
		if strings.HasPrefix(inputArr[s], "$") {
			inputArr[s] = os.ExpandEnv(inputArr[s])
		}
	}
	return inputArr
}

func isBeautiful(name string) bool {
	return strings.HasSuffix(name, ".go")
}

func inputHandler(usr *user.User, input string) (string, error) {
	inputArr := parseEnv(strings.Split(input, " "))
	switch inputArr[0] {
	case "beautiful":
		if len(inputArr) == 1 {
			return "You are beautiful!\n", nil
		}
		ext := filepath.Ext(input)
		if isBeautiful(ext) {
			return "That file is beautiful!\n", nil
		}
		return "That file is not very beautiful...\n", nil
	case "cd":
		if len(inputArr) == 1 {
			return "", nil
		}
		path := inputArr[1]
		dir := usr.HomeDir
		var err error
		if path == "~" {
			err = os.Chdir(dir)
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(dir, path[2:])
			err = os.Chdir(path)
		} else {
			os.Chdir(path)
		}
		if err != nil {
			return "", err
		}
		return "", nil
	// Nothing to see here, carry on
	case "clear":
		return "\033[H\033[2J", nil
	case "echo":
		if len(inputArr) == 0 {
			return "", nil
		}
		return strings.Join(inputArr[1:], " "), nil
	// simple ReadDir
	case "ls":
		// A beautiful colour for beautiful files
		cyan := color.New(color.Bold, color.FgCyan).SprintFunc()
		blue := color.New(color.Bold, color.FgBlue).SprintFunc()
		files, err := ioutil.ReadDir("./")
		if err != nil {
			return "", err
		}
		var ret string
		for _, f := range files {
			name := f.Name()
			// Check if  the file is a directory
			if f.IsDir() {
				ret += fmt.Sprintf("%s\n", blue(name))
				// Check if the file is beautiful
			} else if isBeautiful(name) {
				ret += fmt.Sprintf("%s\n", cyan(name))
			} else {
				// Print out normal files :(
				ret += fmt.Sprintf("%s\n", name)
			}
		}
		return ret, nil
	default:
		err := execCommand(inputArr)
		return "", err
	}
}

func main() {
	// Because we're using gosh
	os.Setenv("SHELL", "gosh")
	reader := bufio.NewReader(os.Stdin)
	cyan := color.New(color.Bold, color.FgCyan).SprintFunc()
	green := color.New(color.Bold, color.FgGreen).SprintFunc()
	usr, _ := user.Current()
	for {
		// Imitate your favourite oh-my-zsh functionality with this one simple trick!
		dir, _ := os.Getwd()
		dir = strings.ReplaceAll(dir, usr.HomeDir, "~")
		fmt.Printf("%s \n%s ", cyan(dir), green("gs>"))
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "exit" || text == "quit" {
			return
		}
		res, err := inputHandler(usr, text)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
	}
}
