package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"sync"
)

const setupScript = `
	sudo apt update
	sudo apt upgrade -y
	sudo apt install curl -y
`

const installPythonScript = `
	sudo apt update
	sudo apt upgrade -y

	sudo apt install software-properties-common build-essential -y

	sudo add-apt-repository ppa:deadsnakes/ppa -y

	sudo apt update

	sudo apt install python3.11 python3.11-venv -y
`

const installNodejsAndYarnScript = `
	sudo apt update
	sudo apt upgrade -y

	curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

	sudo apt install nodejs -y

	sudo apt update
	sudo apt upgrade -y

	curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
	echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list

	sudo apt update
	sudo apt install yarn -y
`

const installDockerScript = `
	sudo apt update
	sudo apt upgrade -y

	sudo apt remove docker docker-engine docker.io containerd runc
	sudo apt install apt-transport-https ca-certificates curl gnupg lsb-release -y

	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
	echo \
	"deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
	$(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

	sudo apt update
	sudo apt install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y

	sudo service docker start
`

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func displayStdout(wg *sync.WaitGroup, stdout io.ReadCloser) {
	defer wg.Done()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(string(scanner.Text()))
	}
}

func executeCommand(order int, command string) {

	cmd := exec.Command("bash", "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error while executing the %dº command: %s\n", order+1, err)
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error while executing the %dº command: %s\n", order, err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go displayStdout(&wg, stdout)

	wg.Wait()
	cmd.Wait()
}

func showProgramVersions() {

	// Python
	cmd := exec.Command("python3.11", "--version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting python's version: %s", err)
	} else {
		fmt.Printf("Python's version: %s", output)
	}

	// Nodejs
	cmd = exec.Command("node", "-v")
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting node's version: %s", err)
	} else {
		fmt.Printf("Node.js's version: %s", output)
	}

	// Npm
	cmd = exec.Command("npm", "-v")
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting npm's version: %s", err)
	} else {
		fmt.Printf("Npm's version: %s", output)
	}

	// Yarn
	cmd = exec.Command("yarn", "-v")
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting yarn's version: %s", err)
	} else {
		fmt.Printf("Yarn's version: %s", output)
	}

	// Docker
	cmd = exec.Command("docker", "-v")
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting docker's version: %s", err)
	} else {
		fmt.Printf(string(output))
	}

	// Docker Compose
	cmd = exec.Command("docker", "compose", "version")
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("An error occurred while getting docker's compose version: %s", err)
	} else {
		fmt.Printf(string(output))
	}

	fmt.Println("Remember to give permission for your default user in order to use docker. Type the following command:")
	fmt.Println("sudo usermod -aG docker $USER")
}

func main() {
	if !isRoot() {
		fmt.Println("Execute this program with root privileges!")
		os.Exit(0)
	}

	commands := []string{setupScript, installPythonScript, installNodejsAndYarnScript, installDockerScript}

	for idx, command := range commands {
		executeCommand(idx, command)
	}

	showProgramVersions()
}
