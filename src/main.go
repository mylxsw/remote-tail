package main

import (
	"fmt"
	"log"
	"sync"
	"flag"
	"strings"
	"os"
	"console"
	"command"
	"github.com/BurntSushi/toml"
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
version: 0.1
` + console.ColorfulText(console.TextMagenta, mossSep)

var filepath *string = flag.String("file", "", "-file=\"/home/data/logs/**/*.log\"")
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

func parseConfig(filepath string, hostStr string, configFile string) (config command.Config) {
	if configFile != "" {
		if _, err := toml.DecodeFile(configFile, &config); err != nil {
			log.Fatal(err)
		}

	} else {

		var hosts []string = strings.Split(hostStr, ",")
		var script string = fmt.Sprintf("tail -f %s", filepath)

		config = command.Config{}
		config.TailFile = script
		config.Servers = make(map[string]command.Server, len(hosts))
		for index, hostname := range hosts {
			hostInfo := strings.Split(hostname, "@")
			config.Servers["server_" + string(index)] = command.Server{
				ServerName: "server_" + string(index),
				Hostname: hostInfo[1],
				User: hostInfo[0],
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

	if (*filepath == "" || *hostStr == "") && *configFile == "" {
		usageAndExit("")
	}

	config := parseConfig(*filepath, *hostStr, *configFile)
	printWelcomeMessage(config)

	outputs := make(chan command.Message, 20)
	var wg sync.WaitGroup

	for _, server := range config.Servers {
		wg.Add(1)
		go func(server command.Server) {

			// 如果单独的服务配置没有tail_file,则使用全局配置
			if server.TailFile == "" {
				server.TailFile = config.TailFile
			}

			cmd, err := command.NewCommand(server)
			if err != nil {
				log.Fatal(err)
			}

			cmd.Execute(outputs)

			wg.Done()
		}(server)
	}

	if len(config.Servers) > 0 {
		wg.Add(1)
		go func() {
			for output := range outputs {
				fmt.Printf(
					"%s %s %s",
					console.ColorfulText(console.TextGreen, output.Host),
					console.ColorfulText(console.TextYellow, "->"),
					output.Content,
				)
			}

			wg.Done()
		}()

		wg.Wait()
	} else {
		log.Fatal("没有可用的目标主机")
	}
}
