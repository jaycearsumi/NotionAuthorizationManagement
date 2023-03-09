package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	rm := createRootManager()
	rm.start()
}

type rootManager struct {
	auth string
	am   *aclManager
	r    *gin.Engine
}

func createRootManager() *rootManager {
	return &rootManager{
		am: createACLManager(),
		r:  gin.Default(),
	}
}

func (rm *rootManager) start() {
	rm.useMiddleware()
	rm.aclRegister()
	rm.rbac3Register()
	rm.r.Run()
}

func (rm *rootManager) useMiddleware() {
	rm.r.POST("/auth", func(c *gin.Context) {
		var params authReqModel
		c.ShouldBind(&params)
		rm.auth = params.Token
		rm.am.api.auth = params.Token
		c.String(http.StatusOK, "OK")
	})
}

func (rm *rootManager) aclRegister() {
	aclGroup := rm.r.Group("acl")
	aclGroup.Use(func(c *gin.Context) {
		if rm.auth == "" {
			c.String(http.StatusForbidden, "need auth token")
			c.Abort()
		} else {
			c.Next()
		}
	})
	aclGroup.POST("/init", rm.am.aclInit())
	aclGroup.POST("/create", rm.am.aclCreate())
	aclGroup.GET("/:rid/:uid", rm.am.aclCheck())
	aclGroup.PATCH("/:rid/:uid", rm.am.aclUpdate())
	aclGroup.DELETE("/:rid/:uid", rm.am.aclRevoke())
}

func (rm *rootManager) rbac3Register() {

}
