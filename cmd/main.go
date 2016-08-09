package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/holys/zerobot"
	"github.com/nlopes/slack"
)

var (
	debug = flag.Bool("debug", false, "debug or not")
)

func main() {
	flag.Parse()
	token := os.Getenv("SLACK_API_TOKEN")
	if strings.TrimSpace(token) == "" {
		fmt.Println("please set $SLACK_API_TOKEN")
		return
	}

	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	zerobot.Logger = logger
	slack.SetLogger(logger)

	bot := zerobot.NewZeroBot(token, *debug)
	bot.Run()

	// TODO: handle Signal
}
