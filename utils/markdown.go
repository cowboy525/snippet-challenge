package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type searchPattern struct {
	Pattern string
	GroupID int8
}

func GetRelatedNoteIdsFromMarkdown(spaceID uint64, projectID uint64, markdown string, selfNoteID *uint64) []uint64 {
	noteIDMap := make(map[uint64]struct{})
	escapedWebURL := strings.ReplaceAll(os.Getenv("WEB_URL"), "/", "\\/")
	searchPatterns := []searchPattern{
		{Pattern: fmt.Sprintf(escapedWebURL+"\\/shared\\/note\\/%v\\/%v\\/(\\d+)", spaceID, projectID), GroupID: 1},
		{Pattern: `(^|\n)\[(\\\[|\\\]|[^\[\]\n])+\]:\s*&(\d+)\s*(\n|$)`, GroupID: 3},
		{Pattern: `\[(\\\[|\\\]|[^\[\]\n])+\]\(&(\d+)(\)|\s+\"[^\"\n]+\"\))`, GroupID: 2},
	}
	for _, sp := range searchPatterns {
		r := regexp.MustCompile(sp.Pattern)
		matches := r.FindAllStringSubmatch(markdown, -1)
		for _, v := range matches {
			id, _ := strconv.ParseUint(v[sp.GroupID], 10, 64)
			if id == 0 {
				continue
			}
			if selfNoteID != nil && *selfNoteID == id {
				continue
			}
			if _, ok := noteIDMap[id]; ok {
				continue
			}
			noteIDMap[id] = struct{}{}
		}
	}

	noteIDs := []uint64{}
	for id := range noteIDMap {
		noteIDs = append(noteIDs, id)
	}
	return noteIDs
}

func GetRelatedTaskIdsFromMarkdown(spaceID uint64, projectID uint64, markdown string, selfTaskID *uint64) []uint64 {
	taskIDMap := make(map[uint64]struct{})
	escapedWebURL := strings.ReplaceAll(os.Getenv("WEB_URL"), "/", "\\/")
	searchPatterns := []searchPattern{
		{Pattern: fmt.Sprintf(escapedWebURL+"\\/shared\\/task\\/%v\\/%v\\/(\\d+)", spaceID, projectID), GroupID: 1},
		{Pattern: `(^|\n)\[(\\\[|\\\]|[^\[\]\n])+\]:\s*\!(\d+)\s*(\n|$)`, GroupID: 3},
		{Pattern: `\[(\\\[|\\\]|[^\[\]\n])+\]\(\!(\d+)(\)|\s+\"[^\"\n]+\"\))`, GroupID: 2},
	}
	for _, sp := range searchPatterns {
		r := regexp.MustCompile(sp.Pattern)
		matches := r.FindAllStringSubmatch(markdown, -1)
		for _, v := range matches {
			id, _ := strconv.ParseUint(v[sp.GroupID], 10, 64)
			if id == 0 {
				continue
			}
			if selfTaskID != nil && *selfTaskID == id {
				continue
			}
			if _, ok := taskIDMap[id]; ok {
				continue
			}
			taskIDMap[id] = struct{}{}
		}
	}

	taskIDs := []uint64{}
	for id := range taskIDMap {
		taskIDs = append(taskIDs, id)
	}
	return taskIDs
}

func GetRelatedMediaIdsFromMarkdown(spaceID uint64, projectID uint64, markdown string) []string {
	mediaIDMap := make(map[string]struct{})
	escapedWebURL := strings.ReplaceAll(os.Getenv("WEB_URL"), "/", "\\/")
	searchPatterns := []searchPattern{
		{Pattern: fmt.Sprintf(escapedWebURL+"\\/media\\/%v\\/%v\\/([a-zA-Z0-9-]+)", spaceID, projectID), GroupID: 1},
	}
	for _, sp := range searchPatterns {
		r := regexp.MustCompile(sp.Pattern)
		matches := r.FindAllStringSubmatch(markdown, -1)
		for _, v := range matches {
			id := v[sp.GroupID]
			if _, ok := mediaIDMap[id]; ok {
				continue
			}
			mediaIDMap[id] = struct{}{}
		}
	}

	mediaIDs := []string{}
	for id := range mediaIDMap {
		id := strings.ReplaceAll(id, "-", "")
		mediaIDs = append(mediaIDs, id)
	}
	return mediaIDs
}

func GetMentions(markdown *string, spaceID, projectID uint64, chargeUserIDs []uint64) (map[uint64]int, bool) {
	mentions := map[uint64]int{}
	mentionProject := false
	if markdown != nil {

		// members
		if exist := strings.Contains(*markdown, "{{{mention:members}}}"); exist {
			for _, chargeUserID := range chargeUserIDs {
				mentions[chargeUserID] = 0
			}
		}

		// mention project members
		if exist := strings.Contains(*markdown, "{{{mention:project}}}"); exist {
			mentionProject = true
		}

		// user
		r := regexp.MustCompile(`\{\{\{mention:(\d+)(:(\d+))?\}([^\r\n{}@]+)\}\}`)
		matches := r.FindAllStringSubmatch(*markdown, -1)
		for _, v := range matches {
			id, _ := strconv.ParseUint(v[1], 10, 64)
			level := 0
			if len(v[3]) > 0 {
				level = 1
			}
			if lv, ok := mentions[id]; ok {
				if lv < level {
					delete(mentions, id)
					mentions[id] = level
				}
			} else {
				mentions[id] = level
			}
		}
	}

	return mentions, mentionProject
}
