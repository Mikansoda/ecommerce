package controller

import (
	"net/http"

	"ecommerce/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(s service.AuthService) *AuthController {
	return &AuthController{service: s}
}

// Struct request for endpoint register, verify, login, and refresh
type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=8,max=20,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type verifyReq struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Register user with email, username, and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      registerReq  true  "Register request"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Example 201 {json} Success Example:
// {
//   "message":"Registered successfully, check your email for OTP"
// }
// @Example 400 {json} Error Example:
// {
//   "message":"Failed to register, try again","detail":"error detail here"
// }
// @Router       /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	if err := a.service.Register(c, req.Username, req.Email, req.Password, "user"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to register, try again",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registered successfully, check your email for OTP",
	})
}

// VerifyOTP godoc
// @Summary      Verify user account with OTP
// @Description  Verify OTP sent to user's email
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      verifyReq  true  "Verify OTP request"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message":"Account successfully verified"
// }
// @Example 400 {json} Error Example:
// {
//   "message":"Failed to verify OTP","detail":"error detail here"
// }
// @Router       /auth/verify-otp [post]
func (a *AuthController) VerifyOTP(c *gin.Context) {
	var req verifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	if err := a.service.VerifyOTP(c, req.Email, req.OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to verify OTP",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account successfully verified",
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return access & refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginReq  true  "Login request"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "access_token":"<access>","refresh_token":"<refresh>"
// }
// @Example 401 {json} Error Example:
// {
//   "message":"Login failed","detail":"invalid credentials"
// }
// @Router       /auth/login [post]
func (a *AuthController) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	access, refresh, err := a.service.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Login failed",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Generate new access & refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshReq  true  "Refresh token request"
// @Success      200      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "access_token":"<new-access>","refresh_token":"<new-refresh>"
// }
// @Example 401 {json} Error Example:
// {
//   "message":"Failed to refresh token","detail":"refresh expired"
// }
// @Router       /auth/refresh [post]
func (a *AuthController) Refresh(c *gin.Context) {
	var req refreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data",
			"detail":  err.Error(),
		})
		return
	}

	access, refresh, err := a.service.Refresh(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Failed to refresh token",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate access token
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Example 200 {json} Success Example:
// {
//   "message":"Logged out successfully"
// }
// @Example 400 {json} Error Example:
// {
//   "message":"Failed to logout","detail":"invalid token"
// }
// @Router       /auth/logout [post]
func (a *AuthController) Logout(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	var token string
	if len(bearer) > 7 && bearer[:7] == "Bearer " {
		token = bearer[7:]
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing access token",
		})
		return
	}

	if err := a.service.Logout(c, token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to logout",
			"detail":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
