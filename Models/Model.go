package models

import (
	"github.com/jinzhu/gorm"
)

var (
	Db1 *gorm.DB
)

type User struct {
	ID       int64  `json:"id" swaggertype:"integer"`                      //用户id
	Username string `json:"username" form:"username" swaggertype:"string"` //用户名
	Email    string `json:"email" form:"email" swaggertype:"string"`       //用户邮箱
	Password string `json:"password" form:"password" swaggertype:"string"` //用户密码
}
type Note struct {
	ID          int64  `json:"id" swaggertype:"integer"`
	Title       string `json:"title" swaggertype:"string"`                        //记录题目
	Content     string `json:"content" swaggertype:"string"`                      //记录内容
	DoneStatus  int    `json:"donestatus" gorm:"default:0" swaggertype:"integer"` //完成状态
	Create_time int64  `json:"create_time" swaggertype:"integer"`
	End_time    int64  `json:"end_time" gorm:"default:0" swaggertype:"integer"` //截止时间
	Start_time  int64  `json:"start_time" swaggertype:"integer"`
}
type Data struct {
	Item  []*Note //数据信息
	Total int64   `json:"total" swaggertype:"integer"` //数据数目

}
type ResponseMessage struct {
	Status      int    `json:"status" swaggertype:"integer"` //状态码
	Message     string `json:"msg" swaggertype:"string"`     //响应信息
	Data_Detail *Data  `json:"data"`                         //数据信息
	Error       string `json:"err" swaggertype:"string"`     //错误信息
}
