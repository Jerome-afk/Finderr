package handlers

import (
	"time"

	"finderr/db"
	"finderr/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var CookieKey = "session_token"

func Register(c *fiber.Ctx) error {
	var form models.AuthForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("auth/signup", fiber.Map{
			"Error": "Invalid form data",
		})
	}

	// Check if user already exists
	if _, err := db.GetUserByEmail(form.Email); err == nil {
		return c.Status(fiber.StatusBadRequest).Render("auth/signup", fiber.Map{
			"Error": "Email already registered",
		})
	}

	if _, err := db.GetUserByUsername(form.Username); err == nil {
		return c.Status(fiber.StatusBadRequest).Render("auth/signup", fiber.Map{
			"Error": "Username already taken",
		})
	}

	user := &models.User{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := db.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/signup", fiber.Map{
			"Error": "Failed to create user",
		})
	}

	// Create session
	token, err := db.GenerateSessionToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/signup", fiber.Map{
			"Error": "Failed to generate session token",
		})
	}

	if err := db.SetUserSession(user.ID, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/signup", fiber.Map{
			"Error": "Failed to set user session",
		})
	}

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     CookieKey,
		Value:    token,
		HTTPOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * 7 * time.Hour), // 1 week
	})

	return c.Redirect("/")
}

func Login(c *fiber.Ctx) error {
	var form models.AuthForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("auth/login", fiber.Map{
			"Error": "Invalid form data",
		})
	}

	user, err := db.GetUserByEmail(form.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).Render("auth/login", fiber.Map{
			"Error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).Render("auth/login", fiber.Map{
			"Error": "Invalid credentials",
		})
	}

	// Create session
	token, err := db.GenerateSessionToken()
	if err != nil {		return c.Status(fiber.StatusInternalServerError).Render("auth/login", fiber.Map{
			"Error": "Failed to generate session token",
		})
	}

	if err := db.SetUserSession(user.ID, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/login", fiber.Map{
			"Error": "Failed to set user session",
		})
	}

	// Set session cookie
	c.Cookie(&fiber.Cookie{
		Name:     CookieKey,
		Value:    token,
		HTTPOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * 7 * time.Hour), // 1 week
	})

	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	token := c.Cookies(CookieKey)
	if token == "" {
		return c.Redirect("/auth/login")
	}

	user, err := db.GetUserBySessionToken(token)
	if err != nil {
		return c.Redirect("/auth/login")
	}

	if err := db.DeleteUserSession(user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("auth/login", fiber.Map{
			"Error": "Failed to logout",
		})
	}

	// Clear session cookie
	c.ClearCookie(CookieKey)

	return c.Redirect("/auth/login")
}

func AuthRequired(c *fiber.Ctx) error {
	token := c.Cookies(CookieKey)
	if token == "" {
		return c.Redirect("/auth/login")
	}

	user, err := db.GetUserBySessionToken(token)
	if err != nil {
		// Clear invalid session cookie
		c.ClearCookie(CookieKey)
		// Redirect to login
		return c.Redirect("/auth/login")
	}

	// Set user in context for use in handlers
	c.Locals("user", user)

	return c.Next()
}
