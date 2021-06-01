package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// syscheckHandler represent the http handler for syscheck
type syscheckHandler struct {
	dUsecase domain.DiskCheckUseCase
	cUsecase domain.CPUCheckUseCase
	mUsecase domain.MemoryCheckUseCase
}

// NewSyscheckHandler initialize the resources of syscheck domain to HTTP API endpoint
func NewSyscheckHandler(r *gin.Engine, du domain.DiskCheckUseCase, cu domain.CPUCheckUseCase, mu domain.MemoryCheckUseCase) {
	h := &syscheckHandler{
		dUsecase: du,
		cUsecase: cu,
		mUsecase: mu,
	}

	r.POST("system-check/types/disk", h.CheckDisk)
	r.POST("system-check/types/cpu", h.CheckCPU)
	r.POST("system-check/types/memory", h.CheckMemory)
}

// CheckDisk method deliver HTTP request to CheckDisk method of domain.DiskCheckUseCase
func (sh *syscheckHandler) CheckDisk(c *gin.Context) {
	switch err := sh.dUsecase.CheckDisk(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check disk status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check disk status").Error(),
		})
	}
}

// CheckCPU method deliver HTTP request to CheckCPU method of domain.CPUCheckUseCase
func (sh *syscheckHandler) CheckCPU(c *gin.Context) {
	switch err := sh.cUsecase.CheckCPU(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check cpu status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check cpu status").Error(),
		})
	}
}

// CheckMemory method deliver HTTP request to CheckMemory method of domain.MemoryCheckUseCase
func (sh *syscheckHandler) CheckMemory(c *gin.Context) {
	switch err := sh.mUsecase.CheckMemory(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check memory status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check memory status").Error(),
		})
	}
}
