package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"

	"github.com/DMS-SMS/v1-health-check/domain"
)

// srvcheckHandler represent the http handler for srvcheck
type srvcheckHandler struct {
	cUsecase domain.ConsulCheckUseCase
	eUsecase domain.ElasticsearchCheckUseCase
	sUsecase domain.SwarmpitCheckUseCase
}

// NewSrvcheckHandler initialize the resources of srvcheck domain to HTTP API endpoint
func NewSrvcheckHandler(r *gin.Engine, cu domain.ConsulCheckUseCase, eu domain.ElasticsearchCheckUseCase, su domain.SwarmpitCheckUseCase) {
	h := &srvcheckHandler{
		cUsecase: cu,
		eUsecase: eu,
		sUsecase: su,
	}

	r.POST("service-check/types/disk", h.CheckConsul)
	r.POST("service-check/types/cpu", h.CheckElasticsearch)
	r.POST("service-check/types/memory", h.CheckSwarmpit)
}

// CheckConsul method deliver HTTP request to CheckConsul method of domain.ConsulCheckUseCase
func (sh *srvcheckHandler) CheckConsul(c *gin.Context) {
	switch err := sh.cUsecase.CheckConsul(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check consul status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check consul status").Error(),
		})
	}
}

// CheckElasticsearch method deliver HTTP request to CheckElasticsearch method of domain.ElasticsearchCheckUseCase
func (sh *srvcheckHandler) CheckElasticsearch(c *gin.Context) {
	switch err := sh.eUsecase.CheckElasticsearch(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check elasticsearch status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check elasticsearch status").Error(),
		})
	}
}

// CheckSwarmpit method deliver HTTP request to CheckSwarmpit method of domain.SwarmpitCheckUseCase
func (sh *srvcheckHandler) CheckSwarmpit(c *gin.Context) {
	switch err := sh.sUsecase.CheckSwarmpit(c.Request.Context()); err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "code": 0, "message": "finished to check swarmpit status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError, "code": 0,
			"message": errors.Wrap(err, "failed to check swarmpit status").Error(),
		})
	}
}
