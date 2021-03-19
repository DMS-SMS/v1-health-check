// Create file in v.1.0.0
// agent_chat.go file define method of slackAgent about slack chat API
// implement agency interface about slack chat defined in each of domain

package slack

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"strconv"
	"time"
)

// SendMessage send message with text & emoji using slack API and return send time & error
func (sa *slackAgent) SendMessage(emoji, text string, opts ...slack.MsgOption) (t time.Time, err error) {
	if emoji != "" {
		text = fmt.Sprintf(":%s: %s", emoji, text)
	}

	opts = append(opts, slack.MsgOptionText(text, false))
	_, _time, _, err := sa.slkCli.SendMessage(sa.chatChannel, opts...)
	if err != nil {
		err = errors.Wrap(err, "failed to send message with slack API")
		return
	}

	if len(_time) >= 10 {
		i, _ := strconv.ParseInt(_time[:10], 10, 64)
		t = time.Unix(i, 0)
		if t.Location().String() == time.UTC.String() {
			t = t.Add(time.Hour * 9)
		}
	}
	return
}
