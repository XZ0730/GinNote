package router

import (
	Controller "NoteGin/Controller"
	Util "NoteGin/Util"

	_ "NoteGin/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	engine := gin.Default()
	engine.NoRoute(Controller.NoRoute_redirect)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	Apigroup := engine.Group("/api/v1")
	{
		//登录请求
		Apigroup.POST("/login", Controller.Login, Util.SetToken)
		//注册请求
		Apigroup.POST("/register", Controller.Register)
		//添加一条事务

		//检测token
		Apigroup.POST("/todo", Util.CheckToken, Controller.Add)
		//查看一条事务
		Apigroup.GET("/todo/:title", Util.CheckToken, Controller.GetByKey)
		//查看所有/已完成/未完成事务
		Apigroup.GET("/todo", Util.CheckToken, Controller.GetAll)
		//将一条/事务设置为待办或者已完成
		Apigroup.PUT("/todo/:status/:title", Util.CheckToken, Controller.UpdateByOneKey)
		//将所有事务设置为代办或者已完成
		Apigroup.PUT("/todo", Util.CheckToken, Controller.UpdateAll)
		//删除一条事务
		Apigroup.DELETE("/todo/:title", Util.CheckToken, Controller.DeleteByKey)
		//删除所有已经完成
		//删除所有未完成
		//删除所有
		Apigroup.DELETE("/todo", Util.CheckToken, Controller.DeleteAll)

	}
	return engine
}
