package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/websockets/controller"
	"github.com/shaileshhb/websockets/db"
	"github.com/shaileshhb/websockets/service"
)

// Server Struct For Start the equisplit service.
type Server struct {
	App      *fiber.App
	Router   fiber.Router
	Database *db.Database
	Log      zerolog.Logger
}

// RegisterRoutes will be implemented by routes package methods to register their routes
type RegisterRoutes interface {
	RegisterRoute(router fiber.Router)
}

// NewServer will initialize the server with logger and fiber router.
func NewServer(log zerolog.Logger, database *db.Database) *Server {
	return &Server{
		Database: database,
		Log:      log,
	}
}

func (ser *Server) InitializeRouter() {
	app := fiber.New(fiber.Config{
		AppName: "Quiz App",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Hello world!!",
		})
	})

	ser.App = app
	ser.Router = app.Group("api/v1")
}

// Register will register all routes from given routes slice.
func (ser *Server) register(routes []RegisterRoutes) {
	for _, route := range routes {
		route.RegisterRoute(ser.Router)
	}
}

func (ser *Server) RegisterModuleRoutes() {
	listingService := service.NewListingService(ser.Database)
	listingController := controller.NewListingController(listingService, ser.Log)

	ser.register([]RegisterRoutes{listingController})
}
