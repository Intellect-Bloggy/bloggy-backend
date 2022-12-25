package main

import (
	"github.com/gin-contrib/cors"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/pgtype"

	_ "github.com/Intellect-Bloggy/bloggy-backend/docs"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/adapters/db/postgres"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/controller/http"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/controller/http/middleware"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/domain/service"
	auth_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/auth"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/post"
	user_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user"
	"github.com/Intellect-Bloggy/bloggy-backend/internal/server"
	"github.com/Intellect-Bloggy/bloggy-backend/pkg/auth"
	pq_client "github.com/Intellect-Bloggy/bloggy-backend/pkg/client/postgres"
	"github.com/Intellect-Bloggy/bloggy-backend/pkg/hash"
)

// @title Bloggy-backend
// @version 1.0
// @description backend for Bloggy (open source twitter-like app)

// @host localhost:8000
// @Basepath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_LOCAL_PORT"))
	if dbPort > 1<<16 || dbPort < 1 || err != nil {
		logrus.Fatal("Подключение к базе данных не удалось: некорректный порт")
	}

	db, err := pq_client.NewPostgresDB(pq_client.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     uint16(dbPort),
		Username: os.Getenv("POSTGRES_USER"),
		DBName:   os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		SSLMode:  os.Getenv("POSTGRES_SSLMode"),
	})
	if err != nil {
		logrus.Fatal("Подключение к базе данных не удалось: ", err)
	}

	logrus.Println("База данных успешно подключена")

	// TODO: Переделать на ENV
	// FIXME: Серьезно
	salt := "random string test"
	signKey := "random string text"

	router := gin.New()

	// cors: allow all origins
	router.Use(cors.Default())

	// utils
	hasher := hash.NewSHA1Hasher(salt)
	tManager, err := auth.NewManager(signKey)
	if err != nil {
		logrus.Fatal("Не удалось создать менеджер токенов ", err)
	}

	// middlewares
	authMiddleware := middleware.NewAuthMiddleware(tManager)

	// storages
	userStorage := postgres.NewUserStorage(db)
	authStorage := postgres.NewAuthStorage(db)
	postStorage := postgres.NewPostStorage(db)

	// services
	authService := service.NewAuthService(authStorage, tManager)
	userService := service.NewUserService(userStorage, hasher)
	postService := service.NewPostService(postStorage)

	// usecases
	authUsecase := auth_usecase.NewAuthUsecase(authService, userService)
	userUsecase := user_usecase.NewUserUsecase(userService, postService)
	postUsecase := post.NewPostUsecase(postService)

	// register handlers
	http.NewAuthHandler(authUsecase).Register(router)
	http.NewUserHandler(authMiddleware, userUsecase).Register(router)
	http.NewPostHandler(authMiddleware, postUsecase).Register(router)
	http.NewDocsHandler().Register(router)

	srv := server.NewServer()
	srvPort, err := strconv.Atoi(os.Getenv("PORT"))
	if srvPort > 1<<16 || srvPort < 1 || err != nil {
		logrus.Fatal("Запуск сервера не удался: некорректный порт")
	}

	if err := srv.Run(uint16(srvPort), router); err != nil {
		logrus.Fatal("Ошибка запуска сервера: ", err)
	}
}
