package main

import (
	"bufio"
	"build/scripts"
	"build/types"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

const expectedAmountOfArguments = 2

func main() {
	if isRoot() {
		fmt.Println("Execute this program with root privileges!")
		os.Exit(1)
	}

	allowedCommands := [expectedAmountOfArguments]string{"dev", "prod"}
	arg := parseCommand(allowedCommands)

	var commands []types.Command

	switch arg {
	case "dev":
		commands = []types.Command{
			{Name: "Python Installation", Lines: scripts.Python},
		}
	case "prod":
		commands = []types.Command{
			{Name: "Docker Installation", Lines: scripts.Docker},
			{Name: "Nginx Installation", Lines: scripts.Nginx},
		}
	}

	installDependencies(commands)

	os.Exit(0)
}

func parseCommand(allowedCommands [expectedAmountOfArguments]string) string {
	args := os.Args

	if len(args) != expectedAmountOfArguments {
		fmt.Println("Execute this program passing only one environment type as argument.")
		os.Exit(1)
	}

	arg := args[1]

	if !isIn(arg, allowedCommands) {
		fmt.Printf("Invalid argument set. Choices are: %+q\n", allowedCommands)
		os.Exit(1)
	}

	return arg
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Printf("[isRoot] Unable to get current user: %s\n", err)
		os.Exit(1)
	}
	return currentUser.Username == "root"
}

func isIn(valueToSearch string, array [expectedAmountOfArguments]string) bool {
	for _, item := range array {
		if item == valueToSearch {
			return true
		}
	}
	return false
}

func displayStdout(wg *sync.WaitGroup, stdout io.ReadCloser) {
	defer wg.Done()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(string(scanner.Text()))
	}
}

func executeCommand(command types.Command) {

	cmd := exec.Command("bash", "-c", command.Lines)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error while executing the '%s' command\n", command.Name)
		os.Exit(1)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error while executing the '%s' command\n", command.Name)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go displayStdout(&wg, stdout)

	cmd.Wait()
	wg.Wait()

}

func installDependencies(commands []types.Command) {
	for _, command := range commands {
		executeCommand(command)
	}
}
