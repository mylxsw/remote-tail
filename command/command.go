package command

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/mylxsw/remote-tail/console"
	"github.com/mylxsw/remote-tail/ssh"
)

type Command struct {
	Host   string
	User   string
	Script string
	Stdout io.Reader
	Stderr io.Reader
	Server Server
}

// Message The message used by channel to transport log line by line
type Message struct {
	Host    string
	Content string
}

// NewCommand Create a new command
func NewCommand(server Server) (cmd *Command) {
	cmd = &Command{
		Host:   server.Hostname,
		User:   server.User,
		Script: fmt.Sprintf("tail %s %s", server.TailFlags, server.TailFile),
		Server: server,
	}

	if !strings.Contains(cmd.Host, ":") {
		cmd.Host = cmd.Host + ":" + strconv.Itoa(server.Port)
	}

	return
}

// Execute the remote command
func (cmd *Command) Execute(output chan Message) {

	client := &ssh.Client{
		Host:           cmd.Host,
		User:           cmd.User,
		Password:       cmd.Server.Password,
		PrivateKeyPath: cmd.Server.PrivateKeyPath,
	}

	if err := client.Connect(); err != nil {
		panic(fmt.Sprintf("[%s] unable to connect: %s", cmd.Host, err))
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(fmt.Sprintf("[%s] unable to create session: %s", cmd.Host, err))
	}
	defer session.Close()

	if err := session.RequestPty("xterm", 80, 40, *ssh.CreateTerminalModes()); err != nil {
		panic(fmt.Sprintf("[%s] unable to create pty: %v", cmd.Host, err))
	}

	cmd.Stdout, err = session.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("[%s] redirect stdout failed: %s", cmd.Host, err))
	}

	cmd.Stderr, err = session.StderrPipe()
	if err != nil {
		panic(fmt.Sprintf("[%s] redirect stderr failed: %s", cmd.Host, err))
	}

	go bindOutput(cmd.Host, output, &cmd.Stdout, "", 0)
	go bindOutput(cmd.Host, output, &cmd.Stderr, "Error:", console.TextRed)

	if err = session.Start(cmd.Script); err != nil {
		panic(fmt.Sprintf("[%s] failed to execute command: %s", cmd.Host, err))
	}

	if err = session.Wait(); err != nil {
		panic(fmt.Sprintf("[%s] failed to wait command: %s", cmd.Host, err))
	}
}

// bing the pipe output for formatted output to channel
func bindOutput(host string, output chan Message, input *io.Reader, prefix string, color int) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				panic(fmt.Sprintf("[%s] faield to execute command: %s", host, err))
			}
			break
		}

		line = prefix + line
		if color != 0 {
			line = console.ColorfulText(color, line)
		}

		output <- Message{
			Host:    host,
			Content: line,
		}
	}
}
