// Create file in v.1.0.0
// syscheck_disk.go is file that define repository implement about disk using elasticsearch
// disk repository struct embed esRepositoryRequiredComponent struct in ./syscheck.go file

package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// esDiskCheckHistoryRepository is to handle DiskCheckHistory model using elasticsearch as data store
type esDiskCheckHistoryRepository struct {
	// get common field from embedding esRepositoryRequiredComponent
	esRepositoryRequiredComponent

	// elasticsearch client connection injected from the outside package
	cli *elasticsearch.Client
}
