package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
	"net/http"
)

func (server *Server) getUsers(ctx *gin.Context) {
	req := &db.Users{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.GetUsers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (server *Server) createUsers(ctx *gin.Context) {
	req := &db.Users{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req.Id = 0 // ai columns
	if req.LocationId < 1 {
		ctx.JSON(http.StatusBadRequest, "missing location")
		return
	}

	user, err := server.store.CreateUsers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) updateUsers(ctx *gin.Context) {
	req := &db.Users{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.LocationId < 1 || req.Id < 1 {
		ctx.JSON(http.StatusBadRequest, "missing location or user")
		return
	}

	err := server.store.UpdateUsers(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

func (server *Server) createUsersLocations(ctx *gin.Context) {
	req := &db.UsersLocations{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.LocationId < 1 || req.UserId < 1 {
		ctx.JSON(http.StatusBadRequest, "missing location or user")
		return
	}

	err := server.store.CreateUsersLocations(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

type UserRequest struct {
	Number int `json:"number"`
}

func (server *Server) randomUsers(ctx *gin.Context) {
	req := &UserRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := 100
	res := make([]*db.Users, 0)

	for i := 0; i < req.Number; i++ {
		res = append(res, &db.Users{
			Name:             util.RandomString(6),
			PermanentAddress: util.RandomString(10),
			CurrentAddress:   util.RandomString(10),
			LocationId:       util.RandomInt(1, 20),
			HealthStatus:     util.RandomInt(0, 3),
		})

		if len(res) >= limit {
			err := server.store.CreateListUsers(ctx, res)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			res = make([]*db.Users, 0)
		}
	}

	if len(res) > 0 {
		err := server.store.CreateListUsers(ctx, res)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "success")
}
