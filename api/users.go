package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/mikosco4real/simple_bank/db/sqlc"
	"github.com/mikosco4real/simple_bank/util"
)

// Create Users
type createUserRequest struct {
	Uuid      string `json:"uuid"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// Create Users
type userResponse struct {
	ID int `json:"id"`
	Uuid      string `json:"uuid"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse {
		ID: int(user.ID),
		Uuid: user.Uuid,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Email: user.Email,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	c := make(chan string, 1)

	go func(){ 
		hp, err := util.HashPassword(req.Password)
		if err != nil {
			c <- ""
		}
		c <- hp
		}()
	
	uuid := util.GenerateUuid()
	p := <- c
	
	if p == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	arg := db.CreateUserParams{
		Uuid: uuid,
		Firstname: req.Firstname,
		Lastname: req.Lastname,
		Email: req.Email,
		Password: p,
	}

	user, err := server.store.CreateUser(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusCreated, res)
}

// Get Users by ID
type getUserByIdRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUserById(ctx *gin.Context) {
	var req getUserByIdRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(context.Background(), req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	res := newUserResponse(user)

	ctx.JSON(http.StatusOK, res)
}


type listAllUsersRequest struct {
	PageID  int32 `json:"page_id" form:"page_id" binding:"required,min=1"`
	PageSize int32 `json:"page_size" form:"page_size" binding:"required,min=5,max=500"`
}
// Get all users in the system
func (server *Server) listAllUsers(ctx *gin.Context) {
	var req listAllUsersRequest
	
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// set a default for Limit
	// if req.PageSize == 0 {
	// 	req.PageSize = 20
	// }

	arg := db.ListUsersParams{
		Limit: req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	res := make([] userResponse, 0)

	for _, u := range users {
		res = append(res, newUserResponse(u))
	}

	ctx.JSON(http.StatusOK, res)
}

type loginUserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(context.Background(), req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ok := util.CheckPassword(user.Password, req.Password)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("password is incorrect")))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserResponse{
		AccessToken: accessToken,
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, res)
}