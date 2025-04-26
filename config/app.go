package config

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/handler"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/handler/route"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/repository"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/service"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/claude"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/gomail"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/google"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/log"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/supabase"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AppConfig struct {
	App *fiber.App
	DB  *gorm.DB
}

func StartApp(config *AppConfig) {
	jwt := jwt.Init()
	val := validator.NewValidator()
	gomail := gomail.NewGomail()
	logger := log.NewLogger()
	oauth := google.NewGoogleOAuth()
	// midtrans := midtrans.NewMidtrans()
	supabase := supabase.NewSupabase()
	claude := claude.NewClaude(logger)

	userRepository := repository.NewUserRepository(config.DB, logger)
	userService := service.NewUserService(userRepository, jwt, gomail, supabase, logger)
	userHandler := handler.NewUserHandler(userService, val)

	googleHandler := handler.NewGoogleHandler(val, userService, oauth)

	chatRepository := repository.NewChatRepository(config.DB, logger)
	chatService := service.NewChatService(chatRepository, logger, claude)
	chatHandler := handler.NewChatHandler(chatService, val)
	
	routes := route.Config{
		App:           config.App,
		UserHandler:   userHandler,
		GoogleHandler: googleHandler,
		ChatHandler: chatHandler,
		Jwt:           jwt,
	}

	routes.Register()
}
