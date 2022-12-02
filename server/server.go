package server

import (
	"crm-service-go/app/clients"
	"crm-service-go/app/controllers"
	controllerInside "crm-service-go/app/controllers/inside"
	"crm-service-go/app/middlewares"
	"crm-service-go/config"
	"crm-service-go/datasources"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	controller       *controllers.Controller
	controllerInside *controllerInside.InsideController
}

func NewServer(
	controller *controllers.Controller,
	controllerInside *controllerInside.InsideController,
) *Server {
	return &Server{
		controller,
		controllerInside,
	}
}

func (s *Server) Run() *gin.Engine {
	println("Start Server")

	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	resourceSaleOpportunity := clients.PolicyResourceSaleOpportunity
	resourceLead := clients.PolicyResourceLead
	actionReadAny := clients.AuthorizationActionReadAny
	actionCreateAny := clients.AuthorizationActionCreateAny
	actionUpdateAny := clients.AuthorizationActionUpdateAny
	actionDeleteAny := clients.AuthorizationActionDeleteAny

	// health check
	router.GET("/health", func(context *gin.Context) {
		db := datasources.MongoDatabase.Name()
		if db != config.GetConfig().DB.Name {
			context.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/live", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.GET("/ready", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	v1 := router.Group("/api/v1", middlewares.Authentication())
	{
		tag := v1.Group("/tags")
		{
			tag.GET("", s.controller.TagController.All)
		}
		lead := v1.Group("/leads")
		{
			lead.GET("", middlewares.Authorize(resourceLead, actionReadAny), s.controller.LeadController.List)
			lead.POST("", s.controller.LeadController.Create)
			lead.GET("/:id", s.controller.LeadController.Show)
			lead.PUT("/:id", s.controller.LeadController.Update)
			lead.DELETE("/:id", s.controller.LeadController.Delete)
		}
		note := v1.Group("/notes")
		{
			note.GET("", s.controller.NoteController.List)
			note.POST("", s.controller.NoteController.Create)
			note.GET("/:id", s.controller.NoteController.Show)
			note.PUT("/:id", s.controller.NoteController.Update)
			note.DELETE("/:id", s.controller.NoteController.Delete)
		}
		saleOpportunity := v1.Group("/sale-opportunities")
		{
			saleOpportunity.GET("", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.SaleOpportunityController.List)
			saleOpportunity.POST("", middlewares.Authorize(resourceSaleOpportunity, actionCreateAny), s.controller.SaleOpportunityController.Create)
			saleOpportunity.GET("/:id", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.SaleOpportunityController.Show)
			saleOpportunity.PUT("/:id", middlewares.Authorize(resourceSaleOpportunity, actionUpdateAny), s.controller.SaleOpportunityController.Update)
			saleOpportunity.DELETE("/:id", middlewares.Authorize(resourceSaleOpportunity, actionDeleteAny), s.controller.SaleOpportunityController.Delete)
			saleOpportunity.GET("/:id/logs", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.SaleOpportunityController.Logs)
		}

		statistics := v1.Group("/statistics")
		{
			statistics.GET("index", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.StatisticsController.Index)
			statistics.GET("sources", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.StatisticsController.Source)
			statistics.GET("stores", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.StatisticsController.Store)
			statistics.GET("employees", middlewares.Authorize(resourceSaleOpportunity, actionReadAny), s.controller.StatisticsController.Employee)
		}
	}

	rInside := router.Group("/api/v1/inside")
	{
		insideSaleOpportunity := rInside.Group("/sale-opportunities")
		{
			insideSaleOpportunity.GET("", s.controllerInside.InsideSaleOpportunityController.List)
			insideSaleOpportunity.POST("", s.controllerInside.InsideSaleOpportunityController.Create)
			insideSaleOpportunity.GET("/:id", s.controllerInside.InsideSaleOpportunityController.Show)
			insideSaleOpportunity.PUT("/:id", s.controllerInside.InsideSaleOpportunityController.Update)
		}

		insideLead := rInside.Group("/leads")
		{
			insideLead.GET("", s.controllerInside.InsideLeadController.List)
			insideLead.GET("/:id", s.controllerInside.InsideLeadController.Show)
		}
	}
	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	return router
}

func (s *Server) Close() error {
	println("Close Server")
	return nil
}
