package controllers

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	"crypto/tls"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/its-ayush-07/go-neux-server/initializers"
	"github.com/its-ayush-07/go-neux-server/models"
	generate "github.com/its-ayush-07/go-neux-server/token"

	gomail "gopkg.in/mail.v2"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Email Or password is incorrect"
		valid = false
	}
	return valid, msg
}

// Goroutine to send welcome email to newly registered users
func sendEmail(email string, username string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("SENDER_EMAIL"))

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Welcome to GoNeux")

	// Set E-Mail body. You can set plain text or html with text/html

	message := fmt.Sprintf(`
Dear %s,

Welcome to GoNeux! Your account has been registered successfully.

Regards,
GoNeux team.
`, username)
	m.SetBody("text/plain", message)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("SENDER_EMAIL"), os.Getenv("SENDER_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// Function to sign up an account
func SignUp(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := HashPassword(user.Password)
	user.Password = hashedPassword

	user.Token, _ = generate.TokenGenerator(user.Email, user.UserName)

	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}

	go sendEmail(user.Email, user.UserName)

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", user.Token, 3600*24*30, "", "go-neux-client.vercel.app", true, true)
	c.JSON(200, models.User{ID: user.ID, Email: user.Email, UserName: user.UserName, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt})
}

// Function to login to an existing account
func Login(c *gin.Context) {
	var user models.User
	var founduser models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := initializers.DB.Where("email = ?", user.Email).First(&founduser).Error; err != nil {
		c.JSON(400, gin.H{"error": "User does not exists"})
		return
	}
	PasswordIsValid, msg := VerifyPassword(user.Password, founduser.Password)
	if !PasswordIsValid {
		c.JSON(500, gin.H{"error": msg})
		return
	}

	newToken, _ := generate.TokenGenerator(founduser.Email, founduser.UserName)
	founduser.Token = newToken
	if err := initializers.DB.Save(&founduser).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update token"})
		return
	}

	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", newToken, 3600*24*30, "", "go-neux-client.vercel.app", true, true)
	c.JSON(200, models.User{ID: founduser.ID, Email: founduser.Email, UserName: founduser.UserName, CreatedAt: founduser.CreatedAt, UpdatedAt: founduser.UpdatedAt})
}

// Function to logout
func Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("Authorization", "", -1, "", "go-neux-client.vercel.app", true, true)
	c.JSON(200, gin.H{"message": "Authorization cookie cleared and logged out"})
}

// Function to fetch user details by ID
func UserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "ID not specified"})
		return
	}
	var user models.User
	if err := initializers.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, models.User{ID: user.ID, Email: user.Email, UserName: user.UserName, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt})
}
