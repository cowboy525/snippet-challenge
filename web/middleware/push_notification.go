package middleware

import (
	"fmt"
	"os"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/topoface/onesignal-go"
)

func SendPushNotificaton(a app.Iface) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in SendPushNotificaton", r)
		}
	}()
	for _, event := range a.Session().SyncEvents {
		notificationData, ok := event.Context["mobile_notification"].(*model.MobileNotificationData)
		if !ok {
			continue
		}

		playerIDs := []string{}
		for _, uid := range notificationData.Users {
			if tokens, err := a.Store().Session().GetDeviceTokens(uid); err == nil {
				playerIDs = append(playerIDs, tokens...)
			}
		}
		if len(playerIDs) == 0 {
			continue
		}

		go sendPushNotification(a, playerIDs, notificationData)
	}
}

func sendPushNotification(a app.Iface, playerIDs []string, data *model.MobileNotificationData) {
	if len(playerIDs) == 0 {
		return
	}

	appID := os.Getenv("ONESIGNAL_APP_ID")
	appKey := os.Getenv("ONESIGNAL_API_KEY")
	userKey := os.Getenv("ONESIGNAL_USER_AUTH_KEY")
	client := onesignal.NewClient(nil)
	client.AppKey = appKey
	client.UserKey = userKey

	notificationReq := &onesignal.NotificationRequest{
		AppID:              appID,
		Headings:           map[string]string{"en": data.Headings, "ja": data.Headings},
		Contents:           data.Content,
		IsIOS:              true,
		IsAndroid:          true,
		IncludePlayerIDs:   playerIDs,
		Data:               data.AdditionalData,
		URL:                data.LaunchURL,
		IOSSound:           "DeepNotification.wav",
		IOSBadgeType:       "Increase",
		IOSBadgeCount:      1,
		AndroidChannelID:   "f71d8823-f1fa-469d-bbae-635ec37fcf74",
		SmallIcon:          "splash",
		AndroidAccentColor: "FFFFFFFF",
	}

	createRes, res, err := client.Notifications.Create(notificationReq)
	if err != nil {
		fmt.Printf("--- res:%+v, err:%+v\n", res, err.Error())
		return
	}

	if createRes != nil {
		if errors, ok := createRes.Errors.(map[string]interface{}); ok {
			if invalidIDs, ok := errors["invalid_player_ids"].([]interface{}); ok {
				for _, id := range invalidIDs {
					if val, ok := id.(string); ok {
						a.RemoveDeviceTokenFromSessions(val)
					}
				}
			}
		}
	}
}
