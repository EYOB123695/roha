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

// Signup godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param body body map[string]string true "User data"
// @Success 200 {object} map[string]string "{"message": "User created successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Could not read request body"}"
// @Router /signup [post]
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

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param body body map[string]string true "Login credentials"
// @Success 200 {object} map[string]string "{"token": "JWT_TOKEN"}"
// @Failure 400 {object} map[string]string "{"error": "Could not read request body"}"
// @Router /login [post]
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

// Validate godoc
// @Summary Validate JWT Token
// @Description Validate if the user's JWT token is valid and active
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /validate [get]
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

// GetUserProfile godoc
// @Summary Get user profile
// @Description Fetch user profile by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid user ID"}"
// @Failure 404 {object} map[string]string "{"error": "User not found"}"
// @Router /users/{id} [get]
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
// FollowUser godoc
// @Summary Follow a user
// @Description Follow another user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Successfully followed user"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid user ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /users/{id}/follow [post]
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

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Unfollow another user by their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Successfully unfollowed user"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid user ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /users/{id}/unfollow [post]
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


// GetFollowers godoc
// @Summary Get user followers
// @Description Fetch followers of a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid user ID"}"
// @Failure 404 {object} map[string]string "{"error": "User not found"}"
// @Router /users/{id}/followers [get]
func (h *UserHandler) GetFollowers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followers, err := h.userUseCase.GetFollowers(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers})
}

// GetFollowing godoc
// @Summary Get user following
// @Description Fetch users followed by a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid user ID"}"
// @Failure 404 {object} map[string]string "{"error": "User not found"}"
// @Router /users/{id}/following [get]
func (h *UserHandler) GetFollowing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	following, err := h.userUseCase.GetFollowing(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following})
}
