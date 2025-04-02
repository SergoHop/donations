package handlers

import (
	"log"
	"mydonate/internal/models"
	"mydonate/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{
    userService services.UserService
}

func NewUserHadler(userService services.UserService) *UserHandler{
    return &UserHandler{userService: userService}
}

// CreateUserHandler обрабатывает запрос на создание нового пользователя.
// POST /users
func(h *UserHandler) CreateUserHandler(c *gin.Context){
    var user models.User
    
    // Декодируем JSON из тела запроса в структуру User
    if err := c.ShouldBindJSON(&user); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "не то написал"})
		return
    }

    // Вызываем сервис для создания пользователя
    err := h.userService.Create(c.Request.Context(), &user)
    if err != nil {
		log.Printf("не тот пользователь: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "оибка пользователя"})
		return
	}
    c.Status(http.StatusCreated)
}


func (h *UserHandler) GetUserByIDHandler(c *gin.Context) {
    // Получаем ID пользователя из параметров маршрута
	idStr := c.Param("id")

    id ,err := strconv.ParseUint(idStr, 10, 32)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка id 400"})
		return
    }

    // Вызываем сервис для получения пользователя по ID
	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
    if err != nil{
        log.Printf("Error не получен id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка пользователя 500"})
		return
    }
    if user == nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "юзер нема 404"})
		return
    }
    c.JSON(http.StatusOK, user)
}


func (h *UserHandler) GetUserByEmailHandler(c *gin.Context){
    email := c.Query("email")
    if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "емаил сори"})
		return
	}
    user, err := h.userService.GetByEmail(c.Request.Context(), email)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "юзер нема"})
		return
	}
    if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "юзерн не найден"})
		return
	}

	c.JSON(http.StatusOK, user)
}