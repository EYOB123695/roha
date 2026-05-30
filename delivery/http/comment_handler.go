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

func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	// 1. Retrieve the post ID from URL path parameters (/posts/:id/comments)
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// 2. Fire the usecase method to fetch comments
	comments, err := h.commentUseCase.GetCommentsByPostID(uint(postID))
	if err != nil {
		// Differentiate status codes for semantic REST responses
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return the comments with a 200 OK status
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// 1. Retrieve the comment ID from URL path parameters (/comments/:id)
	idStr := c.Param("id")
	commentID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// 2. Extract current user from requireAuth middleware context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userInterface.(domain.User)

	// 3. Fire use case to delete the comment
	err = h.commentUseCase.DeleteComment(currentUser.ID, uint(commentID))
	if err != nil {
		if err.Error() == "comment not found" || err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "unauthorized to delete this comment" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Return 200 OK state upon successful deletion
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
