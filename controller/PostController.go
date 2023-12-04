package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"jkdev.cn/api/common"
	"jkdev.cn/api/model"
	"jkdev.cn/api/response"
	"jkdev.cn/api/vo"
	"log"
	"strconv"
)

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

/*
增删查改操作
*/

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	//把 request 里的参数，绑定到 CreatePostRequest 上面
	if err := ctx.ShouldBind(&requestPost); err != nil {
		//如果出错了，打印日志
		log.Print(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	//像 gorm 这种 orm 框架，java 的 类直接对应 表结构，类似于 java 里的 Hibernate
	//mybatis 好像算是一个 半 orm，没有严格的 对应关系
	//有严格对应关系的好处就是，某些 sql 操作可以很方便的完成，比如 下面的 插入操作，不需要写 sql，直接就能完成
	// 创建post
	post := model.Post{
		UserId:     user.(model.User).ID,
		CategoryId: requestPost.CategoryId,
		Title:      requestPost.Title,
		HeadImg:    requestPost.HeadImg,
		Content:    requestPost.Content,
	}

	// 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
	}

	// 成功
	response.Success(ctx, nil, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, "数据验证错误", nil)
		return
	}

	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "文章属于您，请勿非法操作", nil)
		return
	}

	// 更新文章
	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, "更新失败", nil)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(ctx *gin.Context) {
	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Preload("Category").Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	response.Success(ctx, gin.H{"post": post}, "成功")
}

func (p PostController) Delete(ctx *gin.Context) {
	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, "文章不存在", nil)
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != post.UserId {
		response.Fail(ctx, "文章属于您，请勿非法操作", nil)
		return
	}

	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"post": post}, "成功")
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	// strconv.Atoi 用来把 字符串 转为 int
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// 记录的总条数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	// 返回数据
	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")
}

func NewPostController() IPostController {
	//获取 db，main 刚开始时，已经初始化过了
	db := common.GetDB()
	/*
		GORM 的 AutoMigrate() 方法用于自动迁移 ORM 的 Schemas。所谓 “迁移” 就是刷新数据库中的表格定义，使其保持最新（只增不减）。
		AutoMigrate 会创建（新的）表、缺少的外键、约束、列和索引，并且会更改现有列的类型（如果其大小、精度、是否为空可更改的话）。但不会删除未使用的列，以保护现存的数据。
		可以创建表
	*/
	//如果表不存在，就根据 Post 对象，创建表，如果存在，就检查是否需要添加字段
	db.AutoMigrate(model.Post{})

	return PostController{DB: db}
}
