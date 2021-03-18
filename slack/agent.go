// Create package in v.1.0.0
// slack package define struct which is implement various interface about slack agency using in each of domain
// there are kind of slack agency such as chat, conversations, admin, etc.

// in agent.go file, define struct type of slack agent & initializer that are not method.
// Also if exist, custom type or variable used in common in each of method will declared in this file.

package slack

import "github.com/slack-go/slack"

// slackAgent agent various slack API(chat, conversations, admin, etc ...) as implementation
type slackAgent struct {
	// slkCli is slack client connection injected from the outside package
	slkCli *slack.Client
}

// NewAgent return new initialized instance of slackAgent pointer type with slack API token & options
func NewAgent(token string, opts ...slack.Option) *slackAgent {
	return &slackAgent{
		slkCli: slack.New(token, opts...),
	}
}
