package main

import (
	"github.com/gin-gonic/gin"
	"jkdev.cn/api/controller"
	"jkdev.cn/api/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	//添加了 跨域 的 filter 和 异常时 返回 "code 500" 响应的 filter
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	//指定了 /api/auth/register 和 相应的 controller   (controller.Register)
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	/*
		把 categories 相关的 增删查改 放到一起了
	*/
	//用了 group 的概念，即 "公共前缀", 组内的 url，都会有这个 前缀
	categoryRoutes := r.Group("/categories")
	//创建了一个新的 controller
	categoryController := controller.NewCategoryController()
	// 指定了 POST /categories 和 相应的 controller  (categoryController.Create)
	categoryRoutes.POST("", categoryController.Create)
	// 指定了 PUT /categories/{id} 和 相应的 controller  (categoryController.Update)
	categoryRoutes.PUT("/:id", categoryController.Update) //替换
	// 指定了 GET /categories/{id} 和 相应的 controller  (categoryController.Show)
	categoryRoutes.GET("/:id", categoryController.Show)
	// 指定了 DELETE /categories/{id} 和 相应的 controller  (categoryController.Update)
	categoryRoutes.DELETE("/:id", categoryController.Delete)

	/*
		把 posts 相关的 增删查改 放到一起了，同上
	*/
	postRoutes := r.Group("/posts")
	postRoutes.Use(middleware.AuthMiddleware())
	postController := controller.NewPostController()
	postRoutes.POST("", postController.Create)
	postRoutes.PUT("/:id", postController.Update) //替换
	postRoutes.GET("/:id", postController.Show)
	postRoutes.DELETE("/:id", postController.Delete)
	postRoutes.POST("/page/list", postController.PageList)

	return r
}
