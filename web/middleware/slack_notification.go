package middleware

import (
	"fmt"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
	"github.com/slack-go/slack"
)

func SlackNotification(a app.Iface) {
	for _, event := range a.Session().SyncEvents {
		notify(a, event)
	}
}

func notify(a app.Iface, event *model.SyncEvent) {
	users, ok := event.Context["slack_notification_users"].([]uint64)
	if !ok {
		return
	}
	message, ok := event.Context["slack_notification_message"].(string)
	if !ok {
		return
	}
	mention := createMentionString(a, users)
	projectID := utils.GetUintFromInterface(event.EventData.Data["params"].(map[string]interface{})["projectId"])
	go sendSlackNotification(a, message+"\n"+mention, projectID)
}

func createMentionString(a app.Iface, users []uint64) string {
	mention := ""
	for _, id := range users {
		slackUser, _ := a.Store().SlackUser().GetBySpaceUser(id)
		if slackUser != nil {
			mention = fmt.Sprintf("%v<@%v> ", mention, slackUser.UserID)
		}
	}
	return mention
}

func sendSlackNotification(a app.Iface, message string, projectID uint64) {
	channel, _ := a.Store().SlackChannel().GetByProject(projectID)
	if channel == nil {
		fmt.Println("No slack channel is associated with this project: ", projectID)
		return
	}

	team, _ := a.Store().SlackTeam().Get(channel.SlackTeamID)
	if team == nil {
		fmt.Println("No space is associated with this slack workspace: ", channel.SlackTeamID)
		return
	}

	api := slack.New(team.AccessToken)
	if _, _, err := api.PostMessage(channel.ChannelID, slack.MsgOptionPost(), slack.MsgOptionText(message, false)); err != nil {
		fmt.Println("An error occured while sending slack notification: ", err.Error())
	}
}
