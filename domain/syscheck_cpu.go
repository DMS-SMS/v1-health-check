// Create file in v.1.0.0
// syscheck_cpu.go is file that declare model struct & repo interface about cpu health check in syscheck domain.
// also, additional method of model struct is declared in this file, too.

package domain

import (
	"context"
)

// CPUCheckHistory model is used for record cpu health check history and result
type CPUCheckHistory struct {
	// get required component by embedding systemCheckHistoryComponent
	systemCheckHistoryComponent

	// TotalUsageCore specifies current total cpu usage of runtime system looked in cpu check
	TotalUsageCore float64

	// DockerUsageCore specifies current total cpu usage of docker looked in cpu check when weak detected
	DockerUsageCore float64

	// TemporaryFreeCore specifies temporary freed cpu size while recovering cpu health
	TemporaryFreeCore float64

	// MostCPUConsumeContainer specifies the container name which is consumed most CPU
	MostCPUConsumeContainer string
}

// CPUCheckHistoryRepository is interface for repository layer used in usecase layer
// Repository is implemented with elasticsearch in v.1.0.0
type CPUCheckHistoryRepository interface {
	// get required component by embedding systemCheckHistoryRepositoryComponent
	systemCheckHistoryRepositoryComponent

	// Store method save CPUCheckHistory model in repository
	// b in return represents bytes of response body(map[string]interface{})
	Store(*CPUCheckHistory) (b []byte, err error)
}

// DiskCheckUseCase is interface used as business process handler about cpu check
type CPUCheckUseCase interface {
	// CheckCPU method check cpu usage status and store cpu check history using repository
	CheckCPU(ctx context.Context) error
}

// FillPrivateComponent overriding FillPrivateComponent method of systemCheckHistoryComponent
func (ch *CPUCheckHistory) FillPrivateComponent() {
	ch.systemCheckHistoryComponent.FillPrivateComponent()
	ch._type = "CPUCheck"
}

// DottedMapWithPrefix convert CPUCheckHistory to dotted map and return using MapWithPrefixKey of upper struct
// all key value of Map start with prefix received from parameter
func (ch *CPUCheckHistory) DottedMapWithPrefix(prefix string) (m map[string]interface{}) {
	m = ch.systemCheckHistoryComponent.DottedMapWithPrefix(prefix)

	if prefix != "" {
		prefix += "."
	}

	// setting public field value in dotted map
	m[prefix+"total_usage_core"] = ch.TotalUsageCore
	m[prefix+"docker_usage_core"] = ch.DockerUsageCore
	m[prefix+"temporary_free_core"] = ch.TemporaryFreeCore
	m[prefix+"most_cpu_consume_container"] = ch.MostCPUConsumeContainer

	return
}
