package websockets
import (
	"encoding/json"
	"fmt"
)
// messageType is an enum https://www.sohamkamani.com/golang/enums/
// by using an enum for message types we can guarantee to never have
// any kinds of message that we don't explicitly create
type messageType int64
const (
	unknown messageType = iota
	chat
	post
	content
	presence
)
func (m messageType) String() string {
	switch m {
	case chat:
		return "chat"
	case post:
		return "post"
	case content:
		return "content"
	case presence:
		return "presence"
	default:
		return "unknown"
	}
}
func ParseMessageType(s string) (messageType, error) {
	switch s {
	case "chat":
		return chat, nil
	case "post":
		return post, nil
	case "content":
		return content, nil
	case "presence":
		return presence, nil
	default:
		return unknown, fmt.Errorf("unknown message type %s", s)
	}
}
func (m messageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}
func (m *messageType) UnmarshalJSON(data []byte) (err error) {
	var mt string
	if err := json.Unmarshal(data, &mt); err != nil {
		return err
	}
	if *m, err = ParseMessageType(mt); err != nil {
		return err
	}
	return nil
}
