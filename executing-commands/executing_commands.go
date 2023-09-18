package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//execCommand()
	//execDockerComposeCommand()
	//execRawBashCommand()
	execCommandWithSeparateStreams()
}

func execCommandWithSeparateStreams() {
	app := "docker"

	arg0 := "ps"
	arg1 := "-s"
	arg2 := "--format"
	arg3 := "'{\"Names\":\"{{ .Names }}\", \"Size\":\"{{ .Size }}\"}'"

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command: %s. Command Output: %s", err, stderr.String())
	}

	for _, outputLine := range strings.Split(strings.TrimSpace(stdout.String()), "\n") {
		fmt.Println(outputLine)
	}
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
