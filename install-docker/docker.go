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

	// Docker
	cmd := exec.Command("docker", "-v")
	output, err := cmd.Output()
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

	commands := []string{installDockerScript}

	for idx, command := range commands {
		executeCommand(idx, command)
	}

	showProgramVersions()
}
