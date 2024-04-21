package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/marvy-O/fastchat/config"
	"github.com/marvy-O/fastchat/database"
	"github.com/marvy-O/fastchat/users"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	// Load environment variables
	_ = config.LoadEnv()

	//load database
	database.InitDatabase()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api/user")

	api.Post("register", users.Register_user)
	api.Post("login", users.Login_user)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AppConfig.JWT_SECRET)},
	}))

	api.Post("info", users.Get_info)

	app.Listen(":3000")
	log.Println("User service initiated")
}
