package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMySQL() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/ler_inventory?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return
}

func main() {
	// 创建数据库
	// sql
	// 连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	// 模型绑定
	DB.AutoMigrate(&Todo{})
	// 创建服务
	r := gin.Default()
	// 静态资源
	r.Static("/static", "./static")
	// 静态模板
	r.LoadHTMLGlob("templates/*")
	// 主页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	// 路由
	v1Group := r.Group("v1")
	{
		// 查看所有
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			err = DB.Find(&todoList).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		// id查询
		v1Group.GET("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")
			var todo Todo
			err := DB.Where("id = ?", id).First(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": "id 不存在"})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		// 添加
		v1Group.POST("/todo", func(c *gin.Context) {
			// 从请求中取出数据
			var todo Todo
			c.Bind(&todo)
			// 存入数据库
			// 返回响应
			err = DB.Create(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"data": todo})
			}
		})
		// 修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")
			var todo Todo
			err := DB.Where("id = ?", id).First(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			}
			c.Bind(&todo)
			err = DB.Save(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		// 删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")
			var todo Todo
			err := DB.Where("id = ?", id).Delete(&todo).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
	}
	// 运行
	r.Run()
}
