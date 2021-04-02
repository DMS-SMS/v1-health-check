// Create file in v.1.0.0
// srvcheck_swarmpit.go is file that declare model struct & repo interface about swarmpit check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

// SwarmpitCheckHistory model is used for record swarmpit check history and result
type SwarmpitCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent
}
