package http

import (
	"net/http"
	"strconv"

	"github.com/EYOB123695/roha/domain"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUseCase domain.PostUseCase
}

func NewPostHandler(pc domain.PostUseCase) *PostHandler {
	return &PostHandler{postUseCase: pc}
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a post with media URL and caption
// @Tags posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body map[string]string true "Post data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid request body"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var body struct {
		MediaURL  string `binding:"required"`
		MediaType string `binding:"required"`
		Caption   string
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Retrieve user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	post, err := h.postUseCase.CreatePost(currentUser.ID, body.MediaURL, body.MediaType, body.Caption)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// GetPosts godoc
// @Summary Get all posts
// @Description Fetch all public posts
// @Tags posts
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.postUseCase.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// GetPost godoc
// @Summary Get a post
// @Description Fetch a post by ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postUseCase.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update post caption by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param body body map[string]string true "Update data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 403 {object} map[string]string "{"error": "Forbidden"}"
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var body struct {
		Caption string `binding:"required"`
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	post, err := h.postUseCase.UpdatePost(currentUser.ID, uint(id), body.Caption)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post by ID
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Post deleted successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 403 {object} map[string]string "{"error": "Forbidden"}"
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	err = h.postUseCase.DeletePost(currentUser.ID, uint(id))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}


// GetFeed godoc
// @Summary Get user feed
// @Description Fetch personalized feed for authenticated user
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /feed [get]
func (h *PostHandler) GetFeed(c *gin.Context) {
	// 1. Extract current user from middleware context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	// 2. Fetch personalized feed using the usecase
	posts, err := h.postUseCase.GetFeed(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Respond with the list of feed posts
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}


// LikePost godoc
// @Summary Like a post
// @Description Adds a like relation between the authenticated user and a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Post liked successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) { 
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return 
	}

	userInterface, exists := c.Get("user")
	if !exists { 
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	currentUser := userInterface.(domain.User)

	err = h.postUseCase.LikePost(currentUser.ID, uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 4. Return 200 OK on success
	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}


// UnlikePost godoc
// @Summary Unlike a post
// @Description Removes a like relation between the authenticated user and a specific post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Post unliked successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /posts/{id}/like [delete]
func (h *PostHandler) UnlikePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)
	// 3. Fire the usecase method
	err = h.postUseCase.UnlikePost(currentUser.ID, uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 4. Return 200 OK on success
	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// BookmarkPost godoc
// @Summary Bookmark a post
// @Description Saves a post to the authenticated user's bookmarks
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Post bookmarked successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Router /posts/{id}/bookmark [post]
func (h *PostHandler) BookmarkPost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	err = h.postUseCase.BookmarkPost(currentUser.ID, uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post bookmarked successfully"})
}

// UnbookmarkPost godoc
// @Summary Remove bookmark from a post
// @Description Removes a post from the authenticated user's bookmarks
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Bookmark removed successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid post ID"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Router /posts/{id}/bookmark [delete]
func (h *PostHandler) UnbookmarkPost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	err = h.postUseCase.UnbookmarkPost(currentUser.ID, uint(postID))
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bookmark removed successfully"})
}

// GetBookmarkedPosts godoc
// @Summary Get bookmarked posts
// @Description Fetch a list of posts bookmarked by the current logged-in user
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /bookmarks [get]
func (h *PostHandler) GetBookmarkedPosts(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	posts, err := h.postUseCase.GetBookmarkedPosts(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}




// TrackActivity godoc
// @Summary Track engagement activity on a post
// @Description Captures implicit engagement logs (views, clicks, shares)
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param body body map[string]interface{} true "Activity data (action_type: string, watch_duration: int)"
// @Security BearerAuth
// @Success 200 {object} map[string]string "{"message": "Activity tracked successfully"}"
// @Failure 400 {object} map[string]string "{"error": "Invalid input"}"
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Failure 404 {object} map[string]string "{"error": "post not found"}"
// @Router /posts/{id}/activity [post]
func (h *PostHandler) TrackActivity(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var body struct {
		ActionType    string `json:"action_type" binding:"required"`
		WatchDuration int    `json:"watch_duration"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	err = h.postUseCase.TrackActivity(currentUser.ID, uint(postID), body.ActionType, body.WatchDuration)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity tracked successfully"})
}

// GetRecommendations godoc
// @Summary Get post recommendations
// @Description Fetch recommendations based on user interests, excluding liked/bookmarked posts
// @Tags posts
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string "{"error": "Unauthorized"}"
// @Router /posts/recommendations [get]
func (h *PostHandler) GetRecommendations(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	posts, err := h.postUseCase.GetRecommendations(currentUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}
