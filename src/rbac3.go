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

func (rm *rbac3Manager) rbac3UpdateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")
		roleStr := c.Param("role")
		params := &rbac3UpdateRoleReqModel{
			UserID: uid,
			Role:   roleStr,
		}
		res, err := rm.api.updateRole(params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (rm *rbac3Manager) rbac3RevokeRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.Param("uid")
		roleStr := c.Param("role")
		params := &rbac3UpdateRoleReqModel{
			UserID: uid,
			Role:   roleStr,
		}
		res, err := rm.api.revokeRole(params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (rm *rbac3Manager) rbac3UpdateAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params rbac3UpdateAccessReqModel
		c.ShouldBind(&params)
		params.ResourceID = c.Param("rid")
		params.Role = c.Param("role")
		res, err := rm.api.updateAccess(&params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}

func (rm *rbac3Manager) rbac3RevokeAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := &rbac3UpdateAccessReqModel{
			ResourceID: c.Param("rid"),
			Role:       c.Param("role"),
		}
		res, err := rm.api.revokeAccess(params)
		if err != nil {
			c.String(http.StatusNotFound, err.Error())
		} else {
			c.JSON(http.StatusOK, gin.H{
				"res": res,
			})
		}
	}
}
