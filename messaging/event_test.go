package messaging

import (
	"testing"
)

func TestEventBuild(t *testing.T) {
	event := Event{
		Content: "random content",
	}

	actual_json_event := event.Build()

	expected_json_event := `{
  "id": "",
  "pubKey": "",
  "createdAt": 0,
  "kind": 0,
  "tags": null,
  "content": "random content",
  "sig": ""
 }`

	if actual_json_event != expected_json_event {
		t.Error(actual_json_event, expected_json_event)
	}

}

func TestEventSign(t *testing.T) {
	privateKey := "22a47fa09a223f2aa079edf85a7c2d4f8720ee63e502ee2869afab7de234b80c"

	event := Event{
		Content:   "Random Content",
		Kind:      0,
		Tags:      [][]string{{"a tag"}},
		CreatedAt: 33,
	}

	event.Sign(privateKey)

	if !event.Verify() {
		t.Error("Verification of event failed!")
	}

}
