package users

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/marvy-O/fastchat/config"
	"github.com/marvy-O/fastchat/database"
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

	email := user.Email
	password := user.Password

	if email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields must be provided!"})
	}

	query := "SELECT password, id FROM users WHERE email='%s';"
	query = fmt.Sprintf(query, email)

	rows, err := database.ExecuteQuery(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var hashed_password string
	var user_id string

	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		// Scan the current row's values into variables
		err := rows.Scan(&hashed_password, &user_id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	if hashed_password == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Email not found",
		})
	}

	if check_password(password, hashed_password) {
		//jwt stuff

		// Create the claims
		claims := jwt.MapClaims{
			"user_id": user_id,
			"exp":     time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
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

	first_name := user.First_name
	last_name := user.Last_name
	email := user.Email
	password := user.Password

	hashed_password, err := hash_password(password)
	if err != nil {
		log.Fatal("Error in hashing password: ", err)
	}

	if first_name == "" || email == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields must be provided!"})
	}

	query := "INSERT INTO users (email, password, first_name, last_name) VALUES ('%s', '%s', '%s', '%s');"
	query = fmt.Sprintf(query, email, hashed_password, first_name, last_name)

	_, err = database.ExecuteQuery(query)
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

	query := "SELECT id, first_name, last_name, email, created_at FROM users where id='%s'"
	query = fmt.Sprintf(query, user_id)

	rows, err := database.ExecuteQuery(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var (
		id         string
		first_name string
		last_name  string
		email      string
		created_at string
	)

	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		// Scan the current row's values into variables
		err := rows.Scan(&id, &first_name, &last_name, &email, &created_at)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"user_id":    id,
		"first_name": first_name,
		"last_name":  last_name,
		"email":      email,
		"created_at": created_at,
	})
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
