package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mpalumbo7/mattermost-buddybot/pkg/util"
)

// BuddyBot is for sending and handling messages
type BuddyBot struct {
	Client *model.Client4
	User   *model.User
}

// SendMessage commands the bot to send a message
func (b *BuddyBot) SendMessage(channelID string, message string, replyToID string) {
	post := &model.Post{}
	post.ChannelId = channelID
	post.Message = message
	post.RootId = replyToID

	if _, resp := b.Client.CreatePost(post); resp.Error != nil {
		fmt.Printf("Failed to send message to the channel: %s.\n", channelID)
		fmt.Println(resp.Error.Error())
	}
}

// HandleMessage is for processing server events
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

		// if you see any word matching 'buddy' then respond
		if matched, _ := regexp.MatchString(`(?:^|\W)buddy(?:$|\W)`, post.Message); matched {

			// Get a random buddy!
			buddy := b.findBuddy(post.ChannelId)

			b.SendMessage(post.ChannelId, fmt.Sprintf("Why don't you try talking to %s?", buddy.Username), post.Id)
			return
		}

	}

	b.SendMessage(post.ChannelId, "I did not understand you. :(", post.Id)
}

func (b *BuddyBot) findBuddy(channelID string) *model.User {

	// Get the channel for this message, to access TeamId
	channel, _ := b.Client.GetChannel(channelID, "")

	// Get the team size
	teamStats, _ := b.Client.GetTeamStats(channel.TeamId, "")

	// Get random page of users
	rand.Seed(time.Now().Unix())
	const pageSize = 10
	pageNum := rand.Intn(int(teamStats.TotalMemberCount) / pageSize)
	teamMembers, _ := b.Client.GetTeamMembers(channel.TeamId, pageNum, pageSize, "")

	// Pick a random buddy!
	buddy, _ := b.Client.GetUser(teamMembers[rand.Intn(len(teamMembers))].UserId, "")

	return buddy
}
