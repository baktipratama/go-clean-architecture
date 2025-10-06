package main

import (
	"log"

	"go-clean-code/internal/handler"
	"go-clean-code/internal/repository"
	"go-clean-code/internal/usecase"
)

type Container struct {
	UserRepository repository.UserRepositoryInterface
	UserUsecase    usecase.UserUsecaseInterface
	UserHandler    *handler.UserHandler
}

func NewContainer() *Container {
	config := NewConfig()

	// Initialize PostgreSQL connection
	db, err := ConnectDatabase(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Run migrations
	if err := RunMigrations(db, "./migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	log.Println("Using PostgreSQL database")

	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	return &Container{
		UserRepository: userRepo,
		UserUsecase:    userUsecase,
		UserHandler:    userHandler,
	}
}
