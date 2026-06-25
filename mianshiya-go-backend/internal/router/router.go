package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/handler"
	"mianshiya-go-backend/internal/post"
	"mianshiya-go-backend/internal/postfavour"
	"mianshiya-go-backend/internal/postthumb"
	"mianshiya-go-backend/internal/question"
	"mianshiya-go-backend/internal/questionbank"
	"mianshiya-go-backend/internal/questionbankquestion"
	"mianshiya-go-backend/internal/user"
)

func RegisterRouter(r *gin.Engine, database *gorm.DB, tokenStore auth.TokenStore) {
	api := r.Group("/api")

	api.GET("/health", handler.HealthHandler)
	api.GET("/error-demo", handler.ErrorDemoHandler)

	userRepo := user.NewRepository(database)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, tokenStore)

	questionRepo := question.NewRepository(database)
	questionService := question.NewService(questionRepo)
	questionHandler := question.NewHandler(questionService)

	questionBankRepo := questionbank.NewRepository(database)
	questionBankService := questionbank.NewService(questionBankRepo, questionService)
	questionBankHandler := questionbank.NewHandler(questionBankService)

	questionBankQuestionRepo := questionbankquestion.NewRepository(database)
	questionBankQuestionService := questionbankquestion.NewService(
		questionBankQuestionRepo,
		questionRepo,
		questionBankRepo,
	)
	questionBankQuestionHandler := questionbankquestion.NewHandler(questionBankQuestionService)
	postRepo := post.NewRepository(database)
	postService := post.NewService(postRepo, userService)
	postHandler := post.NewHandler(postService)

	postThumbRepo := postthumb.NewRepository(database)
	postThumbService := postthumb.NewService(postThumbRepo, postRepo)
	postThumbHandler := postthumb.NewHandler(postThumbService)

	postFavourRepo := postfavour.NewRepository(database)
	postFavourService := postfavour.NewService(postFavourRepo, postRepo)
	postFavourHandler := postfavour.NewHandler(postFavourService, postService)

	// 公开接口
	api.POST("/user/register", userHandler.RegisterHandler)
	api.POST("/user/login", userHandler.LoginHandler)
	api.GET("/user/get/vo", userHandler.GetUserVOHandler)
	api.GET("/questionBank/get/vo", questionBankHandler.GetQuestionBankVOHandler)
	api.POST("/questionBank/list/page/vo", questionBankHandler.ListQuestionBankVOHandler)
	api.GET("/question/get/vo", questionHandler.GetQuestionVOHandler)
	api.POST("/question/list/page/vo", questionHandler.ListQuestionVOHandler)
	api.POST("/question/search/page/vo", questionHandler.ListQuestionVOHandler)
	api.GET("/questionBankQuestion/get/vo", questionBankQuestionHandler.GetQuestionBankQuestionVOHandler)
	api.POST("/questionBankQuestion/list/page/vo", questionBankQuestionHandler.ListQuestionBankQuestionVOHandler)
	api.GET("/post/get/vo", auth.OptionalAuthMiddleware(tokenStore), postHandler.GetPostVOHandler)
	api.POST("/post/list/page/vo", auth.OptionalAuthMiddleware(tokenStore), postHandler.ListPostVOHandler)
	api.POST("/postFavour/list/page", auth.OptionalAuthMiddleware(tokenStore), postFavourHandler.ListFavourPostHandler)

	// 需要认证的接口
	authAPI := api.Group("")
	authAPI.Use(auth.AuthMiddleware(tokenStore))
	authAPI.GET("/user/get/login", userHandler.GetLoginUserHandler)
	authAPI.POST("/user/logout", userHandler.LogoutHandler)
	authAPI.POST("/user/update/my", userHandler.UpdateMyHandler)
	authAPI.POST("/questionBank/my/list/page/vo", questionBankHandler.ListMyQuestionBankVOHandler)
	authAPI.POST("/questionBank/edit", questionBankHandler.EditQuestionBankHandler)
	authAPI.POST("/question/my/list/page/vo", questionHandler.ListMyQuestionVOHandler)
	authAPI.POST("/question/edit", questionHandler.EditQuestionHandler)
	authAPI.POST("/questionBankQuestion/my/list/page/vo", questionBankQuestionHandler.ListMyQuestionBankQuestionVOHandler)
	authAPI.POST("/post/add", postHandler.AddPostHandler)
	authAPI.POST("/post/delete", postHandler.DeletePostHandler)
	authAPI.POST("/post/edit", postHandler.EditPostHandler)
	authAPI.POST("/post/my/list/page/vo", postHandler.ListMyPostsVOHandler)
	authAPI.POST("/postThumb/do", postThumbHandler.DoPostThumbHandler)
	authAPI.POST("/postFavour/do", postFavourHandler.DoPostFavourHandler)
	authAPI.POST("/postFavour/my/list/page", postFavourHandler.ListMyFavourPostHandler)

	// 管理员接口
	adminAPI := api.Group("")
	adminAPI.Use(auth.AuthMiddleware(tokenStore), user.AdminMiddleware(userService))
	adminAPI.GET("/user/admin/check", userHandler.AdminCheckHandler)
	adminAPI.POST("/user/add", userHandler.AddUserHandler)
	adminAPI.POST("/user/delete", userHandler.DeleteUserHandler)
	adminAPI.POST("/user/update", userHandler.UpdateUserHandler)
	adminAPI.POST("/user/list/page", userHandler.ListUserHandler)
	adminAPI.GET("/user/get", userHandler.GetUserHandler)
	adminAPI.POST("/questionBank/add", questionBankHandler.AddQuestionBankHandler)
	adminAPI.POST("/questionBank/delete", questionBankHandler.DeleteQuestionBankHandler)
	adminAPI.POST("/questionBank/update", questionBankHandler.UpdateQuestionBankHandler)
	adminAPI.POST("/questionBank/list/page", questionBankHandler.ListQuestionBankHandler)
	adminAPI.POST("/question/add", questionHandler.AddQuestionHandler)
	adminAPI.POST("/question/delete", questionHandler.DeleteQuestionHandler)
	adminAPI.POST("/question/update", questionHandler.UpdateQuestionHandler)
	adminAPI.POST("/question/list/page", questionHandler.ListQuestionHandler)
	adminAPI.POST("/question/delete/batch", questionHandler.BatchDeleteQuestionsHandler)
	adminAPI.POST("/questionBankQuestion/add", questionBankQuestionHandler.AddQuestionBankQuestionHandler)
	adminAPI.POST("/questionBankQuestion/delete", questionBankQuestionHandler.DeleteQuestionBankQuestionHandler)
	adminAPI.POST("/questionBankQuestion/update", questionBankQuestionHandler.UpdateQuestionBankQuestionHandler)
	adminAPI.POST("/questionBankQuestion/list/page", questionBankQuestionHandler.ListQuestionBankQuestionHandler)
	adminAPI.POST("/questionBankQuestion/remove", questionBankQuestionHandler.RemoveQuestionBankQuestionHandler)
	adminAPI.POST("/questionBankQuestion/add/batch", questionBankQuestionHandler.BatchAddQuestionsToBankHandler)
	adminAPI.POST("/questionBankQuestion/remove/batch", questionBankQuestionHandler.BatchRemoveQuestionsFromBankHandler)
	adminAPI.POST("/post/update", postHandler.UpdatePostHandler)
	adminAPI.POST("/post/list/page", postHandler.ListPostPageHandler)
}
