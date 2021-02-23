package middleware

import (
	"fmt"
	"strconv"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/types"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
)

func Notification(a app.Iface) {
	for _, event := range a.Session().SyncEvents {
		notification(a, event)
	}
}

func generate(recent *model.Recent, uid uint64, reason string, important bool, watch bool) *model.Notice {
	return &model.Notice{
		ProjectID: recent.ProjectID,
		RecentID:  recent.ID,
		Table:     recent.Table,
		Pk:        model.NewString(recent.Pk),
		SubKey:    recent.SubKey,
		UserID:    uid,
		Reason:    reason,
		Important: types.NewIntBool(important),
		Watch:     types.NewIntBool(watch),
	}
}

func userExistsInInstances(instances []*model.Notice, userId uint64, importantOnly bool) bool {
	for i := range instances {
		if instances[i].UserID == userId {
			if !importantOnly {
				return true
			}
			if instances[i].Important != nil && *instances[i].Important {
				return true
			}
		}
	}
	return false
}

func notification(a app.Iface, syncEvent *model.SyncEvent) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered in notification", r)
		}
	}()

	if syncEvent.Recent == nil {
		return
	}
	if a.Session().Project != nil && a.Session().Project.Archive != nil && *a.Session().Project.Archive {
		return
	}

	event := syncEvent.EventData.Event
	data := syncEvent.EventData.Data
	recent := syncEvent.Recent
	recentData := syncEvent.RecentData
	oldValue := syncEvent.OldValue
	newValue := syncEvent.NewValue

	instances := []*model.Notice{}

	// note comment
	if event == model.CREATE_NOTE_COMMENT_EVENT {
		if syncEvent.Context != nil {
			// mention
			if mentions, ok := syncEvent.Context["mentions"].(map[uint64]int); ok {
				for uid := range mentions {
					instances = append(instances, generate(recent, uid, "mention", true, true))
				}
			}
			// new comment
			if oldNoteValue, ok := syncEvent.Context["old_note_value"].(map[string]interface{}); ok {
				// baton_user
				if userId, ok := oldNoteValue["baton_user"].(*uint64); ok && userId != nil {
					if !userExistsInInstances(instances, *userId, true) {
						instances = append(instances, generate(recent, *userId, "new_comment", true, false))
					}
				}
				// following_users
				note, _ := a.Store().Note().Get(data["params"].(map[string]interface{})["noteId"].(uint64))
				for _, userId := range note.FollowingUserIDs {
					if !userExistsInInstances(instances, userId, true) {
						instances = append(instances, generate(recent, userId, "new_comment", true, false))
					}
				}
				// charge_users
				if chargeUsers, ok := oldNoteValue["charge_users"].([]uint64); ok {
					for _, userId := range chargeUsers {
						if !userExistsInInstances(instances, userId, false) {
							instances = append(instances, generate(recent, userId, "new_comment", false, false))
						}
					}
				}
			}
		}
	} else
	// note
	if event == model.CREATE_NOTE_EVENT || event == model.UPDATE_NOTE_EVENT {
		for _, rd := range recentData {
			// baton_user
			if rd.Field == "baton_user" && rd.AfterValue != nil {
				if uid, err := strconv.ParseUint(*rd.AfterValue, 10, 64); err == nil {
					instances = append(instances, generate(recent, uid, "baton", true, false))
				}
			}
			// charge_users
			if rd.Field == "charge_users" {
				oldChargeUsers := make([]uint64, 0)
				if oldValue != nil {
					if val, ok := oldValue["charge_users"].([]uint64); ok {
						oldChargeUsers = val
					}
				}
				newChargeUsers := make([]uint64, 0)
				if newValue != nil {
					if val, ok := newValue["charge_users"].([]uint64); ok {
						newChargeUsers = val
					}
				}
				leftUsers := utils.Difference(oldChargeUsers, newChargeUsers)
				for _, userId := range leftUsers {
					instances = append(instances, generate(recent, userId, "leave", true, false))
				}
				joinedUsers := utils.Difference(newChargeUsers, oldChargeUsers)
				for _, userId := range joinedUsers {
					if val, ok := newValue["baton_user"].(*uint64); !ok || val == nil || userId != *val {
						instances = append(instances, generate(recent, userId, "join", true, false))
					}
				}
			}
		}
		// update
		if oldValue != nil {
			if userId, ok := oldValue["baton_user"].(*uint64); ok && userId != nil {
				if !userExistsInInstances(instances, *userId, true) {
					instances = append(instances, generate(recent, *userId, "update", true, false))
				}
			}
		}
		if chargeUsers, ok := oldValue["charge_users"].([]uint64); ok {
			for _, userId := range chargeUsers {
				if !userExistsInInstances(instances, userId, false) {
					instances = append(instances, generate(recent, userId, "update", false, false))
				}
			}
		}
	} else
	// task comment
	if event == model.CREATE_TASK_COMMENT_EVENT {
		if syncEvent.Context != nil {
			// mention
			if mentions, ok := syncEvent.Context["mentions"].(map[uint64]int); ok {
				for uid := range mentions {
					instances = append(instances, generate(recent, uid, "mention", true, true))
				}
			}
			// new comment
			if oldTaskValue, ok := syncEvent.Context["old_task_value"].(map[string]interface{}); ok {
				// baton_user
				if userId, ok := oldTaskValue["baton_user"].(*uint64); ok && userId != nil {
					if !userExistsInInstances(instances, *userId, true) {
						instances = append(instances, generate(recent, *userId, "new_comment", true, false))
					}
				}
				// following_users
				task, _ := a.Store().Task().Get(data["params"].(map[string]interface{})["taskId"].(uint64))
				followingUsers, _ := a.Store().Task().GetFollowingUsers(task)
				for _, userId := range followingUsers {
					if !userExistsInInstances(instances, userId, true) {
						instances = append(instances, generate(recent, userId, "new_comment", true, false))
					}
				}
				// charge_users
				if chargeUsers, ok := oldTaskValue["charge_users"].([]uint64); ok {
					for _, userId := range chargeUsers {
						if !userExistsInInstances(instances, userId, false) {
							instances = append(instances, generate(recent, userId, "new_comment", false, false))
						}
					}
				}
			}
		}
	} else
	// task memo
	if event == model.CREATE_TASK_MEMO_EVENT || event == model.UPDATE_TASK_MEMO_EVENT {
		if syncEvent.Context != nil {
			// create_memo or update_memo
			if oldTaskValue, ok := syncEvent.Context["old_task_value"].(map[string]interface{}); ok {
				reason := "update_memo"
				if event == model.CREATE_TASK_MEMO_EVENT {
					reason = "create_memo"
				}
				// baton_user
				if userId, ok := oldTaskValue["baton_user"].(*uint64); ok && userId != nil {
					if !userExistsInInstances(instances, *userId, true) {
						instances = append(instances, generate(recent, *userId, reason, true, false))
					}
				}
				// charge_users
				if chargeUsers, ok := oldTaskValue["charge_users"].([]uint64); ok {
					for _, userId := range chargeUsers {
						if !userExistsInInstances(instances, userId, false) {
							instances = append(instances, generate(recent, userId, reason, false, false))
						}
					}
				}
			}
		}
	} else
	// task
	if event == model.CREATE_TASK_EVENT || event == model.UPDATE_TASK_EVENT {
		for _, rd := range recentData {
			// baton_user
			if rd.Field == "baton_user" && rd.AfterValue != nil {
				if uid, err := strconv.ParseUint(*rd.AfterValue, 10, 64); err == nil {
					instances = append(instances, generate(recent, uid, "baton", true, false))
				}
			}
			// charge_users
			if rd.Field == "charge_users" {
				oldChargeUsers := make([]uint64, 0)
				if oldValue != nil {
					if val, ok := oldValue["charge_users"].([]uint64); ok {
						oldChargeUsers = val
					}
				}
				newChargeUsers := make([]uint64, 0)
				if newValue != nil {
					if val, ok := newValue["charge_users"].([]uint64); ok {
						newChargeUsers = val
					}
				}
				leftUsers := utils.Difference(oldChargeUsers, newChargeUsers)
				for _, userId := range leftUsers {
					instances = append(instances, generate(recent, userId, "leave", true, false))
				}
				joinedUsers := utils.Difference(newChargeUsers, oldChargeUsers)
				for _, userId := range joinedUsers {
					if val, ok := newValue["baton_user"].(*uint64); !ok || val == nil || userId != *val {
						instances = append(instances, generate(recent, userId, "join", true, false))
					}
				}
			}
		}
		// update
		if oldValue != nil {
			if userId, ok := oldValue["baton_user"].(*uint64); ok && userId != nil {
				if !userExistsInInstances(instances, *userId, true) {
					instances = append(instances, generate(recent, *userId, "update", true, false))
				}
			}
		}
		if chargeUsers, ok := oldValue["charge_users"].([]uint64); ok {
			for _, userId := range chargeUsers {
				if !userExistsInInstances(instances, userId, false) {
					instances = append(instances, generate(recent, userId, "update", false, false))
				}
			}
		}
	} else
	// chat comment
	if event == model.CREATE_CHAT_COMMENT_EVENT {
		if syncEvent.Context != nil {
			// mention
			if mentions, ok := syncEvent.Context["mentions"].(map[uint64]int); ok {
				for uid := range mentions {
					instances = append(instances, generate(recent, uid, "mention", true, true))
				}
			}
			// new comment
			if oldChatValue, ok := syncEvent.Context["old_chat_value"].(map[string]interface{}); ok {
				// following_users
				chat, _ := a.Store().Chat().Get(data["params"].(map[string]interface{})["chatId"].(uint64))
				followingUsers, _ := a.Store().Chat().GetFollowingUsers(chat)
				for _, userId := range followingUsers {
					if !userExistsInInstances(instances, userId, true) {
						instances = append(instances, generate(recent, userId, "new_comment", true, false))
					}
				}
				// charge_users
				if chargeUsers, ok := oldChatValue["charge_users"].([]uint64); ok {
					for _, userId := range chargeUsers {
						if !userExistsInInstances(instances, userId, false) {
							instances = append(instances, generate(recent, userId, "new_comment", false, false))
						}
					}
				}
			}
		}
	}

	if len(instances) > 0 {
		// save and send to user
		toUser := data["toUser"].(map[uint64]interface{})
		operateUserId := data["userId"]
		for _, instance := range instances {
			// excluding operator
			if instance.UserID == operateUserId {
				continue
			}
			if !a.UserHasPermissionToProject(instance.UserID, instance.ProjectID) {
				continue
			}
			a.Store().Notice().Create(instance)
			if _, ok := toUser[instance.UserID]; !ok {
				toUser[instance.UserID] = map[string]interface{}{}
			}
			if _, ok := toUser[instance.UserID].(map[string]interface{})["notices"]; !ok {
				toUser[instance.UserID].(map[string]interface{})["notices"] = make([]*model.UncheckedNotice, 0)
			}
			toUser[instance.UserID].(map[string]interface{})["notices"] = append(toUser[instance.UserID].(map[string]interface{})["notices"].([]*model.UncheckedNotice), instance.ToUncheckedNotice())
		}
	}
}
