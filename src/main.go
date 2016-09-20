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
)

var welcomeMessage string = `
 ____                      _      _____     _ _
|  _ \ ___ _ __ ___   ___ | |_ __|_   _|_ _(_) |
| |_) / _ \ '_ ' _ \ / _ \| __/ _ \| |/ _' | | |
|  _ <  __/ | | | | | (_) | ||  __/| | (_| | | |
|_| \_\___|_| |_| |_|\___/ \__\___||_|\__,_|_|_|

author: mylxsw
homepage: github.com/mylxsw/remote-tail
` + "\x1b[0;31m-----------------------------------------------\x1b[0m\n"

var filepath *string = flag.String("file", "", "-file=\"/home/data/logs/**/*.log\"")
var hostStr *string = flag.String("hosts", "", "-hosts=root@192.168.1.225,root@192.168.1.226")

func usageAndExit(message string) {

	if message != "" {
		fmt.Fprintln(os.Stderr, message)
	}

	flag.Usage()
	fmt.Fprint(os.Stderr, "\n")

	os.Exit(1)
}

func main() {

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, welcomeMessage)
		fmt.Fprint(os.Stderr, "Options:\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *filepath == "" || *hostStr == "" {
		usageAndExit("")
	}

	var hosts []string = strings.Split(*hostStr, ",")
	var script string = fmt.Sprintf("tail -f %s", *filepath)

	fmt.Println(welcomeMessage)
	fmt.Println(console.ColorfulText(console.TextMagenta, script) + "\n")

	outputs := make(chan command.Message, 20)
	var wg sync.WaitGroup

	for _, hostname := range hosts {
		wg.Add(1)
		go func(host, script string) {
			cmd, err := command.NewCommand(host, script)
			if err != nil {
				log.Fatal(err)
			}

			cmd.Execute(outputs)

			wg.Done()
		}(hostname, script)
	}

	if len(hosts) > 0 {
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
