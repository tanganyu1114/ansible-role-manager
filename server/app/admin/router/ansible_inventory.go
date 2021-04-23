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
	api, err := inventory.NewInventoryApi()
	if err != nil {
		panic(err)
	}
	r := v1.Group("/ansible/inventory").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.POST("", api.AddHostToGroup)
		r.POST("", api.RenewGroupName)
		r.POST("", api.RemoveHostFromGroup)
		r.POST("", api.RemoveGroupByName)
		r.GET("", api.GetAllHosts)
		r.GET("", api.GetGroups)
	}
}
