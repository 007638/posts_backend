package main
import (
  "github.com/gin-gonic/gin"
  "strconv"
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
  "fmt"
)


type User struct {
  ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
  Username string `json:"username"`
  Password string `json:"password"`
}

type Post struct {
  ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
  Title string `json:"title"`
  Content string `json:"content"`
  UserID int64 `json:"user_id"`
  User User `gorm:"foreignKey:UserID" json:"user"`
  //一条帖子对应多个评论
  Comments []Comment `gorm:"foreignKey:PostsID" json:"comments"`
}

type Comment struct {
  ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
  Content string `json:"content"`//评论内容
  PostsID int64 `json:"posts_id"`//属于哪条帖子
  ParentID *int64 `json:"parent_id"`//回复评论
  UserID int64 `json:"user_id"`//是谁发的这条评论
  User User `foreignKey:UserID" json:"user"`//评论的作者
}

func main() {
  r := gin.Default()
  //解决前端跨域
  r.Use(func(c *gin.Context) {
    //允许所有前端地址访问，开发环境用
    c.Header("Access-Control-Allow-Origin","*")
    c.Header("Access-Control-Allow-Methods","GET,POST,PUT,DELETE,OPTIONS")
    c.Header("Access-Control-Allow-Headers","Content-Type")
    //处理预检OPTIONS请求
    if c.Request.Method == "OPTIONS" {
      c.AbortWithStatus(204)
      return
    }
    c.Next()
  })
  //数据库连接配置
  dsn := "postgres://postgres:88888888@172.20.13.32:5500/post_db?sslmode=disable&timezone=Asia/Shanghai"
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
    panic("数据库连接失败")
  }
  fmt.Println("ok")
  //自动建表：不存在的表就自动创建，已存的就不动
  db.AutoMigrate(&User{},&Post{},&Comment{})

  //注册接口
  r.POST("/api/register",func(c *gin.Context) {
    var user User//定义变量user，用来存前端传过来的用户名和密码
    err := c.ShouldBindJSON(&user)//把前端传过来的json数据解析到user变量里
    if err != nil {
      c.JSON(400,gin.H{"msg":"参数错误"})
      return
    }
    //判断用户名和密码是否为空
    if user.Username == "" || user.Password == "" {
      c.JSON(400,gin.H{"msg":"用户名和密码不能为空"})
      return
    }
    res := db.Create(&user)//把用户数据插入数据库users表
    if res.Error != nil {//判断插入数据库有没有出错
      c.JSON(500,gin.H{"msg":"注册失败"})
      return
    }
    c.JSON(200,gin.H{"msg":"注册成功", "data": user})
  })
  //登录接口
  r.POST("/api/login",func(c *gin.Context) {
    var input User //定义变量input，用来存前端传过来的登录数据
    err := c.ShouldBindJSON(&input)//把前端穿过来的登录数据解析到变量input中
    if err != nil {
      c.JSON(400, gin.H{"msg":"参数错误"})
      return
    }
    var findUser User//定义变量findUser，用来存从数据库查询出来的用户信息
    err = db.Where("username = ?", input.Username).First(&findUser).Error//根据用户名去数据库查找用户
    if err != nil {
      c.JSON(400, gin.H{"msg":"用户不存在"})
      return
    }
    if findUser.Password != input.Password {
      c.JSON(401,gin.H{"msg":"密码错误"})
      return
    }
    c.JSON(200, gin.H{
      "msg":"登录成功",
      "user":gin.H{"id":findUser.ID,"username":findUser.Username},
    })
  })

  //定义一个GET请求的接口，用于返回帖子列表数据
  r.GET("/api/posts",func(c *gin.Context) {
    var list []Post
    err := db.Find(&list).Error
    if err != nil {
      c.JSON(500,gin.H{"msg":"查询数据库出错"})
      return
    }
    c.JSON(200,gin.H{"result":list})
  })

  //新增帖子接口
  r.POST("/api/create",func(c *gin.Context) {
    var blog Post//定义一个Post结构体变量blog，用来接收签到传过来的博客数据
    err := c.ShouldBindJSON(&blog) //把前端穿过来的数据解析到变量blog中
    fmt.Println("err", err)
    if err != nil {
      c.JSON(400,gin.H{"msg":"标题/内容不能为空"})
      return
    }
    //写入数据库
    err = db.Create(&blog).Error
    if err != nil{
      c.JSON(500, gin.H{"msg":"新增帖子失败"})
      return
    }
    c.JSON(200, gin.H{
      "msg":"新增帖子成功",
      "data":blog,//把创建好的帖子数据返回前端
    })
  })

  //添加查看详情路由
  r.GET("/api/post/:id",func(c *gin.Context){
    idStr := c.Param("id")//从路劲获取前端传过来的博客id
    //将字符串id转成int类型
    targetID,err := strconv.ParseInt(idStr,10,64)
    fmt.Println("targetID",targetID)
    if err != nil {
      c.JSON(400, gin.H{"msg":"ID格式错误，必须为数字"})
      return
    }
    var post Post//定义一个Post结构体变量post，用来接收数据库查询结果
    //Preloda预加载关联数据：帖子对应的作者，帖子下的全部评论
    err = db.Preload("User").Preload("Comments").First(&post,targetID).Error//去数据库查id=targetID的帖子，查到数据放到blog里
    if err != nil {
      c.JSON(404,gin.H{"msg":"找不到这篇帖子"})
      return
    }
    c.JSON(200,post)
  })

  //添加回复评论接口
  r.POST("api/comment",func(c *gin.Context){
    //定义commentReq自定义结构体，用来接收前端传过来的数据
    type commentReq struct {
    PostsID int64 `json:"post_id"` //接收前端传的posts_id,代表这条评论属于哪一篇帖子
    ParentID *int64 `json:"parent_id"`//接收parent_id，回复别人评论就填被回复评论的id，直接评论帖子传null
    Content  string `json:"content"`//接收评论文字内容
    UserID int64 `json:"user_id"`//代表是谁发布的这条评论
  }
  var req commentReq//定义一个commentReq结构体变量是req，用来接收用来存放前端传过来的数据
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	if req.Content == "" {
		c.JSON(400, gin.H{"msg": "评论内容不能为空"})
		return
	}

  //创建Comment结构体变量newComment,用来组件准备插入数据库的数据
  newComment := Comment{
		PostsID:  req.PostsID,
		ParentID: req.ParentID,
		Content:  req.Content,
		UserID:   req.UserID,
	}
  err = db.Create(&newComment).Error //把数据插入数据库里，数据库报错信息将赋值给err
  if err != nil{
    c.JSON(500,gin.H{"msg":"评论保存失败"})
    return
  }
  c.JSON(200,gin.H{
    "msg":"评论成功",
    "data":newComment,//把存入数据库的完整评论数据返回给前端
  })
})

  r.Run("0.0.0.0:8099")
}





























