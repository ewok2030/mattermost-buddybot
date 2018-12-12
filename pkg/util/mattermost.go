package mattermost

import (
	"encoding/json"
	"strings"
)

func ParseMentions(data map[string]interface{}) []string {
	var mentions []string
	if val, ok := data["mentions"]; ok {
		// Decode JSON
		json.NewDecoder(strings.NewReader(val.(string))).Decode(&mentions)
	}
	return mentions
}
