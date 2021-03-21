// Create package in v.1.0.0
// usecase package declare implementation of usecase interface about syscheck domain
// all usecase implementation will accept any input from Delivery layer
// This usecase layer will depends to Repository layer

// syscheck.go is file that define structure to embed from another structures.
// It also defines variables or constants, functions used jointly in the package as private.

package usecase

import (
	"github.com/slack-go/slack"
	"time"
)

// global variable used in usecase which is type of processLevel
var (
	healthyLevel      = processLevel("HEALTHY")       // represent that system status is healthy now
	weakDetectedLevel = processLevel("WEAK_DETECTED") // represent that weak of system status is detected
	recoveringLevel   = processLevel("RECOVERING")    // represent that recovering weak of system status now
	recoveredLevel    = processLevel("RECOVERED")     // represent that succeed to recover system status
	unhealthyLevel    = processLevel("UNHEALTHY")     // represent that system status is unhealthy now (not recovered)
	errorLevel        = processLevel("ERROR")         // represent that error occurs while checking system status
)

// systemCheckUsecaseComponent contains required component to syscheck usecase implementation as field
type systemCheckUsecaseComponentConfig interface {}

// slackChatAgency is interface that agent the slack api about chatting
// you can see implementation in slack package
type slackChatAgency interface {
	// SendMessage send message with text & emoji using slack API and return send time & text & error
	SendMessage(emoji, text, uuid string, opts ...slack.MsgOption) (t time.Time, _text string, err error)
}

// processLevel is string custom type used for representing status check process level
type processLevel string

// String method return instance value of processLevel type to string
func (pl processLevel) String() string {
	return string(pl)
}
