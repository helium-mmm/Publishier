package cmd

import (
	"log"
	"net/http"
	"time"

	"github.com/helium-mmm/Publishier/internal/api"
	"github.com/helium-mmm/Publishier/internal/config"
	"github.com/helium-mmm/Publishier/internal/publishier"
	postgre "github.com/helium-mmm/Publishier/internal/repository/postgres"
	"github.com/helium-mmm/Publishier/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// жоски конфиг
	cfg := config.MustLoad()

	// жоская бдшка
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect db:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// жоски клиент
	httpClient := &http.Client{
		Timeout: cfg.Telegram.Timeout,
	}

	// жоская репо
	postRepo := postgre.NewPostRepository(db)
	accountRepo := postgre.NewSocialAccountREpository(db)
	publicationRepo := postgre.NewPublicationRepository(db)

	// жоски публикатор в телегу
	telegramPublishier := publishier.NewTelegramPublishier(httpClient, cfg.Telegram.BaseURL)
	
	// жоски сервис доставки в бдшку
	postService := service.NewPostService(postRepo, accountRepo, publicationRepo, telegramPublishier)

	// жоские ручки
	postHandler := api.NewPostHandler(postService)

	// жоски роутер
	mux := http.NewServeMux()

	mux.HandleFunc("/posts", postHandler.Create)
	mux.HandleFunc("/posts/publish", postHandler.Publish)
	mux.HandleFunc("/posts/get", postHandler.Get)

	// жоски сервер
	server := &http.Server{
		Addr: ":" + cfg.Server.Port,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	// жоски стаааааарт
	log.Println("server started on port", cfg.Server.Port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}