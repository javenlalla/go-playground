package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	//execCommand()
	//execDockerComposeCommand()
	execRawBashCommand()
}

func execRawBashCommand() {
	rawCommand := "ls -l > sample.txt"

	cmd := exec.Command("bash", "-c", rawCommand)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command: %s. Command: %s", err, rawCommand)
	}
}

func execCommand() {
	cmd := exec.Command("ls", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatal(err)
	}

	log.Println(string(out), cmd)
}

func execDockerComposeCommand() {
	// Note: the `-T` flag is needed to get any output from execution within the container.
	cmd := exec.Command("docker-compose", "run", "-T", "aws-cli", "s3", "cp", "--acl", "public-read", "s3://origin_path", "s3://dest_path")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatal(err)
	}

	log.Println(string(out), cmd)
}
