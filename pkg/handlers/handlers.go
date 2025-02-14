package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	hashing "fr_lab_1/pkg/hashing"
	token "fr_lab_1/pkg/token"
	user "fr_lab_1/pkg/user"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home handler reached")
	// http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
	http.ServeFile(w, r, "../static/index.html")
}

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
	// log.Println("Registering user with email:", email)
	// log.Println("Registering user with password:", password)
	// log.Println("Registering user with username:", user_name)
	// log.Println("Registering user with phone:", phone)
	// log.Println("Registering user with birthday:", bday)
	// log.Println("Registering user with gender:", gender)
	// log.Println("Registering user with agreement:", agreement)
	// log.Println("Registering user with contry:", contry)

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
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
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

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("Admin dashboard handler reached")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Admin dashboard"})
}

func UserDashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("User dashboard handler reached")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User dashboard"})
}

func Users(w http.ResponseWriter, r *http.Request) {
	log.Println("Users handler reached")

	user_data, err := user.GetAllUsers()
	if err != nil {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user_data)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("User handler reached")

	_, str_id := path.Split(r.URL.Path)
	id, err := strconv.Atoi(strings.Trim(str_id, "/"))
	if err != nil {
		http.Error(w, "Cant parse id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getUser(w, id)
	case http.MethodPut:
		modifyUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func modifyUser(w http.ResponseWriter, r *http.Request, id int) {
	log.Println("Modify user reached")

	// err := r.ParseMultipartForm(10 << 20) // 10 MB
	// if err != nil {
	// 	http.Error(w, "Error parsing form data", http.StatusBadRequest)
	// 	return
	// }

	var mod_user user.User
	mod_user.ID = id
	mod_user.Email = r.FormValue("email")
	mod_user.Password = r.FormValue("password")
	mod_user.UserName = r.FormValue("username")
	mod_user.Phone = r.FormValue("phone")
	mod_user.BDay = r.FormValue("bday")
	mod_user.Role = r.FormValue("role")
	mod_user.Gender = r.FormValue("gender")
	mod_user.Contry = r.FormValue("contry")

	log.Println("mod_user: ", mod_user)

	// Handle file upload
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()
		filePath := filepath.Join("uploads", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		mod_user.Photo = filePath
	}

	// Update mod_user in the database
	err = user.SaveUserByID(mod_user.ID, mod_user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func getUser(w http.ResponseWriter, id int) {
	log.Println("Get user reached")

	user_data, err := user.GetUserByID(id)
	if err != nil {
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user_data)
}

func deleteUser(w http.ResponseWriter, id int) {
	log.Println("Delete user reached")

	err := user.DeleteUserByID(id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

func Authority(w http.ResponseWriter, r *http.Request) {
	log.Println("Authority handler reached")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"authority": r.Header.Get("Authority")})
}
