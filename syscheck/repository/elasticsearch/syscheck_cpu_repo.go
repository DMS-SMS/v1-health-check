// Create file in v.1.0.0
// syscheck_cpu_repo.go is file that define implement cpu history repository using elasticsearch
// this cpu repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

// esCPUCheckHistoryRepoConfig is the config for cpu check history repository using elasticsearch
type esCPUCheckHistoryRepoConfig interface {
	// get common method from embedding esRepositoryComponentConfig
	esRepositoryComponentConfig
}
