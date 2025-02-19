package product_category_handler

import (
	"easystore/db"
	"easystore/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var productCategory models.ProductCategory
var outlet models.Outlet

// @Summary      Create a product category for an outlet
// @Description  Creates a new product category for an outlet and returns the created product category object
// @Param Authorization header string true "Bearer Token"
// @Param outlet_id path string true "Outlet ID"
// @Tags         Product Category
// @Accept       json
// @Produce      json
// @Param        outlet  body  dtos.ProductCategory  true  "Product Details"
// @Success      202  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/product-category [post]
func Create(c *gin.Context) {
	outlet_id := c.Param("outlet_id")
	err := c.ShouldBindBodyWithJSON(&productCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the request body", "result": gin.H{"error": err.Error()}})
		return
	}

	tx := db.DB.First(&outlet, outlet_id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to find the outlet details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	productCategory.OutletId = outlet.ID
	if productCategory.Title == "" && productCategory.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Title and description should not be empty"})
		return
	}

	tx = db.DB.Create(&productCategory)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to craete category", "result": gin.H{"error": tx.Error.Error()}})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "success", "message": "Category created successfully", "result": gin.H{"category": productCategory}})
}

// @Summary      Create a product category for an outlet
// @Description  Creates a new product category for an outlet and returns the created product category object
// @Param Authorization header string true "Bearer Token"
// @Param outlet_id path string true "Outlet ID"
// @Param category_id path string true "Product Category ID"
// @Tags         Product Category
// @Produce      json
// @Success      200  {object}  dtos.SuccessResponse
// @Failure      400  {object}  dtos.ErrorResponse
// @Failure      500  {object}  dtos.ErrorResponse
// @Security BearerAuth
// @Router       /outlet/{outlet_id}/product-category{:category_id} [get]
func GetProductCategoryDetail(c *gin.Context) {
	// Step 1 -> Get category id from url params
	category_id := c.Param("category_id")

	// Step 2 -> Search the category using the id on db
	tx := db.DB.First(&productCategory, category_id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "message": "Unable to get the category details", "result": gin.H{"error": tx.Error.Error()}})
		return
	}
	// Step 3 -> Return the category objects
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category details fetched successfully", "result": gin.H{"category": productCategory}})
}
