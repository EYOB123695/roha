package http

import (
	"net/http"
	"strconv"

	"github.com/EYOB123695/roha/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(uc domain.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uc}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var body struct {
		Username  string `binding:"required"`
		Email     string `binding:"required"`
		Password  string `binding:"required"`
		AvatarURL string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read request body"})
		return
	}

	err := h.userUseCase.Signup(body.Username, body.Email, body.Password, body.AvatarURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `binding:"required"`
		Password string `binding:"required"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not read request body"})
		return
	}

	tokenString, err := h.userUseCase.Login(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) Validate(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You are logged in successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.userUseCase.GetUserProfile(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}
func (h *UserHandler) FollowUser(c *gin.Context) { 
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userVal.(domain.User)
	
	// Get target user ID from parameter
	idStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Call use case to follow
	err = h.userUseCase.FollowUser(currentUser.ID, uint(targetUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

func (h *UserHandler) UnfollowUser(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userVal.(domain.User)
	
	// Get target user ID from parameter
	idStr := c.Param("id")
	targetUserID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userUseCase.UnFollowUser(currentUser.ID, uint(targetUserID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}
