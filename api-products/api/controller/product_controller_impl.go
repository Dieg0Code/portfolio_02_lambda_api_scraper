package controller

import (
	"github.com/dieg0code/serverles-api-scraper/api/data/request"
	"github.com/dieg0code/serverles-api-scraper/api/service"
	"github.com/dieg0code/shared/json/response"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
}

// GetAll implements ProductController.
func (p *ProductControllerImpl) GetAll(ctx *gin.Context) {
	productResponse, err := p.ProductService.GetAll()
	if err != nil {
		logrus.WithError(err).Error("[ProductControllerImpl.GetAll] Error getting all products")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error getting all products",
			Data:    nil,
		}

		ctx.JSON(500, errorResponse)
		return
	}

	successResponse := response.BaseResponse{
		Code:    200,
		Status:  "OK",
		Message: "Success getting all products",
		Data:    productResponse,
	}

	ctx.JSON(200, successResponse)
}

// GetByID implements ProductController.
func (p *ProductControllerImpl) GetByID(ctx *gin.Context) {
	productId := ctx.Param("productId")
	if productId == "" {
		errorResponse := response.BaseResponse{
			Code:    400,
			Status:  "Bad Request",
			Message: "Product ID is required",
			Data:    nil,
		}

		ctx.JSON(400, errorResponse)
		return
	}

	productResponse, err := p.ProductService.GetByID(productId)
	if err != nil {
		logrus.WithError(err).Error("[ProductControllerImpl.GetByID] Error getting product by ID")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error getting product by ID",
			Data:    nil,
		}

		ctx.JSON(500, errorResponse)
		return
	}

	successResponse := response.BaseResponse{
		Code:    200,
		Status:  "OK",
		Message: "Success getting product by ID",
		Data:    productResponse,
	}

	ctx.JSON(200, successResponse)
}

// UpdateData implements ProductController.
func (p *ProductControllerImpl) UpdateData(ctx *gin.Context) {
	updateReq := request.UpdateDataRequest{}
	err := ctx.BindJSON(&updateReq)
	if err != nil {
		logrus.WithError(err).Error("[ProductControllerImpl.UpdateData] Error binding request")
		errorResponse := response.BaseResponse{
			Code:    400,
			Status:  "Bad Request",
			Message: "Error binding request",
			Data:    nil,
		}

		ctx.JSON(400, errorResponse)
		return
	}

	success, err := p.ProductService.UpdateData(updateReq)
	if err != nil {
		logrus.WithError(err).Error("[ProductControllerImpl.UpdateData] Error updating data")
		errorResponse := response.BaseResponse{
			Code:    500,
			Status:  "Internal Server Error",
			Message: "Error updating data",
			Data:    nil,
		}

		ctx.JSON(500, errorResponse)
		return
	}

	if !success {
		errorResponse := response.BaseResponse{
			Code:    400,
			Status:  "Bad Request",
			Message: "Error updating data",
			Data:    nil,
		}

		ctx.JSON(400, errorResponse)
		return
	}

	successResponse := response.BaseResponse{
		Code:    200,
		Status:  "OK",
		Message: "Data scraping started",
		Data:    nil,
	}

	ctx.JSON(200, successResponse)
}

func NewProductControllerImpl(productService service.ProductService) ProductController {
	return &ProductControllerImpl{ProductService: productService}
}
