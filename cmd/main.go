package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sc
		logger.Printf("signal called, [%d] to exit\n", sig)
		bot.Close()
		os.Exit(0)
	}()

	bot.Run()
}
