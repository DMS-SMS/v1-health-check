// Create package in v.1.0.0
// usecase package declare implementation of usecase interface about service check(srvcheck) domain
// all usecase implementation will accept any input from Delivery layer
// This usecase layer will depends to Repository layer

// srvcheck.go is file that define structure to embed from another structures.
// It also defines variables or constants, functions used jointly in the package as private.

package usecase

import (
	"github.com/slack-go/slack"
	"time"
)

// global variable used in usecase to represent process level
const (
	healthyLevel      = "HEALTHY"       // represent that service status is healthy now
	warningLevel      = "WARNING"       // represent that service status is warning now
	weakDetectedLevel = "WEAK_DETECTED" // represent that weak of service status is detected
	recoveringLevel   = "RECOVERING"    // represent that recovering weak of service status now
	recoveredLevel    = "RECOVERED"     // represent that succeed to recover service status
	unhealthyLevel    = "UNHEALTHY"     // represent that service status is unhealthy now (not recovered)
	errorLevel        = "ERROR"         // represent that error occurs while checking service status
)

// serviceCheckUsecaseComponentConfig contains required component to service usecase implementation as field
type serviceCheckUsecaseComponentConfig interface {}

// slackChatAgency is interface that agent the slack api about chatting
// you can see implementation in slack package
type slackChatAgency interface {
	// SendMessage send message with text & emoji using slack API and return send time & text & error
	SendMessage(emoji, text, uuid string, opts ...slack.MsgOption) (t time.Time, _text string, err error)
}

// intComparator is struct type having int type field which is used for compare with another int
type intComparator struct { V int }

// isMoreThan return boolean if value of instance which call this method is more than parameter's size
func (comparator intComparator) isMoreThan(target int) bool { return comparator.V > target }

// isMoreThan return boolean if value of instance which call this method is less than parameter's size
func (comparator intComparator) isLessThan(target int) bool { return comparator.V < target }
