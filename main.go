package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/ewok2030/mattermost-buddybot/pkg/bot"
	"github.com/mattermost/mattermost-server/model"
)

// Get configuration settings from environment variables
var (
	ServerURL, _     = url.Parse(os.Getenv("SERVER_URL"))
	BotUsername      = os.Getenv("BOT_USERNAME")
	BotPassword      = os.Getenv("BOT_PASSWORD")
	DebugChannel     = os.Getenv("DEBUG_CHANNEL")
	DebugChannelTeam = os.Getenv("DEBUG_CHANNEL_TEAM")
)

var client *model.Client4
var webSocketClient *model.WebSocketClient
var debugChannel *model.Channel

// Documentation for the Go driver can be found
// at https://godoc.org/github.com/mattermost/platform/model#Client
func main() {

	fmt.Printf("Starting up buddybot '%s' for %s\n", BotUsername, ServerURL)
	client = model.NewAPIv4Client(ServerURL.String())

	// Confirm server is running
	if props, resp := client.GetOldClientConfig(""); resp.Error != nil {
		fmt.Printf("There was a problem reaching the Mattermost server '%s'.\n  Exiting.\n\n", ServerURL)
		fmt.Println(resp.Error.Error())
		os.Exit(1)
	} else {
		fmt.Printf("Server '%s' was detected and is running version %s\n", ServerURL, props["Version"])
	}

	// Login as the bot
	botUser, resp := client.Login(BotUsername, BotPassword)
	if resp.Error != nil {
		println("There was an error logging into: %s. Exiting", ServerURL)
		fmt.Println(resp.Error.Error())
		os.Exit(1)
	}

	// Get the Channel to use for debugging (if defined)
	if DebugChannel != "" && DebugChannelTeam != "" {
		// Get the Team
		botTeam, resp := client.GetTeamByName(DebugChannelTeam, "")
		if resp.Error != nil {
			fmt.Printf("Failed to find team '%s', or '%s' may not be member of the team. Exiting", DebugChannelTeam, BotUsername)
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}

		debugChannel, resp = client.GetChannelByName(DebugChannel, botTeam.Id, "")
		if resp.Error != nil {
			println("Failed to find debugging channel '%s'. Exiting.", DebugChannel)
			fmt.Println(resp.Error.Error())
			os.Exit(1)
		}
	}

	// Register a listener to capture CTRL+C commands
	setupGracefulShutdown()

	// Create WebSocket client
	// Lets start listening to some channels via the websocket!
	websocketURL := ServerURL
	websocketURL.Scheme = "wss"
	if ServerURL.Scheme == "http" {
		websocketURL.Scheme = "ws"
	}
	webSocketClient, err := model.NewWebSocketClient4(websocketURL.String(), client.AuthToken)
	if err != nil {
		fmt.Println("Filed to connect to the web socket")
		fmt.Println(err.Error())
	}
	webSocketClient.Listen()

	// Initialize the bot!
	myBot := &bot.BuddyBot{Client: client, User: botUser}
	sendDebugMessage("_buddybot has **STARTED** running!_")

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				myBot.HandleMessage(resp)
			}
		}
	}()

	// You can block forever with
	select {}
}

func sendDebugMessage(message string) {
	fmt.Printf("\nDEBUG::\t%s\n", message)

	if debugChannel == nil {
		return
	}

	post := &model.Post{}
	post.ChannelId = debugChannel.Id
	post.Message = message

	if _, resp := client.CreatePost(post); resp.Error != nil {
		fmt.Printf("Failed to send message to the channel: %s.\n", debugChannel.Id)
		fmt.Println(resp.Error.Error())
	}
}

func setupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if webSocketClient != nil {
			webSocketClient.Close()
		}

		sendDebugMessage("_buddybot has **STOPPED** running!_")
		os.Exit(0)

	}()
}
