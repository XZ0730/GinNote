package controller

import (
	"fmt"
	Model "gin01/Models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// NoRoute_redirect
func NoRoute_redirect(ctx *gin.Context) {
	ctx.Redirect(http.StatusMovedPermanently, "http://localhost:8080/api/v1/login")
}

// Login
// @Description 登录
// @Summary 登录并返回token
// @Accept application/json
// @Produce application/json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /login [POST]

func Login(ctx *gin.Context) {
	//接收前端传送的用户名和密码 json形式
	var user Model.User
	var user1 Model.User
	ctx.ShouldBind(&user)
	//从请求中把数据拿出来					//准备编写bcrpt加密 进行验证
	fmt.Println(user)
	d := Model.Db1.Where("username=?", user.Username).First(&user1)
	if d.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 401,
			"err":    "用户名错误",
			"msg":    "用户名错误",
		})
		ctx.Request.URL.Path = "/api/v1/login"
		ctx.Abort()
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user1.Password), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 401,
			"err":    err,
			"msg":    "密码错误",
		})

		ctx.Request.URL.Path = "/api/v1/login"
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "登录成功",
	})
	ctx.Redirect(http.StatusMovedPermanently, "http://localhost:8080/api/v1/todo")
	ctx.Next()
}

// Register
// @Description 注册
// @Summary 注册用户名和密码
// @Accept application/json
// @Produce application/json
// @Param username application/json string true "用户名"
// @Param password application/json string true "密码"
// @Param email application/json string true "邮箱"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /register [POST]
func Register(ctx *gin.Context) {
	var user Model.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 401,
			"err":    err,
			"msg":    err,
		})
		return
	}
	if user.Email == "" || user.Password == "" || user.Username == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 401,
			"msg":    "全部都要填哦",
		})
		ctx.Redirect(http.StatusMovedPermanently, "http://localhost:8080/api/v1/register")
		return
	}
	//生成密文
	pa, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"err":    err,
			"msg":    "生成密文错误",
		})
		return
	}
	user.Password = string(pa)
	if err := Model.Db1.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"msg":    "注册错误",
			"err":    err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "注册成功",
	})
}

// Add
// @Description 添加记录
// @Summary 添加一条记录
// @Accept application/json
// @Produce application/json
// @Param title json string false "标题"
// @Param content json string false "内容"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo [POST]
func Add(ctx *gin.Context) {
	var note Model.Note
	var data Model.Data
	err := ctx.BindJSON(&note)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 402,
			"msg":    "json绑定失败",
		})
	}
	note.Create_time = time.Now().Unix()
	note.Start_time = time.Now().Unix()
	Model.Db1.Create(&note)
	data.Item = append(data.Item, note)
	data.Total = 1
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"msg":    "添加成功",
		"error":  "",
	})
}

// GetByKey
// @Description 查询关键词为title的记录
// @Summary 查询一条记录
// @Accept path
// @Produce application/json
// @Param title path string true "标题"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo/{title} [GET]
func GetByKey(ctx *gin.Context) {
	var note1 []Model.Note
	var page int
	title1 := ctx.Param("title")
	title := "%" + title1 + "%"
	page1, ok := ctx.GetQuery("page")
	if ok {
		page, _ = strconv.Atoi(page1)
	} else {
		page = 1
	}
	d := Model.Db1.Where("title LIKE ?", title).Limit(2).Offset((page - 1) * 2).Find(&note1)
	if d.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"msg":    "无此记录",
			"error":  d.Error,
		})
		return
	}
	var data Model.Data
	data.Item = append(data.Item, note1...)
	data.Total = d.RowsAffected
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"msg":    "ok",
		"error":  "",
	})
}

// GetAll
// @Description 查询记录
// @Summary 查询所有记录/所有未完成/所有已完成记录
// @Accept query
// @Produce application/json
// @Param status query int false "完成状态"
// @Param page query int false "页码"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo [GET]
func GetAll(ctx *gin.Context) {
	var note2 []Model.Note
	var data Model.Data
	var page int
	var status1 int
	status, ok1 := ctx.GetQuery("status") //  /todo?status=1&page=1
	page1, ok2 := ctx.GetQuery("page")
	if ok2 {
		page, _ = strconv.Atoi(page1)
	} else {
		page = 1
	}

	if !ok1 {
		d := Model.Db1.Limit(2).Offset((page - 1) * 2).Find(&note2)
		if d.Error != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   Model.Data{},
				"msg":    "查询错误",
				"error":  d.Error,
			})
			return
		}
		data.Item = append(data.Item, note2...)
		data.Total = int64(len(note2))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
			"msg":    "ok",
			"error":  "",
		})
	} else {
		status1, _ = strconv.Atoi(status)
		d := Model.Db1.Where("done_status=?", status1).Limit(2).Offset((page - 1) * 2).Find(&note2)
		if d.Error != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   Model.Data{},
				"msg":    "查询错误",
				"error":  d.Error,
			})
			return
		}
		data.Item = append(data.Item, note2...)
		data.Total = int64(len(note2))
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
			"msg":    "ok",
			"error":  "",
		})
	}
}

// UpdateByOneKey
// @Description 更新一条标题为title的记录，状态改为status
// @Summary 更新一条记录
// @Accept path
// @Produce json
// @Param status path int true "修改状态"
// @Param title path string true "被更新记录的标题"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo/{status}/{title} [PUT]
func UpdateByOneKey(ctx *gin.Context) {
	var note3 Model.Note
	var data Model.Data
	title := ctx.Param("title")
	status := ctx.Param("status")

	status1, _ := strconv.Atoi(status)
	d := Model.Db1.Where("title=?", title).First(&note3)
	if d.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "修改失败，该条记录不存在",
			"error":  d.Error,
		})
		return
	}

	if note3.DoneStatus == 0 && status1 == 1 {
		Model.Db1.Model(&note3).Update("end_time", time.Now().Unix())
	} else if note3.DoneStatus == 1 && status1 == 0 {
		Model.Db1.Model(&note3).Update("end_time", 0)
	}
	d2 := Model.Db1.Model(&note3).Update("done_status", status1)
	if d2.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "修改发生错误",
			"error":  d2.Error,
		})
		return
	}
	num := d.RowsAffected
	data.Item = append(data.Item, note3)
	data.Total = num
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"msg":    "ok",
		"error":  "",
	})
}

// UpdateAll
// @Description 更新记录
// @Summary 更新所有未完成或已完成记录的完成状态
// @Accept query
// @Produce application/json
// @Param status query int false "完成状态"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo [PUT]
func UpdateAll(ctx *gin.Context) {
	var note3 []Model.Note
	var data Model.Data
	status, ok := ctx.GetQuery("status") //  /todo?status=
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "未输入要修改后的状态",
			"err":    "未输入要修改后的状态",
		})
		return
	}

	status1, _ := strconv.Atoi(status)
	d := Model.Db1.Find(&note3)
	if d.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "修改失败，无记录",
			"error":  d.Error,
		})
		return
	}
	for _, value := range note3 {
		if value.DoneStatus == 0 && status1 == 1 {
			Model.Db1.Model(&note3).Update("end_time", time.Now().Unix())
		} else if value.DoneStatus == 1 && status1 == 0 {
			Model.Db1.Model(&note3).Update("end_time", 0)
		}
	}

	d2 := Model.Db1.Model(&note3).Update("done_status", status1)
	if d2.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "修改发生错误1",
			"error":  d2.Error,
		})
		return
	}

	data.Item = append(data.Item, note3...)
	data.Total = d.RowsAffected
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"msg":    "ok",
		"error":  "",
	})

}

// DeleteByKey
// @Description 删除标题为title的记录
// @Summary 删除一条记录
// @Accept path
// @Produce application/json
// @Param title path string true "被删除记录的标题"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo/{title} [DELETE]
func DeleteByKey(ctx *gin.Context) {
	var note4 Model.Note
	var data Model.Data
	title := ctx.Param("title")

	db1 := Model.Db1.Where("title=?", title).First(&note4)
	if db1.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   Model.Data{},
			"msg":    "删除失败，该条记录不存在",
			"error":  db1.Error,
		})
		return
	}
	db := Model.Db1.Delete(&note4)
	if db.Error != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"data":   data,
			"msg":    "删除错误",
			"error":  db.Error,
		})
		return
	}
	data.Item = append(data.Item, note4)
	data.Total = db.RowsAffected
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"msg":    "删除成功",
		"err":    "",
	})
}

// DeleteAll
// @Description 删除记录
// @Summary 删除所有/已完成/未完成记录
// @Accept query
// @Produce application/json
// @Param status query int false "完成状态"
// @Param Authorization header string false "Bearer 用户令牌"
// @Success 200 {object} ResponseMessage
// @Failure 500 {object} ResponseMessage
// @Router /todo [DELETE]
func DeleteAll(ctx *gin.Context) {
	var note5 []Model.Note
	var data Model.Data
	status, ok := ctx.GetQuery("status")

	if !ok {
		//查找数据库删除controller
		db2 := Model.Db1.Find(&note5)
		if db2.RowsAffected == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   Model.Data{},
				"msg":    "删除失败，无记录存在",
				"error":  db2.Error,
			})
			return
		}
		d := Model.Db1.Delete(&note5)
		if d.Error != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   data,
				"msg":    "删除错误",
				"error":  d.Error,
			})
			return
		}
		//返回json
		data.Item = append(data.Item, note5...)
		data.Total = d.RowsAffected
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
			"msg":    "删除成功",
			"error":  d.Error,
		})

	} else {
		status1, _ := strconv.Atoi(status)
		db2 := Model.Db1.Where("done_status=?", status1).Find(&note5)
		if db2.RowsAffected == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   Model.Data{},
				"msg":    "删除失败，无记录存在",
				"error":  db2.Error,
			})
			return
		}
		d := Model.Db1.Delete(&note5)
		if d.Error != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"status": 500,
				"data":   data,
				"msg":    "删除错误",
				"error":  d.Error,
			})
			return
		}
		data.Item = append(data.Item, note5...)
		data.Total = d.RowsAffected
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   data,
			"msg":    "删除成功",
			"error":  d.Error,
		})
	}
}
