package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

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

func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home handler reached")
	// http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
	http.ServeFile(w, r, "../static/index.html")
}
