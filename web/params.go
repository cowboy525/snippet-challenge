package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	pageDefault     = 1
	pageSizeDefault = 20
	pageSizeMaximum = 1000
)

// Params : query params structure
type Params struct {
	UserID      uint64
	SpaceID     uint64
	ProjectID   uint64
	ChatID      uint64
	TaskID      uint64
	MediaID     string
	MediaTreeID uint64
	NoteID      uint64
	CommentID   uint64
	TaskMemoID  uint64
	ReactionID  uint64
	SharedToken string
	FilterID    uint64

	Page  int
	Limit int

	Account     string
	DisplayName string
	IDs         string // getRecents, getTasks
	QUserID     uint64
	QProjectID  uint64
	QNoteID     uint64
	QTaskID     uint64

	Table            string // getRecents
	Pk               string // getRecents
	SubKey           string // getRecents
	Event            string // getRecents
	UpdatedField     string // getRecents
	DeferBeforeValue bool   // getRecents
	DeferFields      string // getRecents

	TreeNodes    string
	UnjoinedTree uint64
	Favorite     bool // media, task
	Name         string
	Attached     bool

	Root          *bool  // getTasks, getNotes
	ParentID      uint64 // getTasks, getNotes
	Subject       string // getTasks, getNotes
	Statuses      string // getTasks, getNotes
	BatonUserIDs  string // getTasks, getNotes
	WriteUserID   uint64 // getTasks, getNotes
	ChargeUserIDs string // getTasks, getNotes
	Tags          string // getTasks, getNotes
	Ordering      string // getTasks, getNotes

	CreatedAt            *time.Time // getTaskMemos
	CreatedAtGt          *time.Time // getTaskMemos
	CreatedAtGte         *time.Time // getRecents
	CreatedAtLt          *time.Time // getTaskMemos
	CreatedAtLte         *time.Time // getNoteComments
	CheckedAtGte         *time.Time // getNotices
	ReadAtGte            *time.Time // getNotices
	UpdatedAtGt          *time.Time // getTasks
	UpdatedAtGte         *time.Time
	UpdatedAtLte         *time.Time
	HierarchyUpdatedAtGt *time.Time // getTasks
	StartAt              *time.Time // searchTasks
	EndAt                *time.Time // searchTasks

	SubjectIContains string

	Checked   *bool
	Important *bool
	Watch     *bool

	SpecialPattern string
}

func GetUintFromProps(props map[string]string, key string) uint64 {
	if val, ok := props[key]; ok {
		if v, err := strconv.ParseUint(val, 10, 64); err == nil {
			return v
		}
	}
	return 0
}

func GetStringFromProps(props map[string]string, key string) string {
	if val, ok := props[key]; ok {
		return val
	}
	return ""
}

// ParamsFromRequest : get query params from request
func ParamsFromRequest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	query := r.URL.Query()

	params.UserID = GetUintFromProps(props, "user_id")
	params.SpaceID = GetUintFromProps(props, "space_id")
	params.ProjectID = GetUintFromProps(props, "project_id")
	params.ChatID = GetUintFromProps(props, "chat_id")
	params.TaskID = GetUintFromProps(props, "task_id")
	params.TaskMemoID = GetUintFromProps(props, "task_memo_id")
	params.NoteID = GetUintFromProps(props, "note_id")
	params.MediaTreeID = GetUintFromProps(props, "media_tree_id")
	params.CommentID = GetUintFromProps(props, "comment_id")
	params.ReactionID = GetUintFromProps(props, "reaction_id")
	params.SharedToken = GetStringFromProps(props, "shared_token")
	params.FilterID = GetUintFromProps(props, "filter_id")

	if val, ok := props["media_id"]; ok {
		params.MediaID = val
	}

	if val, err := strconv.Atoi(query.Get("page")); err != nil || val < 0 {
		params.Page = pageDefault
	} else {
		params.Page = val
	}

	if val, err := strconv.Atoi(query.Get("limit")); err != nil || val < 0 {
		params.Limit = pageSizeDefault
	} else if val > pageSizeMaximum {
		params.Limit = pageSizeMaximum
	} else {
		params.Limit = val
	}

	if val, err := strconv.Atoi(query.Get("user_id")); err == nil {
		params.QUserID = uint64(val)
	}

	if val, err := strconv.Atoi(query.Get("project_id")); err == nil {
		params.QProjectID = uint64(val)
	}
	if val, err := strconv.Atoi(query.Get("note_id")); err == nil {
		params.QNoteID = uint64(val)
	}

	if val, err := strconv.Atoi(query.Get("task_id")); err == nil {
		params.QTaskID = uint64(val)
	}

	params.Account = query.Get("account")
	params.DisplayName = query.Get("display_name")
	params.IDs = query.Get("id")                     // getRecents, getTasks
	params.Table = query.Get("table")                // getRecents
	params.Pk = query.Get("pk")                      // getRecents
	params.SubKey = query.Get("sub_key")             // getRecents
	params.Event = query.Get("event")                // getRecents
	params.UpdatedField = query.Get("updated_field") // getRecents
	params.DeferFields = query.Get("defer_fields")   // getRecents

	// getRecents
	if val, err := strconv.ParseBool(query.Get("defer_before_value")); err == nil {
		params.DeferBeforeValue = val
	}

	params.TreeNodes = query.Get("tree_nodes")
	params.Name = query.Get("name")

	if val, err := strconv.Atoi(query.Get("unjoined_tree")); err == nil {
		params.UnjoinedTree = uint64(val)
	}

	// getTasks, media
	if val, err := strconv.ParseBool(query.Get("favorite")); err == nil {
		params.Favorite = val
	}

	if val, err := strconv.ParseBool(query.Get("attached")); err == nil {
		params.Attached = val
	}

	// getTasks, mediaTree
	if val, err := strconv.Atoi(query.Get("parent")); err == nil {
		params.ParentID = uint64(val)
	}

	// getTasks
	if val, err := strconv.ParseBool(query.Get("root")); err == nil {
		params.Root = &val
	}

	// getTasks
	params.Subject = query.Get("subject")
	params.Statuses = query.Get("status")
	params.BatonUserIDs = query.Get("baton_user")

	// getTasks
	if val, err := strconv.Atoi(query.Get("write_user")); err == nil {
		params.WriteUserID = uint64(val)
	}

	// getTasks
	params.ChargeUserIDs = query.Get("charge_users")
	params.Tags = query.Get("tags")
	params.Ordering = query.Get("ordering")

	params.SubjectIContains = query.Get("subject__icontains")

	// getNotices
	if val, err := strconv.ParseBool(query.Get("important")); err == nil {
		params.Important = &val
	}
	if val, err := strconv.ParseBool(query.Get("watch")); err == nil {
		params.Watch = &val
	}
	if val, err := strconv.ParseBool(query.Get("checked")); err == nil {
		params.Checked = &val
	}

	if val, err := time.Parse(time.RFC3339, query.Get("created_at")); err == nil {
		params.CreatedAt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("created_at__gt")); err == nil {
		params.CreatedAtGt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("created_at__lt")); err == nil {
		params.CreatedAtLt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("created_at__gte")); err == nil {
		params.CreatedAtGte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("created_at__lte")); err == nil {
		params.CreatedAtLte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("checked_at__gte")); err == nil {
		params.CheckedAtGte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("read_at__gte")); err == nil {
		params.ReadAtGte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("updated_at__gt")); err == nil {
		params.UpdatedAtGt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("updated_at__gte")); err == nil {
		params.UpdatedAtGte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("updated_at__lte")); err == nil {
		params.UpdatedAtLte = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("hierarchy_updated_at__gt")); err == nil {
		params.HierarchyUpdatedAtGt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("start_at")); err == nil {
		params.StartAt = &val
	}
	if val, err := time.Parse(time.RFC3339, query.Get("end_at")); err == nil {
		params.EndAt = &val
	}

	params.SpecialPattern = query.Get("special_pattern")

	return params
}
