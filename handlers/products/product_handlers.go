package product_handler

import (
	"easystore/db"
	"easystore/dtos"
	"easystore/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var outlet models.Outlet
var product models.Product

// @Summary      Get details of a product of an outlet
// @Description  Get a product of an outlet and returns the product object
// @Param Authorization header string true "Bearer Token"
// @Param  outlet_id path string true "Outlet ID"
// @Param  id path string true "Product ID"
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Product  true  "Product Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/product/{id} [get]
func GetProductDetails(c *gin.Context) {
	if !setOutletFromContext(c) {
		return
	}
	product_id := c.Param("id")
	tx := db.DB.Where("outlet_id = ?", outlet.ID).First(&product, product_id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get product details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Product detail fetched successfully", "result": gin.H{"product": product}})
}

// @Summary      Create a product for an outlet
// @Description  Creates a new product for an outlet and returns the created product object
// @Param Authorization header string true "Bearer Token"
// @Param outlet_id path string true "Outlet ID"
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Product  true  "Product Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/product [post]
func Create(c *gin.Context) {

	if !setOutletFromContext(c) {
		return
	}

	var productDTO dtos.Product
	err := c.ShouldBindBodyWithJSON(&productDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get request body", "result": gin.H{"error": err.Error()}})
		return
	}

	var category models.ProductCategory
	tx := db.DB.First(&category, productDTO.CategoryId)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to get product category", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	product.OutletId = outlet.ID
	product.Title = productDTO.Title
	product.Description = productDTO.Description
	product.CategoryId = category.ID
	product.Status = productDTO.Status

	var productVarients []models.ProductVarient
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		tx = db.DB.Create(&product)
		if tx.Error != nil {
			return tx.Error
		}

		for _, varientDTO := range productDTO.Varients {
			var varient models.ProductVarient
			varient.Name = varientDTO.Name
			varient.ProductId = product.ID
			varient.Mrp = varientDTO.Mrp
			varient.SellingPrice = varientDTO.SellingPrice

			productVarients = append(productVarients, varient)
		}

		tx = db.DB.Create(&productVarients)

		if tx.Error != nil {
			return tx.Error
		}

		return nil

	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to create the product", "result": gin.H{"error": err.Error()}})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "Product created successfully", "result": gin.H{"product": product, "varients": productVarients}})

}

// @Summary      Update a product for an outlet
// @Description  Update a product of an outlet and returns the updated product object
// @Param Authorization header string true "Bearer Token"
// @Param  outlet_id path string true "Outlet ID"
// @Param  id path string true "Product ID"
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.Product  true  "Product Details"
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/product/{id} [put]
func Update(c *gin.Context) {
	if !setOutletFromContext(c) {
		return
	}

	product_id := c.Param("id")
	tx := db.DB.Where("outlet_id = ?", outlet.ID).First(&product, product_id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get product details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	var productDTO dtos.Product
	err := c.ShouldBindBodyWithJSON(&productDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get product details from body", "result": gin.H{"error": err.Error()}})
		return
	}

	var category models.ProductCategory
	tx = db.DB.Where("outlet_id = ?", outlet.ID).First(&category, productDTO.CategoryId)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get category details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	product.Title = productDTO.Title
	product.Description = productDTO.Description
	product.CategoryId = category.ID
	product.Status = productDTO.Status

	tx = db.DB.Where("outlet_id = ?", outlet.ID).Save(&product)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to update product", "result": gin.H{"error": tx.Error.Error()}})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "Product updated successfully", "result": gin.H{"product": product}})
}

func SampleHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// Private methods
var setOutletFromContext = func(c *gin.Context) bool {
	outlet_id := c.Param("outlet_id")

	tx := db.DB.First(&outlet, outlet_id)
	if tx.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to find outlet", "result": gin.H{"error": tx.Error.Error()}})
		return false
	}

	return true
}
