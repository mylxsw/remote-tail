package command

import (
	"bufio"
	"io"
	"fmt"
	"sync"
	"github.com/mylxsw/remote-tail/console"
	"github.com/mylxsw/remote-tail/ssh"
)

type Command struct {
	Host    string
	User    string
	Script  string
	Stdout  io.Reader
	Stderr  io.Reader
	Server  Server
}

type Message struct {
	Host    string
	Content string
}

func NewCommand(server Server) (cmd *Command) {
	cmd = &Command{
		Host: server.Hostname,
		User: server.User,
		Script: fmt.Sprintf("tail -f %s", server.TailFile),
		Server: server,
	}

	return
}

func (cmd *Command) Execute(output chan Message) {
	client := &ssh.Client{
		Host: cmd.Host,
		User: cmd.User,
		Password: cmd.Server.Password,
	}

	if err := client.Connect(); err != nil {
		panic(fmt.Sprintf("[%s] 连接到服务器失败: %s", cmd.Host, err))
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(fmt.Sprintf("[%s] 创建会话失败: %s", cmd.Host, err))
	}
	defer session.Close()

	cmd.Stdout, err = session.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("[%s] 标准输出重定向失败: %s", cmd.Host, err))
	}

	cmd.Stderr, err = session.StderrPipe()
	if err != nil {
		panic(fmt.Sprintf("[%s] 标准错误输出重定向失败: %s", cmd.Host, err))
	}

	if err = session.Start(cmd.Script); err != nil {
		panic(fmt.Sprintf("[%s] 命令执行失败: %s", cmd.Host, err))
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

	if err = session.Wait(); err != nil {
		panic(fmt.Sprintf("[%s] 命令执行等待失败: %s", cmd.Host, err))
	}

	wg.Wait()
}

func bindOutput(host string, output chan Message, input *io.Reader, prefix string, color int) {
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