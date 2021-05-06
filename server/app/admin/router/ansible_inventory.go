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
		r.POST("/groups", api.NewGroup)
		r.PATCH("/groups", api.ModifyGroup)
		r.DELETE("/groups/:group", api.DeleteGroup)
		r.GET("/groups", api.GetGroups)
	}
}
