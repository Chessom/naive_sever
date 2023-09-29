package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	found         = 0
	notfound      = 1
	unknown_error = 2
)
const (
	done              = 0
	empty             = 1
	non_existent      = 2
	incorrect_pswd    = 3
	invalid_token     = 4
	repeated_checkin  = 5
	repeated_username = 6
	internal_error    = 7
)

type User struct {
	gorm.Model
	Username        string
	Password        string
	Access_token    string
	Point           uint
	LastCheckinTime time.Time
}

func RandomToken() string { //使用随机的uuid作为access token
	u, _ := uuid.NewRandom()
	return u.String()
}

func queryname(db *gorm.DB, name string) (User, int) {
	entry := User{}
	if err := db.Where("username = ?", name).First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //找不到username
			logrus.Error("User:" + name + " not found")
			return User{}, notfound
		} else {
			logrus.Error("User:" + name + " " + err.Error())
			return User{}, unknown_error
		}
	}
	return entry, found
}

func default_handler(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": done,
		"msg":  "我要进求是潮！",
		"data": nil,
	})

}

func ping_handler(c *gin.Context) { //   /ping 的回调
	c.JSON(200, gin.H{
		"code": done,
		"msg":  "pong!",
		"data": nil,
	})
}

func signin_handler(c *gin.Context, db *gorm.DB) { //登录操作的回调
	//获取数据
	type form struct {
		Uname string `json:"username"`
		Pswd  string `json:"password"`
	}
	Form := form{}
	c.BindJSON(&Form)
	uname := Form.Uname
	pswd := Form.Pswd

	if uname == "" || pswd == "" {
		c.JSON(http.StatusOK, gin.H{ //用户名或者密码为空
			"code": empty,
			"msg":  "Empty username or password",
			"data": nil,
		})
		return
	}

	entry, code := queryname(db, uname) //查询记录

	if code == found {
		if entry.Password == pswd { //请求成功，返回token
			c.JSON(http.StatusOK, gin.H{
				"code":         done,
				"msg":          "",
				"data":         nil,
				"access_token": entry.Access_token,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{ //密码有误
				"code": incorrect_pswd,
				"msg":  "Incorrect password",
				"data": nil,
			})
		}
	} else if code == notfound { //查无此人
		c.JSON(http.StatusOK, gin.H{
			"code": non_existent,
			"msg":  "User not found",
			"data": nil,
		})
	} else { //内部错误
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": internal_error,
			"msg":  "Internal Error",
			"data": nil,
		})
	}
}

func checkin_handler(c *gin.Context, db *gorm.DB) { //签到回调
	type Token struct {
		Access_token string `json:"access_token"`
	}
	t := Token{}
	c.BindJSON(&t)
	token := t.Access_token
	entry := User{}
	result := db.Where("access_token = ?", token) //依据token查找
	if err := result.First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //如果找不到
			c.JSON(http.StatusOK, gin.H{
				"code": non_existent,
				"msg":  "invalid token or not found",
				"data": nil,
			})
		} else {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": internal_error,
				"msg":  "Internal Error",
				"data": nil,
			})
		}
	} else { //检查上次签到时间：年月日
		t := entry.LastCheckinTime
		y, m, d := t.Year(), t.Month(), t.Day()
		tnow := time.Now()
		ny, nm, nd := tnow.Year(), tnow.Month(), tnow.Day()
		if nd == d && nm == m && ny == y { //年月日相同，则已签到过
			c.JSON(http.StatusOK, gin.H{
				"code": repeated_checkin,
				"msg":  "Multiple check-in forbidden",
				"data": nil,
			})
		} else { //年月日不同，则签到，更新最后一次签到的时间，更新点数
			entry.Point += 1
			entry.LastCheckinTime = time.Now()
			result.Save(&entry)
			c.JSON(http.StatusOK, gin.H{
				"code":  done,
				"msg":   "",
				"data":  nil,
				"point": entry.Point,
			})
		}
	}
}

func signup_handler(c *gin.Context, db *gorm.DB) { //注册的回调
	type form struct {
		Uname string `json:"username"`
		Pswd  string `json:"password"`
	}
	Form := form{}
	c.BindJSON(&Form)
	uname := Form.Uname
	pswd := Form.Pswd

	if uname == "" || pswd == "" { //用户名为空或者密码为空
		c.JSON(http.StatusOK, gin.H{
			"code": empty,
			"msg":  "Empty user name or password",
			"data": nil,
		})
		return
	}

	_, code := queryname(db, uname) //查询是否存在重名

	if code == found {
		c.JSON(http.StatusOK, gin.H{
			"code": repeated_username,
			"msg":  "User name conflict", //名字冲突
			"data": nil,
		})
	} else {
		token := RandomToken()
		db.Save(&User{Username: uname, Password: pswd, Access_token: token}) //保存随机的token
		c.JSON(http.StatusOK, gin.H{
			"code":         done,
			"msg":          "",
			"data":         nil,
			"access_token": token,
		})
	}
}

type Config struct {
	DB   string `json:"db"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	default_config := Config{DB: "main.db", Host: "127.0.0.1", Port: 8080}
	config := default_config
	con, err := os.ReadFile("config.json")
	if err == nil {
		err = json.Unmarshal([]byte(string(con)), &config)
		if err != nil {
			logrus.Error(err.Error())
			config = default_config
		}
	}
	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{})

	if err != nil {
		logrus.Panic("failed to connect database")
	}
	db.AutoMigrate(&User{})

	r := gin.Default()

	r.GET("/", default_handler)

	r.GET("/ping", ping_handler)

	r.POST("/signin", func(c *gin.Context) {
		signin_handler(c, db)
	})

	r.POST("/signup", func(c *gin.Context) {
		signup_handler(c, db)
	})

	r.POST("/checkin", func(c *gin.Context) {
		checkin_handler(c, db)
	})

	r.Run(":" + fmt.Sprint(config.Port))
}
