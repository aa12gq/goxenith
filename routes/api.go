package routes

import (
	"github.com/gin-gonic/gin"
	controllers "goxenith/app/http/controllers/api/v1"
	"goxenith/app/http/controllers/api/v1/auth"
	"goxenith/app/http/middlewares"
)

func RegisterAPIRoutes(r *gin.Engine) {

	v1 := r.Group("/v1")
	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		{
			suc := new(auth.SignupController)
			// 判断手机是否已注册
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsPhoneExist)
			// 判断邮箱是否已注册
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsEmailExist)
			// 用户注册
			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)
			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			// 图片验证码，需要加限流
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)
			lgc := new(auth.LoginController)
			// 使用手机号，短信验证码进行登录
			authGroup.POST("/login/using-phone", lgc.LoginByPhone)
			// 支持手机号，Email 和 用户名
			authGroup.POST("/login/using-password", lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", lgc.RefreshToken)

		}
		usersGroup := v1.Group("/users")
		{
			uc := new(controllers.UsersController)
			// 获取当前用户
			usersGroup.GET("", middlewares.AuthJWT(), uc.CurrentUser)
			usersGroup.GET("/:id/articles", uc.ListArticlesForUser)
			// 获取用户信息
			usersGroup.GET("/:id", uc.GetUserInfo)
			usersGroup.PUT("", middlewares.AuthJWT(), uc.UpdateUserInfo)
			usersGroup.PUT("/avatar", middlewares.AuthJWT(), uc.UpdateUserAvatar)
			usersGroup.PUT("/password", middlewares.AuthJWT(), uc.UpdatePassword)
		}

		// 分类
		categoryGroup := v1.Group("/categories")
		{
			cate := new(controllers.CategoryController)
			categoryGroup.GET("", cate.ListCategory)
			categoryGroup.GET("/:id", cate.GetCategory)
			categoryGroup.POST("", middlewares.AuthJWT(), cate.CreateCategory)
			categoryGroup.GET("/tree", cate.GetMaterialCategoryTree)
		}
		// 博文
		articleGroup := v1.Group("articles")
		{
			article := new(controllers.ArticleController)
			articleGroup.GET("", article.ListArticle)
			articleGroup.GET("/:id", article.GetArticle)

			articleGroup.DELETE("/:id", middlewares.AuthJWT(), article.DeleteArticle)
			articleGroup.PUT("/update", middlewares.AuthJWT(), article.UpdateArticle)
			articleGroup.POST("/create", middlewares.AuthJWT(), article.CreateArticle)
			articleGroup.POST("/view", article.ViewArticle)
			articleGroup.POST("/like", middlewares.AuthJWT(), article.LikeArticle)
			articleGroup.GET("/:id/check-like-status", middlewares.AuthJWT(), article.CheckLikeStatus)
			articleGroup.POST("/:id/collect", middlewares.AuthJWT(), article.ToggleCollectArticle)
			articleGroup.GET("/collected", middlewares.AuthJWT(), article.GetCollectedArticles)
		}
		commentGroup := v1.Group("comments")
		{
			cmt := new(controllers.CommentController)
			commentGroup.POST("/add", middlewares.AuthJWT(), cmt.AddComment)
			commentGroup.DELETE("/:id", middlewares.AuthJWT(), cmt.DeleteComment)
			commentGroup.GET("/article/:articleId", cmt.GetComments)
			commentGroup.GET("/child/:parentId", cmt.GetChildComments)
			commentGroup.GET("/tree/:articleId", cmt.GetFullCommentTree)
		}
		imc := new(controllers.ImageController)
		imcGroup := v1.Group("/upload")
		{
			imcGroup.POST("", imc.Upload)
		}
		lsc := new(controllers.LinksController)
		linksGroup := v1.Group("/links")
		{
			linksGroup.GET("", lsc.Index)
		}
	}
}
