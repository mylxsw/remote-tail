package main

import (
	"fmt"
	"log"
	"sync"
	"flag"
	"strings"
	"os"
	"github.com/BurntSushi/toml"
	"strconv"
	"github.com/mylxsw/remote-tail/console"
	"github.com/mylxsw/remote-tail/command"
)

var mossSep = ".--. --- .-- . .-. . -..   -... -.--   -- -.-- .-.. -..- ... .-- \n"

var welcomeMessage string = `
 ____                      _      _____     _ _
|  _ \ ___ _ __ ___   ___ | |_ __|_   _|_ _(_) |
| |_) / _ \ '_ ' _ \ / _ \| __/ _ \| |/ _' | | |
|  _ <  __/ | | | | | (_) | ||  __/| | (_| | | |
|_| \_\___|_| |_| |_|\___/ \__\___||_|\__,_|_|_|

author: mylxsw
homepage: github.com/mylxsw/remote-tail
version: 0.1.1
` + console.ColorfulText(console.TextMagenta, mossSep)

var filePath *string = flag.String("file", "", "-file=\"/home/data/logs/**/*.log\"")
var hostStr *string = flag.String("hosts", "", "-hosts=root@192.168.1.225,root@192.168.1.226")
var configFile *string = flag.String("conf", "", "-conf=example.toml")

func usageAndExit(message string) {

	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}

	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")

	os.Exit(1)
}

func printWelcomeMessage(config command.Config) {
	fmt.Println(welcomeMessage)

	for _, server := range config.Servers {
		// 如果单独的服务配置没有tail_file,则使用全局配置
		if server.TailFile == "" {
			server.TailFile = config.TailFile
		}

		serverInfo := fmt.Sprintf("%s@%s:%s", server.User, server.Hostname, server.TailFile)
		fmt.Println(console.ColorfulText(console.TextMagenta, serverInfo))
	}
	fmt.Printf("\n%s\n", console.ColorfulText(console.TextCyan, mossSep))
}

func parseConfig(filePath string, hostStr string, configFile string) (config command.Config) {
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			log.Fatal(err)
		}

	} else {

		var hosts []string = strings.Split(hostStr, ",")
		var script string = fmt.Sprintf("tail -f %s", filePath)

		config = command.Config{}
		config.TailFile = script
		config.Servers = make(map[string]command.Server, len(hosts))
		for index, hostname := range hosts {
			hostInfo := strings.Split(strings.Replace(hostname, ":", "@", -1), "@")
			var port int
			if len(hostInfo) > 2 {
				port, _ = strconv.Atoi(hostInfo[2])
			}
			config.Servers["server_" + string(index)] = command.Server{
				ServerName: "server_" + string(index),
				Hostname: hostInfo[1],
				User: hostInfo[0],
				Port: port,
			}
		}
	}

	return
}

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, welcomeMessage)
		fmt.Fprint(os.Stderr, "Options:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if (*filePath == "" || *hostStr == "") && *configFile == "" {
		usageAndExit("")
	}

	config := parseConfig(*filePath, *hostStr, *configFile)
	printWelcomeMessage(config)

	outputs := make(chan command.Message, 20)
	var wg sync.WaitGroup

	for _, server := range config.Servers {
		wg.Add(1)
		go func(server command.Server) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf(console.ColorfulText(console.TextRed, "Error: %s\n"), err)
				}
			}()
			defer wg.Done()

			// 如果单独的服务配置没有tail_file,则使用全局配置
			if server.TailFile == "" {
				server.TailFile = config.TailFile
			}

			// 如果服务配置没有port，则使用默认值22
			if server.Port == 0 {
				server.Port = 22
			}

			cmd, err := command.NewCommand(server)
			if err != nil {
				panic(err)
			}

			cmd.Execute(outputs)
		}(server)
	}

	if len(config.Servers) > 0 {
		go func() {
			for output := range outputs {
				fmt.Printf(
					"%s %s %s",
					console.ColorfulText(console.TextGreen, output.Host),
					console.ColorfulText(console.TextYellow, "->"),
					output.Content,
				)
			}
		}()
	} else {
		fmt.Println(console.ColorfulText(console.TextRed, "没有可用的目标主机"))
	}

	wg.Wait()
}
