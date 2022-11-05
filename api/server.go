package api

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/martikan/carrental_cars-api/config"
	dbConn "github.com/martikan/carrental_cars-api/db/sqlc"
	"github.com/martikan/carrental_common/middleware"
	"github.com/martikan/carrental_common/util"
)

type Api struct {
	router     *gin.Engine
	db         *dbConn.Queries
	sqlDb      *sql.DB
	config     *config.Config
	tokenMaker util.Maker
}

func InitApi() *Api {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v\n", err)
	}

	tokenMaker, err := util.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %v", err)
	}

	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", config.PostgreUser, config.PostgrePassword,
		config.PostgreHost, config.PostgrePort, config.PostgreDb, config.SSLMode)
	log.Println(dbUrl)

	conn, err := sql.Open(config.PostgreDriver, dbUrl)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("could not reach the database: %v", err)
	}

	log.Println("successfully connected to database")

	api := &Api{
		db:         dbConn.New(conn),
		sqlDb:      conn,
		config:     &config,
		tokenMaker: tokenMaker,
	}
	api.setupRouter()

	return api
}

func (a *Api) setupRouter() {
	router := gin.Default()

	// Open routes

	// FIXME: Implement it for k8s
	// Health
	// router.GET("/health/live", a.live)
	// router.GET("/health/ready", a.ready)

	// Authenticated routes

	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(a.tokenMaker))

	// Cars
	authRoutes.GET("/api/v1/cars", a.getAllCars)
	authRoutes.GET("/api/v1/cars/:id", a.getCarById)
	authRoutes.GET("/api/v1/cars/search", a.searchCars)
	authRoutes.POST("/api/v1/cars", a.createCar)
	authRoutes.PUT("/api/v1/cars/:id", a.updateCar)
	authRoutes.DELETE("/api/v1/cars/:id", a.deleteCar)

	// Brands
	authRoutes.POST("/api/v1/brands", a.createBrand)
	authRoutes.DELETE("/api/v1/brands/:id", a.deleteBrand)

	a.router = router
}

func (a *Api) Start() {
	log.Fatal(a.router.Run(":" + a.config.Port))
}
