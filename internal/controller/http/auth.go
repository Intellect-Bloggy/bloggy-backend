package http

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/controller/http/response"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    auth_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/auth/dto"
    user_usecase "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/user/dto"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type authUsecase interface {
    SignUp(ctx context.Context, dto user_usecase.CreateUserDTO) (entity.Auth, error)
    SignIn(ctx context.Context, dto user_usecase.GetByCredentialsDTO) (entity.Auth, error)
    Refresh(ctx context.Context, dto auth_usecase.RefreshDTO) (entity.Auth, error)
    Logout(ctx context.Context, dto auth_usecase.LogoutDTO) error
}

type authHandler struct {
    authUsecase
}

func NewAuthHandler(authUsecase authUsecase) *authHandler {
    return &authHandler{authUsecase: authUsecase}
}

func (h *authHandler) Register(router *gin.Engine) {
    auth := router.Group("/auth")
    {
        auth.POST("/signup", h.signUp)
        auth.POST("/signin", h.signIn)
        auth.POST("/refresh", h.refresh)
        auth.DELETE("/logout", h.logout)
    }
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body dto.CreateUserDTO true "account info"
// @Success 201 {object} response.TokenResponse
// @Failure 400,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/signup [post]
func (h *authHandler) signUp(c *gin.Context) {
    dto := user_usecase.CreateUserDTO{}
    if err := c.BindJSON(&dto); err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    res, err := h.authUsecase.SignUp(c.Request.Context(), dto)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, res)
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept json
// @Produce json
// @Param input body dto.GetByCredentialsDTO true "account username and password"
// @Success 200 {object} response.TokenResponse
// @Failure 400,403 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/signin [post]
func (h *authHandler) signIn(c *gin.Context) {
    dto := user_usecase.GetByCredentialsDTO{}
    if err := c.BindJSON(&dto); err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    res, err := h.authUsecase.SignIn(c.Request.Context(), dto)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, res)
}

// @Summary Refresh
// @Tags auth
// @Description get new access and refresh token
// @ID get-new-access-and-refresh-token
// @Accept json
// @Produce json
// @Param input body dto.RefreshDTO true "refresh token"
// @Success 200 {object} response.TokenResponse
// @Failure 400,403,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) refresh(c *gin.Context) {
    dto := auth_usecase.RefreshDTO{}
    if err := c.BindJSON(&dto); err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    res, err := h.authUsecase.Refresh(c.Request.Context(), dto)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, res)
}

// @Summary Logout
// @Tags auth
// @Description logout user session
// @ID logout
// @Accept json
// @Produce json
// @Param input body dto.LogoutDTO true "refresh token"
// @Success 200 {object} response.TokenResponse
// @Failure 400,403,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /auth/logout [delete]
func (h *authHandler) logout(c *gin.Context) {
    dto := auth_usecase.LogoutDTO{}
    if err := c.BindJSON(&dto); err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    if err := h.authUsecase.Logout(c.Request.Context(), dto); err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.Status(200)
}
