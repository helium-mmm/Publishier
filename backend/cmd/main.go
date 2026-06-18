package main

import (
	"log"
	"net/http"
	"time"

	"github.com/helium-mmm/Publishier/internal/api"
	"github.com/helium-mmm/Publishier/internal/api/middleware"
	"github.com/helium-mmm/Publishier/internal/auth"
	"github.com/helium-mmm/Publishier/internal/config"
	"github.com/helium-mmm/Publishier/internal/crypto"
	"github.com/helium-mmm/Publishier/internal/publishier"
	postgre "github.com/helium-mmm/Publishier/internal/repository/postgres"
	"github.com/helium-mmm/Publishier/internal/repository/postgres/models"
	"github.com/helium-mmm/Publishier/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.SocialAccount{},
		&models.Publication{},
	); err != nil {
		log.Fatal("failed to migrate db:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	encryptor, err := crypto.NewEncryptor(cfg.Crypto.Key)
	if err != nil {
		log.Fatal("failed to init encryptor:", err)
	}

	tokenManager := auth.NewTokenManager(cfg.JWT.Secret, cfg.JWT.TTL)

	httpClient := &http.Client{
		Timeout: cfg.Telegram.Timeout,
	}

	userRepo := postgre.NewUserRepository(db)
	postRepo := postgre.NewPostRepository(db)
	accountRepo := postgre.NewSocialAccountRepository(db)
	publicationRepo := postgre.NewPublicationRepository(db)

	telegramPublishier := publishier.NewTelegramPublishier(httpClient, cfg.Telegram.BaseURL)

	authService := service.NewAuthService(userRepo, tokenManager)
	accountService := service.NewAccountService(accountRepo, encryptor)
	postService := service.NewPostService(postRepo, accountRepo, publicationRepo, accountService, telegramPublishier)

	authHandler := api.NewAuthHandler(authService)
	accountHandler := api.NewAccountHandler(accountService)
	postHandler := api.NewPostHandler(postService)

	authMiddleware := middleware.Auth(tokenManager)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("GET /accounts/telegram", authMiddleware(http.HandlerFunc(accountHandler.GetTelegram)))
	mux.Handle("POST /accounts/telegram", authMiddleware(http.HandlerFunc(accountHandler.ConnectTelegram)))
	mux.Handle("POST /posts", authMiddleware(http.HandlerFunc(postHandler.Create)))
	mux.Handle("GET /posts/{id}", authMiddleware(http.HandlerFunc(postHandler.Get)))
	mux.Handle("POST /posts/{id}/publish", authMiddleware(http.HandlerFunc(postHandler.Publish)))

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("server started on port", cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
