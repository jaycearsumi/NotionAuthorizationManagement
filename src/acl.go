package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type aclManager struct {
	databaseID string
	parentID   string
	api        *aclAPI
}

func createACLManager() *aclManager {
	am := &aclManager{
		api: createACLAPI(),
	}
	return am
}

func (am *aclManager) aclInit() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params initReqModel
		c.ShouldBind(&params)
		res, err := am.api.retrieveDatabase(&params)
		//log.Println(res)
		//log.Println(err)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			am.databaseID = params.DatabaseID
			am.parentID = res.Parent.id()
			am.api.databaseID = params.DatabaseID
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (am *aclManager) aclCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params createReqModel
		c.ShouldBind(&params)
		res, err := am.api.createDatabase(&params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			am.databaseID = res.DatabaseID
			am.api.databaseID = res.DatabaseID
			am.parentID = res.parent.id()
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (am *aclManager) aclCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Param("rid")
		uid := c.Param("uid")
		params := &checkReqModel{
			ResourceID: rid,
			UserID:     uid,
		}
		res, err := am.api.checkAccess(params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (am *aclManager) aclUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Param("rid")
		uid := c.Param("uid")
		checkParams := &checkReqModel{
			ResourceID: rid,
			UserID:     uid,
		}
		var bodyParams updateReqModel
		c.ShouldBind(&bodyParams)

		res, err := am.api.updateAccess(&bodyParams, checkParams)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (am *aclManager) aclRevoke() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Param("rid")
		uid := c.Param("uid")
		checkParams := &checkReqModel{
			ResourceID: rid,
			UserID:     uid,
		}
		res, err := am.api.revokeAccess(checkParams)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}
