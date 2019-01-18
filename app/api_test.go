package app

import (
	"encoding/json"
	"testing"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestResponse4xx(t *testing.T) {
	rawResp := []byte(`{"ok":false,"error_code":429,"description":"Too Many Requests: retry after 44","parameters":{"retry_after":44}}`)

	var respAPI telegram.APIResponse
	err := json.Unmarshal(rawResp, &respAPI)
	if err != nil {
		t.Errorf("expected err is nil, actual %v", err)
		return
	}

	assertEqual(t, false, respAPI.Ok)
	assertEqual(t, 429, respAPI.ErrorCode)
}
