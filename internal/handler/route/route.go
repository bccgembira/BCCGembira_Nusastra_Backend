package route

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/handler"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/middleware"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	App           *fiber.App
	UserHandler   handler.IUserHandler
	GoogleHandler handler.IGoogleHander
	ChatHandler   handler.IChatHandler
	Jwt           jwt.IJWT
}

func (c *Config) Register() {
	api := c.App.Group("/api/v1")

	api.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	c.userRoutes(api)
	c.googleRoutes(api)
	c.chatRoutes(api)
}

func (c *Config) userRoutes(r fiber.Router) {
	user := r.Group("/users")
	user.Post("/register", c.UserHandler.Register())
	user.Post("/login", c.UserHandler.Login())
	user.Get("/me", middleware.Authenticate(c.Jwt), c.UserHandler.GetUser())
	user.Patch("/update-account", middleware.Authenticate(c.Jwt), c.UserHandler.Update())
	user.Delete("/delete-account", middleware.Authenticate(c.Jwt), c.UserHandler.Delete())
	user.Patch("/upload-image", middleware.Authenticate(c.Jwt), c.UserHandler.UploadProfileImage())
	user.Post("/notification", middleware.Authenticate(c.Jwt), c.UserHandler.Notification())
}

func (c *Config) googleRoutes(r fiber.Router) {
	google := r.Group("/oauth")
	google.Get("/redirect", c.GoogleHandler.GoogleLoginRedirect())
	google.Get("/callback", c.GoogleHandler.GoogleCallback())
	google.Post("/login", c.GoogleHandler.GoogleLogin())
}

func (c *Config) chatRoutes(r fiber.Router) {
	chat := r.Group("/chats")
	chat.Post("/create-chat", middleware.Authenticate(c.Jwt), c.ChatHandler.CreateChat())
	chat.Get("/:id", middleware.Authenticate(c.Jwt), c.ChatHandler.GetChatByID())
}