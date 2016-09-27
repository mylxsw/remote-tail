package command

import (
	"os/exec"
	"bufio"
	"io"
	"fmt"
	"sync"
	"strconv"
	"github.com/mylxsw/remote-tail/console"
)

type Command struct {
	command *exec.Cmd
	Host    string
	Script  string
	Stdout  io.ReadCloser
	Stderr  io.ReadCloser
	Server  Server
}

type Message struct {
	Host    string
	Content string
}

func NewCommand(server Server) (cmd *Command, err error) {
	cmd = &Command{
		Host: fmt.Sprintf("%s@%s", server.User, server.Hostname),
		Script: fmt.Sprintf("tail -f %s", server.TailFile),
		Server: server,
	}

	commandParameters := []string{}
	// 如果提供了密码,则使用明文密码, 不安全!
	command := "ssh"
	if server.Password != "" {
		command = "sshpass"
		commandParameters = append(commandParameters, "-p", server.Password, "ssh")
	}
	commandParameters = append(commandParameters, "-p", strconv.Itoa(server.Port), "-f", cmd.Host, cmd.Script)

	cmd.command = exec.Command(command, commandParameters...)

	stdout, err := cmd.command.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stdout = stdout

	stderr, err := cmd.command.StderrPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = stderr

	return
}

func (cmd *Command) Execute(output chan Message) {
	if err := cmd.command.Start(); err != nil {
		panic(fmt.Sprintf("[%s] 命令启动失败: %s", cmd.Host, err))
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func () {
		defer wg.Done()
		bindOutput(cmd.Host, output, &cmd.Stdout, "", 0)
	}()
	go func () {
		defer wg.Done()
		bindOutput(cmd.Host, output, &cmd.Stderr, "Error:", console.TextRed)
	}()

	wg.Wait()

	if err := cmd.command.Wait(); err != nil {
		panic(fmt.Sprintf("[%s] 等待命令执行失败: %s", cmd.Host, err))
	}
}

func bindOutput(host string, output chan Message, input *io.ReadCloser, prefix string, color int) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				panic(fmt.Sprintf("[%s] 命令执行失败: %s", host, err))
			}
			break
		}

		line = prefix + line
		if color != 0 {
			line = console.ColorfulText(color, line)
		}

		output <- Message{
			Host: host,
			Content: line,
		}
	}
}