package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/tanganyu1114/ansible-role-manager/app/admin/apis/ansible/roles"
	"github.com/tanganyu1114/ansible-role-manager/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerAnsibleRolesRouter)
}

// 需认证的路由代码
func registerAnsibleRolesRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := roles.NewRolesApi()
	r := v1.Group("ansible/roles").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		r.POST("/:role", api.AddRoleByCompressedData)
		r.GET("/:role", api.DownloadRoleCompressedData)
		r.DELETE("/:role", api.RemoveRole)
	}
}
