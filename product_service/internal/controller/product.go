package controller

import (
	"net/http"
	"product_service/internal/model/dto"
	"product_service/internal/oauth"
	"product_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Product struct {
	uc usecase.Product
}

func NewProduct(uc usecase.Product) *Product {
	return &Product{uc: uc}
}

// GetById godoc
// @Summary      Get product by id
// @Tags         product
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  model.Product
// @Failure      400
// @Router       /product/{id} [get]
func (p *Product) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid id")
		return
	}

	product, err := p.uc.GetById(c.Request.Context(), id)
	if err != nil {
		abortJson(c, 404, "not found")
	}

	c.IndentedJSON(200, product)
}

// Page godoc
// @Summary      Get product page
// @Tags         product
// @Produce      json
// @Param        page   query      int  true  "page number"
// @Param        limit  query      int  true  "limit number"
// @Success      200  {object}  	dto.ProductPage
// @Failure      400
// @Router       /product/ [get]
func (p *Product) Page(c *gin.Context) {
	page, err := atoiZeroIfEmpty(c.Query("page"))
	if err != nil {
		abortJson(c, http.StatusInternalServerError, "invalid page")
		return
	}

	limit, err := atoiZeroIfEmpty(c.Query("limit"))
	if err != nil {
		abortJson(c, http.StatusInternalServerError, "invalid limit")
		return
	}

	products, err := p.uc.GetPage(c.Request.Context(), page, limit)
	if err != nil {
		abortJson(c, 404, "not found")
		return
	}

	c.IndentedJSON(200, dto.ProductPage{Products: products, Page: page, Limit: limit})
}

func (p *Product) BatchByIds(c *gin.Context) {
	var ids dto.GetBatchRequest
	err := c.ShouldBindJSON(&ids)

	if err != nil {
		abortJson(c, 400, "invalid body json")
		return
	}

	batch, err := p.uc.GetBacth(c.Request.Context(), ids.ProductIds)
	if err != nil {
		c.IndentedJSON(200, &dto.GetBatchResponse{Products: make([]dto.GetBatchItem, 0)})
		return
	}
	c.IndentedJSON(200, &dto.GetBatchResponse{Products: batch})
}

// Update godoc
// @Summary      Update Product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param 		 id path string true "product id"
// @Param        updateProduct body dto.CreateProduct true "new product"
// @Success      200
// @Failure      400
// @Router       /product/{id} [put]
func (p *Product) Update(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)

	if !ok {
		abortJson(c, http.StatusBadRequest, "no valid token")
		return
	}

	if !user.HasRole("seller") {
		abortNotAllowed(c)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid id")
		return
	}

	var updateProduct dto.CreateProduct

	err = c.ShouldBindJSON(&updateProduct)

	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid body json")
		return
	}

	err = p.uc.Update(c.Request.Context(), id, updateProduct, user)
	if err != nil {
		abortJson(c, http.StatusBadRequest, "bad request")
		return
	}
}
