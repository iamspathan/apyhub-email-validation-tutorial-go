package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var users = make(map[string]string)

func main() {
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve index.html file
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	w.Write(html)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == "" {
		http.Error(w, "Username, password and email are required", http.StatusBadRequest)
		return
	}

	// Validate email using API
	validEmail, err := validateEmail(email)
	if err != nil {
		http.Error(w, "Error validating email", http.StatusInternalServerError)
		return
	}

	if !validEmail {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	// Create account
	// ...

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully", username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	storedPassword, exists := users[username]
	if !exists || storedPassword != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome, %s!", username)
}

func validateEmail(email string) (bool, error) {
	url := "https://api.apyhub.com/validate/email/dns"

	payload := strings.NewReader(fmt.Sprintf("{\n    \"email\":\"%s\"\n}", email))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return false, err
	}

	req.Header.Add("apy-token", "APY0q37lVqWKoW6ggl7T9CmdsqFlZvigsiR70b0KIiGkrVuSR6aA8KsQU9O4WvbMjZAWHikXCZR")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  struct {
			Valid bool `json:"valid"`
		} `json:"result"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	if response.Status != "success" {
		return false, fmt.Errorf("API error: %s", response.Message)
	}

	return response.Result.Valid, nil
}
