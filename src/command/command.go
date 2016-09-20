package command

import (
	"os/exec"
	"bufio"
	"io"
	"log"
)

type Command struct {
	command *exec.Cmd
	Host    string
	Script  string
	Stdout  io.ReadCloser
}

type Message struct {
	Host    string
	Content string
}

func NewCommand(host string, script string) (cmd *Command, err error) {
	cmd = &Command{
		Host: host,
		Script: script,
	}

	cmd.command = exec.Command("ssh", "-f", host, script)
	stdout, err := cmd.command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stdout = stdout


	return
}

func (cmd *Command) Execute(stdout chan Message) {
	if err := cmd.command.Start(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	bindOutput(cmd.Host, &cmd.Stdout, stdout)

	if err := cmd.command.Wait(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func bindOutput(host string, input *io.ReadCloser, output chan Message) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				log.Fatalf("Error: %d", err)
			}
			break
		}

		output <- Message{
			Host: host,
			Content: line,
		}
	}
}