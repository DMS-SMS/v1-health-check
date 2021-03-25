// Create file in v.1.0.0
// syscheck_mem.go is file that declare model struct & repo interface about memory health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

// MemCheckHistory model is used for record memory health check history and result
type MemCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// UsageSize specifies current cpu usage of runtime system looked in cpu check
	//UsageSize bytesize.ByteSize
	//
	//// FreeSize specifies freed cpu size while recovering cpu health
	//FreeSize float64
	//
	//// MostCPUConsumeContainer specifies the container name which is consumed most CPU
	//MostCPUConsumeContainer string
}
