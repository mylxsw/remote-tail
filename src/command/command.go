package command

import (
	"os/exec"
	"bufio"
	"io"
	"log"
	"fmt"
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
	commandParameters = append(commandParameters, "-f", cmd.Host, cmd.Script)

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

func (cmd *Command) Execute(stdout chan Message) {
	if err := cmd.command.Start(); err != nil {
		log.Fatalf("[%s] 命令启动失败: %s", cmd.Host, err)
	}

	bindOutput(cmd.Host, &cmd.Stdout, stdout)

	// TODO 处理标准错误输出

	if err := cmd.command.Wait(); err != nil {
		log.Fatalf("[%s] 等待命令执行失败: %s", cmd.Host, err)
	}
}

func bindOutput(host string, input *io.ReadCloser, output chan Message) {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				log.Fatalf("[%s] 命令执行失败: %d", host, err)
			}
			break
		}

		output <- Message{
			Host: host,
			Content: line,
		}
	}
}