package route

import (
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/handler"
	"github.com/CRobinDev/BCCGembira_Nusastra/internal/middleware"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	App               *fiber.App
	UserHandler       handler.IUserHandler
	PaymentHandler    handler.IPaymentHandler
	ChatHandler       handler.IChatHandler
	ConnectionHandler handler.IConnectionHandler
	Jwt               jwt.IJWT
}

func (c *Config) Register() {
	api := c.App.Group("/api/v1")

	api.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	c.userRoutes(api)
	c.chatRoutes(api)
	c.paymentRoutes(api)
	c.connectionRoutes(api)
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

func (c *Config) chatRoutes(r fiber.Router) {
	chat := r.Group("/chats")
	chat.Post("/create-chat", middleware.Authenticate(c.Jwt), c.ChatHandler.CreateChat())
	chat.Post("/create-chat-ocr", middleware.Authenticate(c.Jwt), c.ChatHandler.CreateChatWithOCR())
	chat.Get("/:id", middleware.Authenticate(c.Jwt), c.ChatHandler.GetChatByID())
}

func (c *Config) paymentRoutes(r fiber.Router) {
	payment := r.Group("/payments")
	payment.Post("/update-status", c.PaymentHandler.UpdatePaymentStatus())
	payment.Post("/create-payment", middleware.Authenticate(c.Jwt), c.PaymentHandler.CreatePayment())
}

func (c *Config) connectionRoutes(r fiber.Router) {
	connection := r.Group("/connections")
	connection.Post("/create-connection", middleware.Authenticate(c.Jwt), c.ConnectionHandler.CreateConnection())
	connection.Delete("/delete-connection", middleware.Authenticate(c.Jwt), c.ConnectionHandler.DeleteConnection())
	connection.Get("/all-connections", middleware.Authenticate(c.Jwt), c.ConnectionHandler.GetAllConnection())
}