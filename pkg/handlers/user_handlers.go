package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	hashing "fr_lab_1/pkg/hashing"
	user "fr_lab_1/pkg/user"
)

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
	mod_user.Contry = r.FormValue("country")
	mod_user.Agreement = r.FormValue("agreement") == "on"

	log.Println("Modifying user agreement to", r.FormValue("agreement"))

	// Handle file upload
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		uploadsDir := "../static/uploads/"

		filePath := uploadsDir + hashing.HashImageName(handler.Filename) + ".jpg"
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating destination file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
		mod_user.Photo = filePath
	} else {
		mod_user.Photo = r.FormValue("photoName")
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
