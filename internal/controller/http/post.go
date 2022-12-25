package http

import (
    "context"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "github.com/Intellect-Bloggy/bloggy-backend/internal/controller/http/response"
    "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/entity"
    post_dto "github.com/Intellect-Bloggy/bloggy-backend/internal/domain/usecase/post/dto"
    "github.com/Intellect-Bloggy/bloggy-backend/pkg/errors"
)

type postUsecase interface {
    Create(ctx context.Context, dto post_dto.CreatePostDTO, userId int) (int, error)
    GetById(ctx context.Context, id int) (entity.Post, error)
    DeleteById(ctx context.Context, id int, userId int) error
}

type postHandler struct {
    authMiddleware
    postUsecase
}

func NewPostHandler(authMiddleware authMiddleware, postUsecase postUsecase) *postHandler {
    return &postHandler{postUsecase: postUsecase, authMiddleware: authMiddleware}
}

func (h *postHandler) Register(router *gin.Engine) {
    post := router.Group("/post")
    {
        post.GET("/:id", h.getById)
        protected := post.Group("", h.authMiddleware.UserIdentity)
        {
            protected.POST("", h.create)
            protected.DELETE("/:id", h.deleteById)
        }
    }
}

// @Summary CreatePost
// @Tags post
// @Description create post
// @Security ApiKeyAuth
// @ID create-post
// @Accept json
// @Produce json
// @Param input body dto.CreatePostDTO true "post information"
// @Success 200 {object} entity.CreatePostResponse
// @Failure 400,401,403,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /post [post]
func (h *postHandler) create(c *gin.Context) {
    userIdI, exist := c.Get("user_id")
    if !exist {
        response.ResponseWithError(c, errors.ErrIdNotFound)
        return
    }

    dto := post_dto.CreatePostDTO{}

    if err := c.BindJSON(&dto); err != nil {
        response.ResponseWithError(c, errors.NewHTTPError(http.StatusBadRequest, err))
        return
    }

    userId, _ := userIdI.(int)

    id, err := h.postUsecase.Create(c.Request.Context(), dto, userId)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusCreated, entity.CreatePostResponse{Id: id})
}

// @Summary GetPostById
// @Tags post
// @Description get post
// @ID get-post-by-id
// @Accept json
// @Produce json
// @Param id path integer true "post id"
// @Success 200 {object} entity.Post
// @Failure 400,401,403,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /post/{id} [get]
func (h *postHandler) getById(c *gin.Context) {
    postIdS := c.Param("id")
    postId, err := strconv.Atoi(postIdS)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    post, err := h.postUsecase.GetById(c.Request.Context(), postId)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.JSON(http.StatusOK, post)
}

// @Summary DeletePostById
// @Tags post
// @Description delete post
// @Security ApiKeyAuth
// @ID delete-post-by-id
// @Accept json
// @Produce json
// @Param id path integer true "post id"
// @Success 200
// @Failure 400,401,403,404,409 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router /post/{id} [delete]
func (h *postHandler) deleteById(c *gin.Context) {
    postIdS := c.Param("id")
    userIdI, exist := c.Get("user_id")
    if !exist {
        response.ResponseWithError(c, errors.ErrIdNotFound)
        return
    }

    userId, _ := userIdI.(int)
    postId, err := strconv.Atoi(postIdS)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    err = h.postUsecase.DeleteById(c.Request.Context(), postId, userId)
    if err != nil {
        response.ResponseWithError(c, errors.EtoHe(err))
        return
    }

    c.Status(http.StatusOK)
}
