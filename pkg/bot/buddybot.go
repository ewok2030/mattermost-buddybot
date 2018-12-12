package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ewok2030/mattermost-buddybot/pkg/util"
	"github.com/mattermost/mattermost-server/model"
)

// Get configuration settings from environment variables
type BuddyBot struct {
	Client       *model.Client4
	User         *model.User
	DebugChannel *model.Channel
}

func NewBuddyBot(client *model.Client4, user *model.User, channel *model.Channel) *BuddyBot {
	return &BuddyBot{client, user, channel}
}

func (b *BuddyBot) SendMessage(channelId string, msg string, replyToId string) {
	post := &model.Post{}
	post.ChannelId = channelId
	post.Message = msg
	post.RootId = replyToId

	if _, resp := b.Client.CreatePost(post); resp.Error != nil {
		fmt.Printf("Failed to send message to the channel: %s.\n", channelId)
		fmt.Println(resp.Error.Error())
	}
}

func (b *BuddyBot) LogMessage(message string) {

	fmt.Printf("DEBUG\t::\t%s\n", message)

	if b.DebugChannel == nil {
		fmt.Println("channel is nil")
		return
	}

	post := &model.Post{}
	post.ChannelId = b.DebugChannel.Id
	post.Message = message

	if _, resp := b.Client.CreatePost(post); resp.Error != nil {
		fmt.Printf("Failed to send message to the channel: %s.\n", b.DebugChannel.Id)
		fmt.Println(resp.Error.Error())
	}
}

func (b *BuddyBot) HandleMessage(event *model.WebSocketEvent) {

	// Lets only reponded to messaged posted events
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	// Buddybot must be the first mention
	mentions := mattermost.ParseMentions(event.Data)
	if len(mentions) < 1 || mentions[0] != b.User.Id {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {

		// ignore my events
		if post.UserId == b.User.Id {
			return
		}

		// if you see any word matching 'alive' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)alive(?:$|\W)`, post.Message); matched {
			b.SendMessage(post.ChannelId, "Yes, I'm alive.", post.Id)
			return
		}

		// if you see any word matching 'up' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)up(?:$|\W)`, post.Message); matched {
			b.SendMessage(post.ChannelId, "Yes, I'm up and running", post.Id)
			return
		}

		// if you see any word matching 'running' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)running(?:$|\W)`, post.Message); matched {
			b.SendMessage(post.ChannelId, "Yes, I'm running", post.Id)
			return
		}

		// if you see any word matching 'hello' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)hello(?:$|\W)`, post.Message); matched {
			b.SendMessage(post.ChannelId, "Hello!", post.Id)
			return
		}
	}

	b.SendMessage(post.ChannelId, "I did not understand you. :(", post.Id)
}
