package users

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/marvy-O/fastchat/config"
	"golang.org/x/crypto/bcrypt"
)

type User_register struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type User_login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login_user(c *fiber.Ctx) error {
	user := new(User_login)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields must be provided!"})
	}

	user_credentials, err := login_user(user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user_credentials.hashed_password == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email not found",
		})
	}

	if check_password(user.Password, user_credentials.hashed_password) {
		claims := jwt.MapClaims{
			"user_id": user_credentials.user_id,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(config.AppConfig.JWT_SECRET))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Incorrect email or password"})
}

func Register_user(c *fiber.Ctx) error {
	user := new(User_register)

	// Parse the JSON input from the request body
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	hashed_password, err := hash_password(user.Password)
	if err != nil {
		log.Fatal("Error in hashing password: ", err)
	}

	user.Password = hashed_password

	if user.First_name == "" || user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields must be provided!"})
	}

	err = create_user(*user)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "An account with the provied email already exists!",
		})
	}

	return c.JSON(fiber.Map{"message": "Account created successfully!"})
}

func Get_info(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	user_id := claims["user_id"].(string)

	userInfo, err := get_user_info("id", user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(userInfo)
}

func Get_users_info(c *fiber.Ctx) error {
	email := c.Query("email")
	id := c.Query("id")

	if email == "" && id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "either 'email' or 'id' parameter is required",
		})
	}

	var key, val string
	switch {
	case email != "":
		key, val = "email", email
	case id != "":
		key, val = "id", id
	}

	userInfo, err := get_user_info(key, val)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(userInfo)
}

func hash_password(password string) (string, error) {
	// Generate a hashed version of the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func check_password(password, hashedPassword string) bool {
	// Compare the password with its hashed version using bcrypt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
