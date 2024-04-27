package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/marvy-O/fastchat/database"
)

type User_Info struct {
	Id         uuid.UUID `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
}

type User_credentials struct {
	user_id         string
	hashed_password string
}

func get_user_info(key string, value string) (User_Info, error) {

	query := "SELECT id, first_name, last_name, email, created_at FROM users where %s='%s'"
	query = fmt.Sprintf(query, key, value)

	rows, err := database.ExecuteQuery(query)
	if err != nil {
		return User_Info{}, err
	}
	defer rows.Close()

	var userInfo User_Info
	found := false

	for rows.Next() {
		found = true
		err := rows.Scan(&userInfo.Id, &userInfo.First_name, &userInfo.Last_name, &userInfo.Email, &userInfo.Created_at)
		if err != nil {
			return User_Info{}, err
		}
	}

	if !found {
		return User_Info{}, errors.New("user not found")
	}

	if err := rows.Err(); err != nil {
		return User_Info{}, err
	}

	return userInfo, nil
}

func create_user(new_user User_register) error {
	query := "INSERT INTO users (email, password, first_name, last_name) VALUES ('%s', '%s', '%s', '%s');"
	query = fmt.Sprintf(query, new_user.Email, new_user.Password, new_user.First_name, new_user.Last_name)

	_, err := database.ExecuteQuery(query)

	return err
}

func login_user(email string) (User_credentials, error) {
	query := "SELECT password, id FROM users WHERE email='%s';"
	query = fmt.Sprintf(query, email)

	rows, err := database.ExecuteQuery(query)
	if err != nil {
		return User_credentials{}, err
	}

	var user_credentials User_credentials

	defer rows.Close()

	// Iterate over the rows
	for rows.Next() {
		// Scan the current row's values into variables
		err := rows.Scan(&user_credentials.hashed_password, &user_credentials.user_id)
		if err != nil {
			return User_credentials{}, err
		}
	}

	return user_credentials, nil

}
