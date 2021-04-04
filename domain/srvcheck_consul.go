// Create file in v.1.0.0
// srvcheck_consul.go is file that declare model struct & repo interface about consul check in srvcheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

// ConsulCheckHistory model is used for record consul check history and result
type ConsulCheckHistory struct {
	// get required component by embedding serviceCheckHistoryComponent
	serviceCheckHistoryComponent

	// InstancesNumPerSrv specifies instances number per service in consul
	InstancesNumPerSrv map[string]int

	// DeregisteredServiceIDs specifies id list of deregistered service in consul check
	DeregisteredServiceIDs []string

	// IfServiceDeregister specifies if any service in consul was deregistered
	IfServiceDeregistered bool
}
