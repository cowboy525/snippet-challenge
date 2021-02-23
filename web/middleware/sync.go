package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ernie-mlg/ErniePJT-main-api-go/app"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
)

func Synchronization(a app.Iface) {
	for _, event := range a.Session().SyncEvents {
		syncURL := os.Getenv("SYNC_URL")
		if len(syncURL) > 0 {
			go postSyncEvent(syncURL, event.EventData)
		}

	}
}

func postSyncEvent(syncURL string, eventData *model.EventData) {
	client := &http.Client{}

	b, _ := json.Marshal(eventData)
	rq, err := http.NewRequest("POST", syncURL+"/serve/events/", bytes.NewReader(b))
	if err != nil {
		mlog.Error("An error occured while posting sync data" + string(b))
		return
	}

	rq.Header.Set("Content-Type", "application/json")

	rp, err := client.Do(rq)
	if err != nil || rp == nil {
		mlog.Error("An error occured while posting sync data")
	}

	defer func(r *http.Response) {
		if r != nil && r.Body != nil {
			_, _ = ioutil.ReadAll(r.Body)
			_ = r.Body.Close()
		}
		if r := recover(); r != nil {
			fmt.Println("Recovered in postSyncEvent", r)
		}
	}(rp)
}
