package http

import (
	"net/http"
	"strconv"

	"github.com/EYOB123695/roha/domain"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUseCase domain.CommentUseCase
}

// NewCommentHandler constructs a new HTTP adapter for Comment UseCase.
func NewCommentHandler(cc domain.CommentUseCase) *CommentHandler {
	return &CommentHandler{commentUseCase: cc}
}

func (h *CommentHandler) AddComment(c *gin.Context) {
	// 1. Retrieve the post ID from URL path parameters (/posts/:id/comments)
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 2. Extract and parse incoming JSON payload
	var body struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment body is required"})
		return
	}

	// 3. Extract the current user from context injected by requireAuth middleware
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	// 4. Fire the use case method
	comment, err := h.commentUseCase.AddComment(currentUser.ID, uint(postID), body.Body)
	if err != nil {
		// Differentiate status codes for semantic REST responses
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Return 201 Created state upon successful insertion
	c.JSON(http.StatusCreated, gin.H{"comment": comment})
}
