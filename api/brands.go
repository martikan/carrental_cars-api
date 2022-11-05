package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	e "github.com/martikan/carrental_common/exception"
)

type createBrandRequest struct {
	Name string `json:"name" binding:"required"`
}

type getBrandByIdRequest struct {
	ID int64 `form:"id" binding:"required"`
}

func (a *Api) createBrand(ctx *gin.Context) {
	var req createBrandRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	brand, err := a.db.CreateBrand(ctx, req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	ctx.JSON(http.StatusCreated, brand)
}

func (a *Api) deleteBrand(ctx *gin.Context) {
	var req getBrandByIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	err := a.db.DeleteBrand(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	ctx.JSON(http.StatusOK, e.ApiMessage("brand has been deleted", 200))
}
