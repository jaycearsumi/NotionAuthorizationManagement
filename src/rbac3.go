package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type rbac3Manager struct {
	homePageID string
	parentID   string
	api        *rbac3API
}

func createRbac3Manager() *rbac3Manager {
	rm := &rbac3Manager{
		api: createRbac3API(),
	}
	return rm
}

func (rm *rbac3Manager) rbac3Init() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params rbac3InitReqModel
		c.ShouldBind(&params)
		res, err := rm.api.homePageInit(&params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			rm.homePageID = params.HomePageID
			//rm.parentID = res.Parent.id()
			rm.api.homePageID = params.HomePageID
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (rm *rbac3Manager) rbac3Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Param("rid")
		uid := c.Param("uid")
		params := &checkReqModel{
			ResourceID: rid,
			UserID:     uid,
		}
		res, err := rm.api.rbac3CheckAccess(params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}
