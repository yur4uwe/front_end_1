package handlers

import (
	"encoding/json"
	"log"
	"os"
	"fmt"

	hashing "fr_lab_1/pkg/hashing"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserName  string `json:"username"`
	Phone     string `json:"phone"`
	BDay      string `json:"bday"`
	Role      string `json:"role"`
	Gender    string `json:"gender"`
	Agreement bool   `json:"agreement"`
	Photo     string `json:"photo"`
	Contry    string `json:"contry"`
}

func getNextUserID() (int, error) {
	config, err := os.ReadFile("../data/conf.json")
	if err != nil {
		log.Println("Error reading config file:", err)
		return -1, err
	}

	var config_data map[string]int

	err = json.Unmarshal(config, &config_data)
	if err != nil {
		log.Println("Error unmarshalling json:", err)
		return -1, err
	}

	config_data["last_usr_id"]++

	go func() {
		marshalled_data, err := json.Marshal(config_data)
		if err != nil {
			log.Println("Error marshalling json:", err)
			return
		}

		err = os.WriteFile("../data/conf.json", marshalled_data, 0644)
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}
	}()

	return config_data["last_usr_id"], nil
}

func getUserData() ([]User, error) {
	users_data, err := os.ReadFile("../data/user_data.json")
	if err != nil {
		log.Println("Error reading json file as byte[]", err)
		return nil, err
	}

	var all_user_data []User

	err = json.Unmarshal(users_data, &all_user_data)
	if err != nil {
		log.Println("Error unmarshalling json", err)
		return nil, err
	}

	return all_user_data, nil
}

func GetAllUsers() ([]User, error) {
	return getUserData()
}

func writeUserData(users []User) error {
	marshalled_data, err := json.Marshal(users)
	if err != nil {
		log.Println("Error marshalling json:", err)
		return err
	}

	err = os.WriteFile("../data/user_data.json", marshalled_data, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func NewUser(email, password, username, phone, bday, role, gender, photo, contry, authority string, agreement bool) (User, error) {
	id, err := getNextUserID()
	if err != nil {
		log.Println("Error getting next user id:", err)
		return User{}, err
	}

	new_user := User{
		ID:        id,
		Email:     email,
		Password:  password,
		UserName:  username,
		Phone:     phone,
		BDay:      bday,
		Role:      role,
		Gender:    gender,
		Agreement: agreement,
		Photo:     photo,
		Contry:    contry,
	}

	return new_user, nil
}

func AddUser(user User) error {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return err
	}

	user_data = append(user_data, user)

	return writeUserData(user_data)
}

func CheckUserExists(email string) bool {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return false
	}

	for i := 0; i < len(user_data); i++ {
		if user_data[i].Email == email {
			return true
		}
	}

	return false
}

func GetUser(email, password string) (User, error) {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return User{}, err
	}

	password = hashing.HashPassword(password)

	for i := 0; i < len(user_data); i++ {
		if user_data[i].Email == email && user_data[i].Password == password {
			return user_data[i], nil
		}
	}

	return User{}, fmt.Errorf("No such user")
}

func GetUserByID(id int) (User, error) {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return User{}, err
	}

	for i := 0; i < len(user_data); i++ {
		if user_data[i].ID == id {
			return user_data[i], nil
		}
	}

	return User{}, nil
}

func UpdateUser(ID int, email, new_password, new_username, new_phone, new_bday, new_role, new_gender, new_photo, new_contry string, new_agreement bool) error {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return err
	}

	for i := 0; i < len(user_data); i++ {
		if user_data[i].ID == ID {
			user_data[i].Password = new_password
			user_data[i].UserName = new_username
			user_data[i].Phone = new_phone
			user_data[i].BDay = new_bday
			user_data[i].Role = new_role
			user_data[i].Gender = new_gender
			user_data[i].Photo = new_photo
			user_data[i].Contry = new_contry
			user_data[i].Agreement = new_agreement
			break
		}
	}

	return writeUserData(user_data)
}

func DeleteUser(email string) error {
	user_data, err := getUserData()
	if err != nil {
		log.Println("Error getting user data:", err)
		return err
	}

	for i := 0; i < len(user_data); i++ {
		if user_data[i].Email == email {
			user_data = append(user_data[:i], user_data[i+1:]...)
			break
		}
	}

	return writeUserData(user_data)
}
