package messaging

import (
	"testing"
)

func TestMsgBuild(t *testing.T) {
	msg := Message{
		Content: "random content",
	}

	actual_json_msg := msg.Build()

	expected_json_msg := `{
  "id": "",
  "pubKey": "",
  "createdAt": 0,
  "kind": 0,
  "tags": null,
  "content": "random content",
  "sig": ""
 }`

	if actual_json_msg != expected_json_msg {
		t.Error(actual_json_msg, expected_json_msg)
	}

}

func TestMsgSign(t *testing.T) {
	privateKey := "22a47fa09a223f2aa079edf85a7c2d4f8720ee63e502ee2869afab7de234b80c"

	msg := Message{
		Content:   "Random Content",
		Kind:      0,
		Tags:      [][]string{[]string{"a tag"}},
		CreatedAt: 33,
	}

	msg.Sign(privateKey)

	if !msg.Verify() {
		t.Error("Verification of msg failed!")
	}

}
