package models

import (
	"github.com/jinzhu/gorm"
)

var (
	Db1 *gorm.DB
)

type User struct {
	ID       int64  `json:"id"`                       //用户id
	Username string `json:"username" form:"username"` //用户名
	Email    string `json:"email" form:"email"`       //用户邮箱
	Password string `json:"password" form:"password"` //用户密码
}
type Note struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`                       //记录题目
	Content     string `json:"content"`                     //记录内容
	DoneStatus  int    `json:"donestatus" gorm:"default:0"` //完成状态
	Create_time int64  `json:"create_time"`
	End_time    int64  `json:"end_time" gorm:"default:0"` //截止时间
	Start_time  int64  `json:"start_time"`
}
type Data struct {
	Item  []Note //数据信息
	Total int64  //数据数目
}
type ResponseMessage struct {
	Status      int    `json:"status"` //状态码
	Message     string `json:"msg"`    //响应信息
	Data_Detail *Data  `json:"data"`   //数据信息
	Error       string `json:"error"`  //错误信息
}
