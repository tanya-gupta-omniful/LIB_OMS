package order

import (
	"OMS/service"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service service.Service
}

func NewController(s service.Service) *Controller {
	return &Controller{
		service: s,
	}
}

// HandleCSVFilePath handles the file path provided by the user

func (c *Controller) HandleOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		filePath := ctx.DefaultQuery("filePath", "")

		decodedPath, _ := url.QueryUnescape(filePath)

		if decodedPath == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "File path is required",
			})
			return
		}

		// Delegate CSV file processing to the service layer
		err := c.service.ProcessOrder(decodedPath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Respond with success
		ctx.JSON(http.StatusOK, gin.H{
			"message": "File processed successfully",
		})
	}
	// Get the file path from query parameters

}