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

func execCommand(input_arr []string) error {
	bin := input_arr[0]
	path, err := exec.LookPath(bin)
	if err != nil {
		fmt.Printf("gosh: command not found: %s", bin)
	}
	var cmd *exec.Cmd
	if len(input_arr) == 0 {
		cmd = exec.Command(path)
	} else {
		rest := input_arr[1:]
		cmd = exec.Command(path, rest...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func parseEnv(input string) string {
	if strings.HasPrefix(input, "$") {
		input = os.ExpandEnv(input)
	}
	return input
}

func inputHandler(input string) (string, error) {
	input_arr := strings.Split(input, " ")
	switch input_arr[0] {
	// Nothing to see here, carry on
	case "clear":
		return "\033[H\033[2J", nil
	case "echo":
		if len(input_arr) == 0 {
			return "", nil
		}
		for s := range input_arr {
			input_arr[s] = parseEnv(input_arr[s])
		}
		return strings.Join(input_arr[1:], " "), nil
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
		err := execCommand(input_arr)
		return "", err
	}
}

func main() {
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
