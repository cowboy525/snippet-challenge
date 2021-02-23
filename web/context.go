package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
)

// Context structure
type Context struct {
	App           app.Iface
	Log           *mlog.Logger
	Params        *Params
	Err           *model.AppError
	siteURLHeader string
}

func (c *Context) LogError(err *model.AppError) {
	// Filter out 404s, endless reconnects and browser compatibility errors
	var isDebug bool
	isDebug = false
	if err.StatusCode == http.StatusNotFound {
		isDebug = true
	}
	if detail, ok := c.Err.Errors[0]["detail"]; ok {
		if errArr, ok := detail.([]*model.Error); ok && errArr[0].ID == "web.check_browser_compatibility.app_error" {
			isDebug = true
		}
	}
	if isDebug {
		c.LogDebug(err)
	} else {
		c.Log.Error(
			err.SystemMessage(utils.TDefault),
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogInfo(err *model.AppError) {
	// Filter out 401s
	if err.StatusCode == http.StatusUnauthorized {
		c.LogDebug(err)
	} else {
		c.Log.Info(
			err.SystemMessage(utils.TDefault),
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	c.Log.Debug(
		err.SystemMessage(utils.TDefault),
		mlog.String("err_where", err.Where),
		mlog.Int("http_code", err.StatusCode),
		mlog.String("err_details", err.DetailedError),
	)
}

func (c *Context) RemoveSessionCookie(w http.ResponseWriter, r *http.Request) {
	subpath, _ := utils.GetSubpathFromConfig(c.App.Config())

	cookie := &http.Cookie{
		Name:     model.SESSION_COOKIE_TOKEN,
		Value:    "",
		Path:     subpath,
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (c *Context) SetParamNotFound(where, parameter string) {
	c.Err = model.ParamNotFoundError(where, parameter)
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = model.InvalidParamError(parameter)
}

func (c *Context) SetInvalidUrlParam(parameter string) {
	c.Err = model.NewInvalidUrlParamError(parameter)
}

func (c *Context) SetPermissionError(permission *model.Permission) {
	c.Err = c.App.MakePermissionError(permission)
}

func (c *Context) SetSiteURLHeader(url string) {
	c.siteURLHeader = strings.TrimRight(url, "/")
}

func (c *Context) RequireUserID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.UserID == model.Me {
		c.Params.UserID = c.App.Session().UserID
	}

	return c
}

func (c *Context) RequireSpaceID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.SpaceID == 0 {
		c.SetInvalidUrlParam("space_id")
	}

	return c
}

func (c *Context) RequireProjectID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.ProjectID == 0 {
		c.SetInvalidUrlParam("project_id")
	}

	return c
}

func (c *Context) RequireChatID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.ChatID == 0 {
		c.SetInvalidUrlParam("chat_id")
	}

	return c
}

func (c *Context) RequireTaskID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.TaskID == 0 {
		c.SetInvalidUrlParam("task_id")
	}

	return c
}

func (c *Context) RequireTaskMemoID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.TaskMemoID == 0 {
		c.SetInvalidUrlParam("task_memo_id")
	}

	return c
}

func (c *Context) RequireMediaID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.MediaID == "00000000000000000000000000000000" {
		c.SetInvalidUrlParam("media_id")
	}

	return c
}

func (c *Context) RequireMediaTreeID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.MediaTreeID == 0 {
		fmt.Println(c.Params)
		c.SetInvalidUrlParam("media_tree_id")
	}

	return c
}

func (c *Context) RequireNoteID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.NoteID == 0 {
		c.SetInvalidUrlParam("note_id")
	}

	return c
}

func (c *Context) RequireCommentID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.CommentID == 0 {
		c.SetInvalidUrlParam("comment_id")
	}

	return c
}

func (c *Context) RequireReactionID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.ReactionID == 0 {
		c.SetInvalidUrlParam("reaction_id")
	}

	return c
}

func (c *Context) RequireSharedToken() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.SharedToken == "" {
		c.SetInvalidUrlParam("shared_token")
	}

	return c
}

func (c *Context) RequireFilterID() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.FilterID == 0 {
		c.SetInvalidUrlParam("filter_id")
	}

	return c
}
