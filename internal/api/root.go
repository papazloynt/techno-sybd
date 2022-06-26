package api

import (
	"SYBD/internal/api/controllers"
	"SYBD/internal/db"
	"SYBD/internal/service"
	"context"
	"github.com/spf13/viper"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type APIService struct {
	log    *logrus.Entry
	router *echo.Echo
}

func (svc *APIService) Serve() {
	//svc.log.Info("Starting HTTP server")
	listenAddr := viper.GetString("service.bind.address") + ":" + viper.GetString("service.bind.port")
	svc.log.Fatal(svc.router.Start(listenAddr))
}

func (svc *APIService) Shutdown(ctx context.Context) error {
	if err := svc.router.Shutdown(ctx); err != nil {
		svc.log.Fatal(err)
	}
	return nil
}

func NewAPIService(log *logrus.Entry, db_ *pgxpool.Pool) (*APIService, error) {
	svc := &APIService{
		log:    log,
		router: echo.New(),
	}

	//svc.router.Validator = NewValidator()
	//svc.router.Binder = NewBinder()

	repository, err := db.NewRepository(db_)
	if err != nil {
		log.Fatal(err)
	}

	registry := service.NewRegistry(log, repository)

	userCtrl := controllers.NewUserController(log, registry)
	forumCtrl := controllers.NewForumController(log, registry)
	threadCtrl := controllers.NewThreadController(log, registry)
	postCtrl := controllers.NewPostController(log, registry)
	serviceCtrl := controllers.NewServiceController(log, repository)

	api := svc.router.Group("/api")

	api.POST("/user/:nickname/create", userCtrl.CreateUser)
	api.GET("/user/:nickname/profile", userCtrl.GetProfile)
	api.POST("/user/:nickname/profile", userCtrl.UpdateProfile)

	api.POST("/forum/create", forumCtrl.CreateForum)
	api.GET("/forum/:slug/details", forumCtrl.GetForum)
	api.GET("/forum/:slug/threads", forumCtrl.GetForumThreads)
	api.POST("/forum/:slug/create", threadCtrl.CreateThread)
	api.GET("/forum/:slug/users", forumCtrl.GetUsers)

	api.POST("/thread/:slug_or_id/create", postCtrl.CreatePost)
	api.POST("/thread/:slug_or_id/vote", threadCtrl.UpdateVote)
	api.GET("/thread/:slug_or_id/details", threadCtrl.GetDetails)
	api.GET("/thread/:slug_or_id/posts", postCtrl.GetPost)
	api.POST("/thread/:slug_or_id/details", threadCtrl.UpdateForumThread)

	api.GET("/post/:id/details", postCtrl.GetPostDetails)
	api.POST("/post/:id/details", postCtrl.UpdatePost)

	api.GET("/service/status", serviceCtrl.Status)
	api.POST("/service/clear", serviceCtrl.Delete)

	return svc, nil
}
