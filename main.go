package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var users = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/signup", signupHandler)
	mux.HandleFunc("/login", loginHandler)

	corsHandler := corsMiddleware(mux)

	fmt.Println("Server running at http://localhost:3000")
	http.ListenAndServe(":3000", corsHandler)
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		// Here "*" is used to allow any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
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

type SignupForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var form SignupForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	password := form.Password
	email := form.Email

	if password == "" || email == "" {
		http.Error(w, "password and email are required", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s created successfully")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var form LoginForm
	err = json.Unmarshal(body, &form)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	username := form.Username
	password := form.Password

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

	payload := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	// Add your apy-token here.
	
	req.Header.Add("apy-token", "*********** ADD YOUR SECRET APY TOKEN HERE **************")
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
		Valid bool `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	return response.Valid, nil
}
