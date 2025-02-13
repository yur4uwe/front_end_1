package token

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	user "fr_lab_1/pkg/user"
)

type ActiveUser struct {
	ID         int    `json:"id"`
	Token      string `json:"token"`
	LastAccess string `json:"last_access"`
}

func getActiveUsers() ([]ActiveUser, error) {
	users_data, err := os.ReadFile("../data/tokens.json")
	if err != nil {
		log.Println("Error reading json file as byte[]", err)
		return nil, err
	}

	var all_user_data []ActiveUser

	err = json.Unmarshal(users_data, &all_user_data)
	if err != nil {
		log.Println("Error unmarshalling json", err)
		return nil, err
	}

	return all_user_data, nil
}

func writeActiveUsers(users []ActiveUser) error {
	marshalled_data, err := json.Marshal(users)
	if err != nil {
		log.Println("Error marshalling json:", err)
		return err
	}

	err = os.WriteFile("../data/tokens.json", marshalled_data, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func CheckTokenExists(token string) bool {
	users, err := getActiveUsers()
	if err != nil {
		return false
	}

	for _, user := range users {
		lastAccessTime, err := time.Parse(time.RFC3339, user.LastAccess)
		if err != nil {
			log.Println("Error parsing time:", err)
			continue
		}
		if user.Token == token && time.Since(lastAccessTime).Hours() < 2 {
			return true
		}
	}

	return false
}

func NewActiveUser(id int, token string) ActiveUser {
	return ActiveUser{id, token, time.Now().Format(time.RFC3339)}
}

func AddActiveUser(id int, token string) error {
	users, err := getActiveUsers()
	if err != nil {
		return err
	}

	users = append(users, ActiveUser{id, token, time.Now().Format(time.RFC3339)})

	err = writeActiveUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func RemoveActiveUser(token string) error {
	users, err := getActiveUsers()
	if err != nil {
		return err
	}

	for i, user := range users {
		if user.Token == token {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	err = writeActiveUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func GenerateToken(id int, email string, password string) string {
	hash := sha256.New()
	hash.Write([]byte(strconv.Itoa(id) + email + password))
	return hex.EncodeToString(hash.Sum(nil))
}

func GetActiveUser(token string) user.User {
	users, err := getActiveUsers()
	if err != nil {
		return user.User{}
	}

	for _, user_it := range users {
		if user_it.Token == token {
			found_user, usr_err := user.GetUserByID(user_it.ID)
			if usr_err != nil {
				return user.User{}
			}

			return found_user
		}
	}

	return user.User{}

}
