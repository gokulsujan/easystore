package product_varient_handler

import (
	"easystore/db"
	"easystore/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var productVarient models.ProductVarient

// @Summary      Create a product varient for an outlet
// @Description  Creates a new product varient for an outlet and returns the created product varient object
// @Param Authorization header string true "Bearer Token"
// @Param product_id path string true "Product ID"
// @Tags         Product Varient
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.ProductVarient  true  "Product Details"
// @Success      202  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /product/{product_id}/product-varient [post]
func Create(c *gin.Context) {
	err := c.ShouldBindBodyWithJSON(&productVarient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the request body", "result": gin.H{"error": err.Error()}})
		return
	}

	productIdStr := c.Param("product_id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil && productIdStr == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Invalid product id", "result": gin.H{"error": err.Error()}})
		return
	}

	product := &productVarient.Product
	tx := db.DB.First(product, productId)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the product details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	productVarient.ProductId = product.ID

	tx = db.DB.Create(&productVarient)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to create product varient", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "Product varient created successfully", "result": gin.H{"varient": productVarient}})
}

// @Summary      Update a product varient for an outlet
// @Description  Update a product varient for an outlet and returns the updated product varient object
// @Param Authorization header string true "Bearer Token"
// @Param product_id path string true "Product ID"
// @Param varient_id path string true "Product Variend ID"
// @Tags         Product Varient
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.ProductVarient  true  "Product Details"
// @Success      202  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /product/{product_id}/product-varient{varient_id} [put]
func Update(c *gin.Context) {
	varient_id := c.Param("varient_id")
	product_id := c.Param("product_id")
	tx := db.DB.Where("product_id = ?", product_id).First(&productVarient, varient_id)

	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid product varient id"})
		return
	}

	var updatedProductVarient models.ProductVarient
	err := c.ShouldBindBodyWithJSON(&updatedProductVarient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to get the product varient", "result": err.Error()})
		return
	}

	updatedProductVarient.ID = productVarient.ID
	tx = db.DB.Where("product_id = ?", product_id).Updates(&updatedProductVarient)

	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable update product varient", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "Product varient updated successfully", "result": gin.H{"varient": updatedProductVarient}})
}

// @Summary      Get all product varient for product
// @Description  Get all product varients for a product and returns the product varient list object
// @Param Authorization header string true "Bearer Token"
// @Param product_id path string true "Product ID"
// @Tags         Product Varient
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.ProductVarient  true  "Product Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /product/{product_id}/product-varient [get]
func GetProductVarients(c *gin.Context) {
	productIdStr := c.Param("product_id")
	var varients []models.ProductVarient

	tx := db.DB.Where("product_id = ?", productIdStr).Find(&varients)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the product varients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product varients successfully fetched", "result": gin.H{"varients": varients}})

}

// @Summary      Get a product varient for product
// @Description  Get a product varient for a product and returns the product varient object
// @Param Authorization header string true "Bearer Token"
// @Param product_id path string true "Product ID"
// @Param varient_id path string true "Product Varient ID"
// @Tags         Product Varient
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.ProductVarient  true  "Product Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /product/{product_id}/product-varient/{varient_id} [get]
func GetProductVarient(c *gin.Context) {
	productIdStr := c.Param("product_id")
	vaientIdStr := c.Param("varient_id")

	tx := db.DB.Where("product_id = ?", productIdStr).First(&productVarient, vaientIdStr)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the product varient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product varient successfully fetched", "result": gin.H{"varient": productVarient}})
}
