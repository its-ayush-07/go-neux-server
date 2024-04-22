package controllers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/its-ayush-07/go-neux-server/initializers"
	"github.com/its-ayush-07/go-neux-server/models"
)

// Function to fetch most recent 5 articles
func RecentArticles(c *gin.Context) {
	var articles []models.Article
	err := initializers.DB.Order("created_at desc").Limit(5).Find(&articles).Error
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to retrieve articles",
		})
		return
	}

	c.JSON(200, gin.H{
		"articles": articles,
	})
}

// Function to create an article
func CreateArticle(c *gin.Context) {
	var article models.Article

	if err := c.BindJSON(&article); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	cookie_username, _ := c.Get("username")
	if cookie_username != article.Author {
		c.JSON(401, gin.H{"error": "User unauthenticated"})
		return
	}

	if err := initializers.DB.Create(&article).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, "Article created successfully")
}

// Function to fetch an article by ID
func ArticleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID not specified"})
		return
	}
	var article models.Article
	if err := initializers.DB.First(&article, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(200, article)
}

// Function to fetch articles by search query
func SearchArticles(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		c.JSON(400, gin.H{"error": "Missing query parameter"})
		return
	}

	var articles []models.Article

	if err := initializers.DB.Where("LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(author) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").Find(&articles).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to search articles"})
		return
	}

	c.JSON(200, articles)
}

// Function to update like of an article by ID
func UpdateLikes(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID not specified"})
		return
	}

	var body struct {
		IsLike bool
	}
	c.Bind(&body)
	var article models.Article
	result := initializers.DB.First(&article, id)
	if result.Error != nil {
		c.JSON(404, "Article not found")
		return
	}
	// increment like
	if body.IsLike {
		article.Likes++
	} else { // decrement like
		article.Likes--
	}
	result = initializers.DB.Save(&article)

	if result.Error != nil {
		c.JSON(500, "Could not update database")
		return
	}

	c.JSON(200, gin.H{
		"article": article,
	})
}

// Function to delete an article by ID
func DeleteArticleByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID not specified"})
		return
	}

	var article models.Article
	if err := initializers.DB.Where("id = ?", id).First(&article).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to find specified article"})
		return
	}

	cookie_username, _ := c.Get("username")
	if cookie_username != article.Author {
		c.JSON(401, gin.H{"error": "User unauthenticated"})
		return
	}

	// Delete the article
	if err := initializers.DB.Delete(&article).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to delete specified article"})
		return
	}

	c.JSON(200, "Article deleted successfully")
}

// Function to fetch articles by author
func ArticleByAuthor(c *gin.Context) {
	authorID := c.Param("id")
	if authorID == "" {
		c.JSON(400, gin.H{"error": "ID not specified"})
		return
	}

	var articles []models.Article
	if err := initializers.DB.Where("author_id = ?", authorID).Find(&articles).Error; err != nil {
		c.JSON(500, gin.H{"error": "Unable to fetch articles"})
		return
	}

	c.JSON(200, articles)
}
