package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/martikan/carrental_cars-api/db/sqlc"
	e "github.com/martikan/carrental_common/exception"
)

type pageableRequest struct {
	Page int32 `form:"page" binding:"min=1"`
	Size int32 `form:"size" binding:"min=5,max=50"`
}

type getCarByIdRequest struct {
	ID int64 `form:"id" binding:"required"`
}

type createCarRequest struct {
	BrandID int64  `json:"brand_id" binding:"required"`
	Color   string `json:"color" binding:"required"`
	Serial  string `json:"serial" binding:"required"`
	Comfort string `json:"comfort" binding:"required,oneof=S A B C D"`
}

type updateCarRequest struct {
	BrandID   int64  `json:"brand_id" binding:"required"`
	Color     string `json:"color" binding:"required"`
	Serial    string `json:"serial" binding:"required"`
	Comfort   string `json:"comfort" binding:"required,oneof=S A B C D"`
	Available bool   `json:"available" binding:"required"`
}

type carResponse struct {
	ID      int64  `json:"id" binding:"required"`
	BrandID int64  `json:"brand_id" binding:"required"`
	Color   string `json:"color" binding:"required"`
	Serial  string `json:"serial" binding:"required"`
	Comfort string `json:"comfort" binding:"required"`
}

func newCarResponse(car db.Car) carResponse {
	return carResponse{
		ID:      car.ID,
		BrandID: car.BrandID,
		Color:   car.Color,
		Serial:  car.Serial,
		Comfort: car.Comfort,
	}
}

func (a *Api) getAllCars(ctx *gin.Context) {
	var req pageableRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	args := db.GetAllCarsParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	cars, err := a.db.GetAllCars(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	} else if len(cars) == 0 {
		ctx.JSON(http.StatusOK, make([]string, 0))
		return
	}

	var carsDTO []carResponse
	for _, c := range cars {
		dto := carResponse{
			ID:      c.ID,
			BrandID: c.BrandID,
			Color:   c.Color,
			Serial:  c.Serial,
			Comfort: c.Comfort,
		}

		carsDTO = append(carsDTO, dto)
	}

	ctx.JSON(http.StatusOK, carsDTO)

}

func (a *Api) searchCars(ctx *gin.Context) {
	var req pageableRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	q := ctx.Query("q")
	var qParam sql.NullString
	if q != "" {
		qParam.String = q
		qParam.Valid = true
	}

	available := ctx.Query("available")
	var availableParam sql.NullBool
	switch available {
	case "false":
		availableParam.Bool = false
		availableParam.Valid = true
	case "true":
		availableParam.Bool = true
		availableParam.Valid = true
	default:
		availableParam.Valid = false
	}

	args := db.SearchCarsParams{
		AvailableParam: availableParam,
		QParam:         qParam,
		LimitParam:     req.Size,
		OffsetParam:    (req.Page - 1) * req.Size,
	}

	cars, err := a.db.SearchCars(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	} else if len(cars) == 0 {
		ctx.JSON(http.StatusOK, make([]string, 0))
		return
	}

	var carsDTO []carResponse
	for _, c := range cars {
		dto := carResponse{
			ID:      c.ID,
			BrandID: c.BrandID,
			Color:   c.Color,
			Serial:  c.Serial,
			Comfort: c.Comfort,
		}

		carsDTO = append(carsDTO, dto)
	}

	ctx.JSON(http.StatusOK, carsDTO)
}

func (a *Api) getCarById(ctx *gin.Context) {
	var req getCarByIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	car, err := a.db.GetCarById(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, e.ApiMessage("car is not found", 404))
			return
		}

		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	dto := carResponse{
		ID:      car.ID,
		BrandID: car.BrandID,
		Color:   car.Color,
		Serial:  car.Serial,
		Comfort: car.Comfort,
	}

	ctx.JSON(http.StatusOK, dto)
}

func (a *Api) createCar(ctx *gin.Context) {
	var req createCarRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	args := db.CreateCarParams{
		BrandID: req.BrandID,
		Color:   req.Color,
		Serial:  req.Serial,
		Comfort: req.Comfort,
	}

	car, err := a.db.CreateCar(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	dto := newCarResponse(car)

	ctx.JSON(http.StatusOK, dto)
}

func (a *Api) updateCar(ctx *gin.Context) {
	var idReq getCarByIdRequest
	if err := ctx.BindQuery(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	var updateReq updateCarRequest
	if err := ctx.BindQuery(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	args := db.UpdateCarParams{
		ID:        idReq.ID,
		BrandID:   updateReq.BrandID,
		Color:     updateReq.Color,
		Serial:    updateReq.Color,
		Comfort:   updateReq.Comfort,
		Available: updateReq.Available,
	}

	car, err := a.db.UpdateCar(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	dto := newCarResponse(car)

	ctx.JSON(http.StatusOK, dto)
}

func (a *Api) deleteCar(ctx *gin.Context) {
	var req getCarByIdRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, e.ApiError(err))
		return
	}

	err := a.db.DeleteCar(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, e.ApiError(err))
		return
	}

	ctx.JSON(http.StatusOK, e.ApiMessage("car has been deleted", 200))
}
