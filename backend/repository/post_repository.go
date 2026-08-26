package repository

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/EYOB123695/roha/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postRepository struct {
	db          *gorm.DB
	cacheMutex  sync.RWMutex
	cachedPosts []domain.Post
	lastFetched time.Time
	cacheTTL    time.Duration
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{
		db:       db,
		cacheTTL: 3 * time.Second, // Cache valid for 3 seconds
	}
}

func (r *postRepository) Create(p *domain.Post) error {
	gormPost := Post{
		UserID:    p.UserID,
		MediaURL:  p.MediaURL,
		MediaType: p.MediaType,
		Caption:   p.Caption,
	}

	result := r.db.Create(&gormPost)
	if result.Error != nil {
		return result.Error
	}

	p.ID = gormPost.ID
	p.CreatedAt = gormPost.CreatedAt
	p.UpdatedAt = gormPost.UpdatedAt

	// Invalidate posts cache
	r.cacheMutex.Lock()
	r.lastFetched = time.Time{}
	r.cacheMutex.Unlock()

	return nil
}

func (r *postRepository) GetAll() ([]domain.Post, error) {
	// Fast path: serve from RAM if cache is still fresh (sub-microsecond)
	r.cacheMutex.RLock()
	if !r.lastFetched.IsZero() && time.Since(r.lastFetched) < r.cacheTTL {
		posts := r.cachedPosts
		r.cacheMutex.RUnlock()
		return posts, nil
	}
	r.cacheMutex.RUnlock()

	// Slow path: refresh cache from database
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	// Double-check after acquiring write lock (another goroutine may have refreshed)
	if !r.lastFetched.IsZero() && time.Since(r.lastFetched) < r.cacheTTL {
		return r.cachedPosts, nil
	}

	var gormPosts []Post
	result := r.db.Preload("User").Limit(50).Find(&gormPosts)
	if result.Error != nil {
		return nil, result.Error
	}

	posts := make([]domain.Post, 0, len(gormPosts))
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID:        gp.ID,
			UserID:    gp.UserID,
			MediaURL:  gp.MediaURL,
			MediaType: gp.MediaType,
			Caption:   gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID:        gp.User.ID,
				Username:  gp.User.Username,
				Email:     gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}

	// Update cache
	r.cachedPosts = posts
	r.lastFetched = time.Now()

	return posts, nil
}

func (r *postRepository) GetByID(id uint) (*domain.Post, error) {
	var gp Post
	result := r.db.Preload("User").First(&gp, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &domain.Post{
		ID:        gp.ID,
		UserID:    gp.UserID,
		MediaURL:  gp.MediaURL,
		MediaType: gp.MediaType,
		Caption:   gp.Caption,
		CreatedAt: gp.CreatedAt,
		UpdatedAt: gp.UpdatedAt,
		User: domain.User{
			ID:        gp.User.ID,
			Username:  gp.User.Username,
			Email:     gp.User.Email,
			AvatarURL: gp.User.AvatarURL,
		},
	}, nil
}

func (r *postRepository) Update(p *domain.Post) error {
	var gp Post
	result := r.db.First(&gp, p.ID)
	if result.Error != nil {
		return result.Error
	}

	result = r.db.Model(&gp).Updates(Post{
		Caption: p.Caption,
	})
	if result.Error != nil {
		return result.Error
	}

	p.UpdatedAt = gp.UpdatedAt

	// Invalidate posts cache
	r.cacheMutex.Lock()
	r.lastFetched = time.Time{}
	r.cacheMutex.Unlock()

	return nil
}

func (r *postRepository) Delete(id uint) error {
	result := r.db.Delete(&Post{}, id)
	if result.Error != nil {
		return result.Error
	}

	// Invalidate posts cache
	r.cacheMutex.Lock()
	r.lastFetched = time.Time{}
	r.cacheMutex.Unlock()

	return nil
}


func (r *postRepository) GetFeed(userID uint) ([]domain.Post, error) { 
	var followingIDs []uint
	err := r.db.Table("followers").Where("follower_id = ?", userID).Pluck("following_id", &followingIDs).Error
	if err != nil { 
		return nil, err
	}

	if len(followingIDs) == 0 { 
		return []domain.Post{}, nil
	}
    var gormPosts []Post
    err = r.db.Preload("User").Where("user_id IN ?", followingIDs).Order("created_at desc").
		Find(&gormPosts).Error
	if err != nil {
		return nil, err
	}


	var posts []domain.Post
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID: gp.ID,
			UserID: gp.UserID,
			MediaURL: gp.MediaURL,
			MediaType: gp.MediaType,
			Caption: gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID: gp.User.ID,
				Username: gp.User.Username,
				Email: gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}
	return posts, nil



}

func (r *postRepository) LikePost(userID uint, postID uint) error {

	like := Like{
		UserID: userID,
		PostID : postID,
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&like).Error
}

func (r*postRepository) UnlikePost(userID uint, postID uint) error { 

	result := r.db.Where("user_id = ? AND post_id = ?" , userID, postID).Delete(&Like{})

	return result.Error
	
}

func (r *postRepository) BookmarkPost(userID uint, postID uint) error {
	bookmark := Bookmark{
		UserID: userID,
		PostID: postID,
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&bookmark).Error
}

func (r *postRepository) UnbookmarkPost(userID uint, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&Bookmark{})
	return result.Error
}

func (r *postRepository) GetBookmarkedPosts(userID uint) ([]domain.Post, error) {
	var bookmarkPostIDs []uint
	err := r.db.Model(&Bookmark{}).Where("user_id = ?", userID).Pluck("post_id", &bookmarkPostIDs).Error
	if err != nil {
		return nil, err
	}

	if len(bookmarkPostIDs) == 0 {
		return []domain.Post{}, nil
	}

	var gormPosts []Post
	err = r.db.Preload("User").Where("id IN ?", bookmarkPostIDs).Order("created_at desc").Find(&gormPosts).Error
	if err != nil {
		return nil, err
	}

	var posts []domain.Post
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID:        gp.ID,
			UserID:    gp.UserID,
			MediaURL:  gp.MediaURL,
			MediaType: gp.MediaType,
			Caption:   gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID:        gp.User.ID,
				Username:  gp.User.Username,
				Email:     gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}
	return posts, nil
}



func (r* postRepository) TrackActivity(userID uint,postID uint,actionType string, watchDuration int) error { 

	log:= UserActivityLog { 
		UserID:        userID,
		PostID:        postID,
		ActionType:    actionType,
		WatchDuration: watchDuration,

	}
	return r.db.Create(&log).Error

}

func (r *postRepository) GetRecommendations(userID uint) ([]domain.Post, error) {
	// Subqueries to exclude liked and bookmarked posts
	likedSubQuery := r.db.Table("likes").Select("post_id").Where("user_id = ?", userID)
	bookmarkedSubQuery := r.db.Table("bookmarks").Select("post_id").Where("user_id = ?", userID)
	// 1. Get the user's top tags ordered by InterestScore descending
	var userInterests []UserInterest
	err := r.db.Where("user_id = ?", userID).Order("interest_score desc").Find(&userInterests).Error
	if err != nil {
		return nil, err
	}
	var gormPosts []Post
	// If interests exist, query posts matching those tags
	if len(userInterests) > 0 {
		var tagIDs []uint
		tagScoreMap := make(map[uint]float64)
		for _, ui := range userInterests {
			tagIDs = append(tagIDs, ui.TagID)
			tagScoreMap[ui.TagID] = ui.InterestScore
		}
		err = r.db.Preload("User").Preload("Tags").
			Joins("JOIN post_tags ON post_tags.post_id = posts.id").
			Where("post_tags.tag_id IN ?", tagIDs).
			Where("posts.id NOT IN (?)", likedSubQuery).
			Where("posts.id NOT IN (?)", bookmarkedSubQuery).
			Find(&gormPosts).Error
		if err != nil {
			return nil, err
		}
		// Deduplicate and score the posts
		if len(gormPosts) > 0 {
			uniquePostsMap := make(map[uint]Post)
			for _, p := range gormPosts {
				uniquePostsMap[p.ID] = p
			}
			type ScoredPost struct {
				Post  Post
				Score float64
			}
			var scoredPosts []ScoredPost
			for _, gp := range uniquePostsMap {
				maxScore := 0.0
				for _, tag := range gp.Tags {
					if score, exists := tagScoreMap[tag.ID]; exists {
						if score > maxScore {
							maxScore = score
						}
					}
				}
				scoredPosts = append(scoredPosts, ScoredPost{Post: gp, Score: maxScore})
			}
			// Sort by tag InterestScore descending, then by created_at descending
			sort.Slice(scoredPosts, func(i, j int) bool {
				if scoredPosts[i].Score == scoredPosts[j].Score {
					return scoredPosts[i].Post.CreatedAt.After(scoredPosts[j].Post.CreatedAt)
				}
				return scoredPosts[i].Score > scoredPosts[j].Score
			})
			// Map back to domain model
			var posts []domain.Post
			for _, sp := range scoredPosts {
				posts = append(posts, domain.Post{
					ID:        sp.Post.ID,
					UserID:    sp.Post.UserID,
					MediaURL:  sp.Post.MediaURL,
					MediaType: sp.Post.MediaType,
					Caption:   sp.Post.Caption,
					CreatedAt: sp.Post.CreatedAt,
					UpdatedAt: sp.Post.UpdatedAt,
					User: domain.User{
						ID:        sp.Post.User.ID,
						Username:  sp.Post.User.Username,
						Email:     sp.Post.User.Email,
						AvatarURL: sp.Post.User.AvatarURL,
					},
				})
			}
			return posts, nil
		}
	}
	// Fallback: If no interests or matching posts, return latest posts not liked/bookmarked
	err = r.db.Preload("User").
		Where("id NOT IN (?)", likedSubQuery).
		Where("id NOT IN (?)", bookmarkedSubQuery).
		Order("created_at desc").
		Limit(20).
		Find(&gormPosts).Error
	if err != nil {
		return nil, err
	}
	var posts []domain.Post
	for _, gp := range gormPosts {
		posts = append(posts, domain.Post{
			ID:        gp.ID,
			UserID:    gp.UserID,
			MediaURL:  gp.MediaURL,
			MediaType: gp.MediaType,
			Caption:   gp.Caption,
			CreatedAt: gp.CreatedAt,
			UpdatedAt: gp.UpdatedAt,
			User: domain.User{
				ID:        gp.User.ID,
				Username:  gp.User.Username,
				Email:     gp.User.Email,
				AvatarURL: gp.User.AvatarURL,
			},
		})
	}
	return posts, nil
}





