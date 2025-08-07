package v1

import (
	"apps/config"
	"apps/internal/infrastructure/db"
	userRepository "apps/internal/modules/user/repository"
	userUsecase "apps/internal/modules/user/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Routes(r fiber.Router, cfg *config.Config, db db.MysqlDBInterface) {
	userRepo := userRepository.NewUserRepository(db)
	userUc := userUsecase.NewUserUsecase(userRepo, cfg.App.Timeout*time.Second)
	userHandler := NewUserHandler(userUc)

	r.Post("/register", userHandler.Create)
	r.Post("/login", userHandler.Login)

	userRoute := r.Group("/user")
	userRoute.Get("check-by-email/:email", userHandler.GetByEmail)
	userRoute.Get("check-by-phone/:phone", userHandler.GetByPhone)
	userRoute.Put("update-password/:email", userHandler.Update)
}
