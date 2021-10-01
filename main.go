package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func execCommand(inputArr []string) error {
	bin := inputArr[0]
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

func inputHandler(input string) (string, error) {
	inputArr := parseEnv(strings.Split(input, " "))
	switch inputArr[0] {

	case "cd":
		err := os.Chdir(inputArr[1])
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
		files, err := ioutil.ReadDir("./")
		if err != nil {
			return "", fmt.Errorf("something strange happened")
		}
		var ret string
		for _, f := range files {
			ret += fmt.Sprintf("%s\n", f.Name())
		}
		return ret, nil
	default:
		err := execCommand(inputArr)
		return "", err
	}
}

func main() {
	// Because we're using goshell
	os.Setenv("SHELL", "goshell")
	reader := bufio.NewReader(os.Stdin)
	cyan := color.New(color.Bold, color.FgCyan).SprintFunc()
	for {
		// Imitate your favourite oh-my-zsh functionality with this one simple trick!
		dir, _ := os.Getwd()
		fmt.Printf("%s \n> ", cyan(dir))
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		if text == "exit" || text == "quit" {
			return
		}
		res, err := inputHandler(text)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res + "\n")
	}
}
