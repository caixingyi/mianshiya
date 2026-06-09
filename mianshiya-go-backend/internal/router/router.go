package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mianshiya-go-backend/internal/auth"
	"mianshiya-go-backend/internal/handler"
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

	questionBankRepo := questionbank.NewRepository(database)
	questionBankService := questionbank.NewService(questionBankRepo)
	questionBankHandler := questionbank.NewHandler(questionBankService)

	questionRepo := question.NewRepository(database)
	questionService := question.NewService(questionRepo)
	questionHandler := question.NewHandler(questionService)

	questionBankQuestionRepo := questionbankquestion.NewRepository(database)
	questionBankQuestionService := questionbankquestion.NewService(
		questionBankQuestionRepo,
		questionRepo,
		questionBankRepo,
	)
	questionBankQuestionHandler := questionbankquestion.NewHandler(questionBankQuestionService)

	// 公开接口
	api.POST("/user/register", userHandler.RegisterHandler)
	api.POST("/user/login", userHandler.LoginHandler)
	api.GET("/user/get/vo", userHandler.GetUserVOHandler)
	api.GET("/questionBank/get/vo", questionBankHandler.GetQuestionBankVOHandler)
	api.POST("/questionBank/list/page/vo", questionBankHandler.ListQuestionBankVOHandler)
	api.GET("/question/get/vo", questionHandler.GetQuestionVOHandler)
	api.POST("/question/list/page/vo", questionHandler.ListQuestionVOHandler)
	// 需要认证的接口
	authAPI := api.Group("")
	authAPI.Use(auth.AuthMiddleware(tokenStore))
	authAPI.GET("/user/get/login", userHandler.GetLoginUserHandler)
	authAPI.POST("/user/logout", userHandler.LogoutHandler)
	authAPI.POST("/user/update/my", userHandler.UpdateMyHandler)
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
	adminAPI.POST("/question/add", questionHandler.AddQuestionHandler)
	adminAPI.POST("/question/delete", questionHandler.DeleteQuestionHandler)
	adminAPI.POST("/question/update", questionHandler.UpdateQuestionHandler)
	adminAPI.POST("/questionBankQuestion/add/batch", questionBankQuestionHandler.BatchAddQuestionsToBankHandler)
	adminAPI.POST("/questionBankQuestion/remove/batch", questionBankQuestionHandler.BatchRemoveQuestionsFromBankHandler)
}
