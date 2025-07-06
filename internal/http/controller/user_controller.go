package controller

import (
	"be/neurade/v2/internal/model"
	"be/neurade/v2/internal/service"
	"be/neurade/v2/internal/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type UserController struct {
	UserService *service.UserService
	JWTUtil     *util.JWTUtil
	log         *logrus.Logger
}

func NewUserController(userService *service.UserService, log *logrus.Logger, jwtSecret string) *UserController {
	jwtUtil := util.NewJWTUtil(jwtSecret, 24*time.Hour)
	return &UserController{
		UserService: userService,
		JWTUtil:     jwtUtil,
		log:         log,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.log.Error("Failed to parse form: ", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	fmt.Println(r.FormValue("email"))
	fmt.Println(r.FormValue("password"))
	fmt.Println(r.FormValue("role"))
	request := model.UserCreateRequest{
		Email:        r.FormValue("email"),
		PasswordHash: r.FormValue("password"),
		Role:         r.FormValue("role"),
	}

	if request.Email == "" || request.PasswordHash == "" {
		c.log.Error("Missing required fields")
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, err := c.UserService.Register(r.Context(), &request)
	if err != nil {
		c.log.Error("Failed to register user: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.log.Error("Failed to parse form: ", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	request := model.LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if request.Email == "" || request.Password == "" {
		c.log.Error("Missing required fields")
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	authResponse, err := c.UserService.Login(r.Context(), &request, c.JWTUtil)
	if err != nil {
		c.log.Error("Failed to login: ", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse)
}
