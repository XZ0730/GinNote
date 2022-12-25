package Jwtutil

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtkey = []byte("loginkey")

// SetToken
// @Description 创建token
// @Summary 创建token
// @Accept json
// @Produce json
// @Param username json string true "用户名"
// @Success 200 {json} json "status":200,"token":"
// @Failure 500 "修改失败"
// @Router /login [POST]
func SetToken(ctx *gin.Context) {
	value, _ := ctx.Get("username")
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.MapClaims{
		"ExpiresAt": expireTime.Unix(), //过期时间
		"IssuedAt":  time.Now().Unix(),
		"Issuer":    "127.0.0.1",  // 签名颁发者
		"Subject":   "user token", //签名主题
		"Username":  value,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token1, err := token.SignedString(jwtkey)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": 500,
			"err":    err.Error,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"token":  token1,
	})
	ctx.Next()
}

// CheckToken
// @Description 检验token
// @Summary  检验token
// @Accept json
// @Produce json
// @Param token json string true "token"
// @Param Authorization header string true "用户令牌"
// @Success 200 {json} json ""msg":"token 验证成功","status":200"
// @Failure 401 "权限不足"
// @Router /todo [ANY]
func CheckToken(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")

	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"msg":  "权限不足1",
		})
		ctx.Abort()
		return
	}
	token, _, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status": 401,
			"msg":  "权限不足2",
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "token 验证成功",
	})
	ctx.Next()
}

// cannot convert token.Claims.(jwt.MapClaims)["ExpiresAt"] (map index expression of type interface{}) to int64 (need type assertion)
func ParseToken(tokenString string) (*jwt.Token, *jwt.MapClaims, error) {
	Claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtkey, nil
	})
	return token, Claims, err
}
