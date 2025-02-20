package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	hashing "fr_lab_1/pkg/hashing"
	token "fr_lab_1/pkg/token"
	user "fr_lab_1/pkg/user"
)

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Register handler reached")

	email := r.FormValue("email")
	password := r.FormValue("password")
	user_name := r.FormValue("username")
	phone := r.FormValue("phone")
	bday := r.FormValue("bday")
	gender := r.FormValue("gender")
	agreement := r.FormValue("agreement")
	contry := r.FormValue("contry")

	// log.Println("Request: ", r)
	log.Println("Registering user with email:", email)
	log.Println("Registering user with password:", password)
	log.Println("Registering user with username:", user_name)
	log.Println("Registering user with phone:", phone)
	log.Println("Registering user with birthday:", bday)
	log.Println("Registering user with gender:", gender)
	log.Println("Registering user with agreement:", agreement)
	log.Println("Registering user with contry:", contry)

	if user_name == "" || phone == "" || bday == "" || gender == "" || agreement == "" || contry == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	password = hashing.HashPassword(password)

	if user.CheckUserExists(email) {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	photo, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error uploading photo", http.StatusBadRequest)
	}
	defer photo.Close()

	uploadsDir := "../static/uploads/"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadsDir, 0755)
		if err != nil {
			http.Error(w, "Error creating uploads directory", http.StatusInternalServerError)
			return
		}
	}

	filePath := uploadsDir + hashing.HashImageName(header.Filename) + ".jpg"
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, photo)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	new_user, err := user.NewUser(email, password, user_name, phone, bday, "user", gender, filePath, contry, "user", agreement == "true")
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	err = user.AddUser(new_user)
	if err != nil {
		http.Error(w, "Error adding user", http.StatusInternalServerError)
		return
	}

	user_token := token.GenerateToken(new_user.ID, new_user.Email, new_user.Password)

	err = token.AddActiveUser(new_user.ID, user_token)
	if err != nil {
		http.Error(w, "Error adding active user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token":     user_token,
		"authority": "user",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler reached")

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := user.GetUser(email, password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	user_token := token.GenerateToken(user.ID, user.Email, user.Password)

	if !token.CheckTokenExists(user_token) {
		err = token.AddActiveUser(user.ID, user_token)
		if err != nil {
			http.Error(w, "Error adding active user", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"authority": user.Role,
		"token":     user_token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("IsAuthorized") == "false" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user_token := r.Header.Get("Authorization")

	if !token.CheckTokenExists(user_token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := token.RemoveActiveUser(user_token)
	if err != nil {
		http.Error(w, "Error removing active user", http.StatusInternalServerError)
		return
	}

	log.Println("Logout handler reached")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func Authority(w http.ResponseWriter, r *http.Request) {
	log.Println("Authority handler reached")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"authority": r.Header.Get("Authority")})
}
