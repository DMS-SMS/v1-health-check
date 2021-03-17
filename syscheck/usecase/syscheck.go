// Create package in v.1.0.0
// usecase package declare implementation of usecase interface about syscheck domain
// all usecase implementation will accept any input from Delivery layer
// This usecase layer will depends to Repository layer

// syscheck.go is file that define structure to embed from another structures.
// It also defines variables or constants, functions used jointly in the package as private.

package usecase

// global variable used in usecase which is type of processLevel
var (
	healthyLevel       = processLevel("HEALTHY")       // represent that system status is healthy now
	weakDetectedLevel  = processLevel("WEAK_DETECTED") // represent that weak of system status is detected
	recoveringLevel    = processLevel("RECOVERING")    // represent that recovering weak of system status now
	unhealthyLevel     = processLevel("UNHEALTHY")     // represent that system status is unhealthy now (not recovered)
	errorLevel         = processLevel("ERROR")         // represent that error occurs while checking system status
)

// systemCheckUsecaseComponent contains required component to syscheck usecase implementation as field
type systemCheckUsecaseComponentConfig interface {}

// processLevel is string custom type used for representing status check process level
type processLevel string

// String method return instance value of processLevel type to string
func (pl processLevel) String() string {
	return string(pl)
}
