// Create package in v.1.0.0
// usecase package declare implementation of usecase interface about service check(srvcheck) domain
// all usecase implementation will accept any input from Delivery layer
// This usecase layer will depends to Repository layer

// srvcheck.go is file that define structure to embed from another structures.
// It also defines variables or constants, functions used jointly in the package as private.

package usecase

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
