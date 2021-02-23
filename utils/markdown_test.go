package utils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRelatedNoteIdsFromMarkdown(t *testing.T) {
	spaceID := uint64(1)
	projectID := uint64(100)
	noteID := uint64(1000)
	selfNoteID := uint64(10000)
	markdown := "aaaa" + fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/note/%v/%v/%v", spaceID, projectID, noteID) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/note/%v/%v/%v", spaceID, projectID, noteID+1) + "\n" +
		fmt.Sprintf("[testtest]:&%v", noteID+2) + "\n" +
		fmt.Sprintf("[testtest](&%v)", noteID+3) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/note/%v/%v/%v", spaceID, projectID, selfNoteID) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/note/%v/%v/%v", spaceID+1, projectID, selfNoteID+1) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/note/%v/%v/%v", spaceID, projectID+1, selfNoteID+2) + "\n" +
		fmt.Sprintf("[testtest]:!%v", selfNoteID+3) + "\n" +
		fmt.Sprintf("[testtest](!%v)", selfNoteID+4) + "\n" +
		fmt.Sprintf("[testtest]:&%v", noteID+2)
	results := GetRelatedNoteIdsFromMarkdown(spaceID, projectID, markdown, &selfNoteID)
	assert.ElementsMatch(t, []uint64{noteID, noteID + 1, noteID + 2, noteID + 3}, results)
}

func TestGetRelatedTaskIdsFromMarkdown(t *testing.T) {
	spaceID := uint64(1)
	projectID := uint64(100)
	taskID := uint64(1000)
	selfTaskID := uint64(10000)
	markdown := "aaaa" + fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID, projectID, taskID) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID, projectID, taskID+1) + "\n" +
		fmt.Sprintf("[testtest]:!%v", taskID+2) + "\n" +
		fmt.Sprintf("[testtest](!%v)", taskID+3) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID, projectID, selfTaskID) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID+1, projectID, selfTaskID+1) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID, projectID+1, selfTaskID+2) + "\n" +
		fmt.Sprintf("[testtest]:&%v", selfTaskID+3) + "\n" +
		fmt.Sprintf("[testtest](&%v)", selfTaskID+4) + "\n" +
		fmt.Sprintf("[testtest]:!%v", taskID+2)

	results := GetRelatedTaskIdsFromMarkdown(spaceID, projectID, markdown, &selfTaskID)
	assert.ElementsMatch(t, []uint64{taskID, taskID + 1, taskID + 2, taskID + 3}, results)
}

func TestGetRelatedMediaIdsFromMarkdown(t *testing.T) {
	spaceID := uint64(1)
	projectID := uint64(100)
	mediaID := "2fac123431f811b4a22208002b34c003"
	markdown := "aaaa" + fmt.Sprintf(os.Getenv("WEB_URL")+"/media/%v/%v/%v", spaceID, projectID, mediaID) + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID+1, projectID, mediaID+"1") + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/shared/task/%v/%v/%v", spaceID, projectID+1, mediaID+"2") + "\n" +
		fmt.Sprintf(os.Getenv("WEB_URL")+"/media/%v/%v/%v", spaceID, projectID, mediaID)
	results := GetRelatedMediaIdsFromMarkdown(spaceID, projectID, markdown)
	assert.ElementsMatch(t, []string{mediaID}, results)
}
