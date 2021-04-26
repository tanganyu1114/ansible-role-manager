package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/tanganyu1114/ansible-role-manager/app/admin/apis/ansible/inventory"
	"github.com/tanganyu1114/ansible-role-manager/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAnsibleInventoryRouter)
}

// 需认证的路由代码
func registerAnsibleInventoryRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := inventory.NewInventoryApi()
	r := v1.Group("/ansible/inventory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.POST("/groups", api.AddHostToGroup)
		r.PUT("/groups", api.RenewGroupName)
		r.PATCH("/groups", api.RemoveHostFromGroup)
		r.DELETE("/groups/:group", api.RemoveGroupByName)
		r.GET("/hosts", api.GetAllHosts)
		r.GET("/groups", api.GetGroups)
	}
}
