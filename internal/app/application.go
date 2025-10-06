package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"laba6/docs"
	"laba6/internal/handlers"
	"laba6/internal/routes"
	"net/http"
	"time"

	"laba6/internal/processors"
	"laba6/internal/repositories"
	"laba6/pkg/config"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Application struct {
	config   *config.Configuration
	server   *http.Server
	database *sqlx.DB
}

func NewApplication(ctx context.Context, cnfg *config.Configuration) (*Application, error) {
	db, err := config.NewPostgresDB(cnfg.Database)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	engine := gin.Default()

	repos := repositories.NewRepositories(db)
	procs := processors.NewProcessors(repos)
	handler := handlers.NewHandler(procs)

	router := routes.NewRouter(engine)
	router.SetupRoutes(handler)

	docs.SwaggerInfo.BasePath = "/"
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:         ":" + cnfg.Application.Port,
		Handler:      engine,
		ReadTimeout:  time.Duration(cnfg.Application.RequestTimeout) * time.Second,
		WriteTimeout: time.Duration(cnfg.Application.ResponseTimeout) * time.Second,
	}

	return &Application{
		config:   cnfg,
		server:   server,
		database: db,
	}, nil
}

func (app *Application) Start() error {
	fmt.Printf("Server is running on port %s\n", app.config.Application.Port)
	return app.server.ListenAndServe()
}

func (app *Application) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down server...")

	if err := app.server.Shutdown(ctx); err != nil {
		return err
	}

	if app.database != nil {
		if err := app.database.Close(); err != nil {
			return err
		}
	}

	fmt.Println("Application shutdown completed")
	return nil
}
