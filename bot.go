package zerobot

import (
	"encoding/json"
	"fmt"

	"github.com/nlopes/slack"
)

type ZeroBot struct {
	defaultChannel slack.Channel
	client         *slack.Client
	rtm            *slack.RTM
}

func NewZeroBot(token string, options ...bool) *ZeroBot {
	var debug bool
	if len(options) > 0 {
		debug = options[0]
	}
	bot := &ZeroBot{}

	client := slack.New(token)
	client.SetDebug(debug)
	bot.client = client
	bot.rtm = client.NewRTM()
	bot.initDefaultChannel()
	return bot
}

// default channel is random
func (b *ZeroBot) initDefaultChannel() {
	channels, err := b.client.GetChannels(true)
	if err != nil {
		Logger.Println(err.Error())
		return
	}

	for _, ch := range channels {
		if ch.Name == "random" {
			b.defaultChannel = ch
			return
		}
	}

	// the first one
	b.defaultChannel = channels[0]
}

func (b *ZeroBot) Close() {
	// nothing to do now.
}

func (b *ZeroBot) Run() {
	go b.rtm.ManageConnection()

	b.sendMsg("zerobot starting ...")

Loop:
	for {
		select {
		case msg := <-b.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				Logger.Println("Infos:", ev.Info)
				Logger.Println("Connection counter:", ev.ConnectionCount)

			case *slack.MessageEvent:
				b.handlerMsg(ev.Msg.Text)

			case *slack.PresenceChangeEvent:
				//

			case *slack.LatencyReport:
				// Logger.Printf("Current latency: %+v\n", ev.Value)

			case *slack.RTMError:
				b.sendMsg(ev.Error(), b.defaultChannel.ID)
				Logger.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				Logger.Println("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
			}
		}
	}
}

func (b *ZeroBot) handlerMsg(msg string) {
	switch msg {
	case "channel":
		channels, err := b.client.GetChannels(true)
		if err == nil {
			b.sendMsg(marshal(channels))
		}

	case "team":
		team, err := b.client.GetTeamInfo()
		if err == nil {
			b.sendMsg(marshal(team))
		}

	case "user":
		users, err := b.client.GetUsers()
		if err == nil {
			b.sendMsg(marshal(users))
		}

	case "ping":
		b.sendMsg("pong")

	case "build":
		b.handleBuild()

	case "restart":
		b.handleRestart()

	case "default":
		b.sendMsg("default channelID: %s", b.defaultChannel.ID)

	default:
		// b.sendMsg("unknown command: %s", msg)
	}

}

func (b *ZeroBot) sendMsg(text string, a ...interface{}) {
	b.send(fmt.Sprintf(text, a...), b.defaultChannel.ID)
}

func (b *ZeroBot) sendErr(err error) {
	Logger.Println("err: ", err.Error())
	b.send(err.Error(), b.defaultChannel.ID)
}

func (b *ZeroBot) send(text, channelID string) {
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(text, channelID))
}

func marshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		Logger.Println(err.Error())
		return ""
	}

	return string(b)
}
