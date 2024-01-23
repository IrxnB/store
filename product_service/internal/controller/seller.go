package controller

import (
	"net/http"
	"product_service/internal/model/dto"
	"product_service/internal/oauth"
	"product_service/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Seller struct {
	uc usecase.Seller
}

func NewSeller(usecase usecase.Seller) *Seller {
	return &Seller{uc: usecase}
}

// Create godoc
// @Summary      Create seller
// @Tags         seller
// @Accept       json
// @Produce      json
// @Param        createSeller body dto.CreateSeller true "new seller"
// @Success      200
// @Failure      400
// @Router       /seller/ [post]
func (s *Seller) Create(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)

	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
		return
	}

	if !user.HasRole("seller") {
		abortNotAllowed(c)
		return
	}

	var createSeller dto.CreateSeller

	if err := c.ShouldBindJSON(&createSeller); err != nil {
		abortJson(c, http.StatusBadRequest, "invalid body json")
		return
	}

	if err := s.uc.Create(c.Request.Context(), createSeller, user); err != nil {
		abortJson(c, http.StatusBadRequest, "exists")
		return
	}

	jsonOK(c)
}

// GetAll godoc
// @Summary      Get seller list
// @Tags         seller
// @Produce      json
// @Success      200 {object} []model.Seller
// @Failure      400
// @Router       /seller/ [get]
func (s *Seller) GetAll(c *gin.Context) {
	sellers, err := s.uc.GetAll(c.Request.Context())
	if err != nil {
		abortJson(c, 404, "not found")
		return
	}

	response := dto.SellerList{Sellers: sellers}

	c.IndentedJSON(200, &response)
}

// Update godoc
// @Summary      Update seller
// @Tags         seller
// @Accept       json
// @Param 		 id path string true "seller id"
// @Param        updateSeller body dto.UpdateSeller true "new seller"
// @Success      200
// @Failure      400
// @Router       /seller/{id} [put]
func (s *Seller) Update(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)

	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
	}

	if !user.HasRole("seller") {
		abortNotAllowed(c)
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		abortJson(c, http.StatusBadRequest, "wrong id format")
	}

	if user.Id != id {
		abortJson(c, http.StatusForbidden, "not yours")
	}
	var updateSeller dto.UpdateSeller

	if err := c.ShouldBindJSON(&updateSeller); err != nil {
		abortJson(c, http.StatusBadRequest, "invalid body json")
		return
	}

	if err := s.uc.Update(c.Request.Context(), id, updateSeller, user); err != nil {
		abortJson(c, http.StatusBadRequest, "bad request")
		return
	}

	jsonOK(c)
}

// GetById godoc
// @Summary      Get seller  by id
// @Tags         seller
// @Param        id path string true "seller id"
// @Success      200 {object} model.Seller
// @Failure      400
// @Router       /seller/{id} [get]
func (s *Seller) GetById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		abortJson(c, http.StatusBadRequest, "wrong id format")
	}

	seller, err := s.uc.GetById(c.Request.Context(), id)

	if err != nil {
		abortJson(c, 404, "not found")
		return
	}

	c.IndentedJSON(200, &seller)
}

// GetProductPage
// @Summary      Get product page of a seller
// @Tags         seller
// @Produce      json
// @Param        id path string true "seller id"
// @Param        page   query      int  true  "page number"
// @Param        limit  query      int  true  "limit number"
// @Success      200 {object} 		dto.ProductPage
// @Failure      400
// @Router       /seller/{id}/products [get]
func (s *Seller) GetProductPage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		abortJson(c, http.StatusBadRequest, "wrong id format")
		return
	}

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

	products, err := s.uc.GetProductPage(c.Request.Context(), id, page, limit)
	if err != nil {
		abortJson(c, 404, "not found")
		return
	}

	c.IndentedJSON(200, dto.ProductPage{Products: products, Page: page, Limit: limit})
}

// Add product
// @Summary      Add product to seller
// @Tags         seller
// @Accept	     json
// @Param        id path string true "seller id"
// @Param        createProduct body dto.CreateProduct true "new product"
// @Success      200
// @Failure      400
// @Router       /seller/{id}/products [post]
func (s *Seller) AddProduct(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)
	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
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

	if id != user.Id {
		abortJson(c, http.StatusForbidden, "not yours")
	}

	var createProduct dto.CreateProductRequest

	err = c.ShouldBindJSON(&createProduct)
	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid body")
		return
	}

	err = s.uc.AddProduct(c.Request.Context(),
		dto.CreateProduct(createProduct), user.Id)

	if err != nil {
		abortJson(c, http.StatusInternalServerError, "internal")
		return
	}

}

func abortNotAllowed(c *gin.Context) {
	status := http.StatusForbidden
	c.AbortWithStatusJSON(status, gin.H{"status": status, "message": "not allowed"})
}

func abortJson(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"status": status, "message": message})
}

func jsonOK(c *gin.Context) {
	c.IndentedJSON(200, gin.H{"status": 200, "message": "ok"})
}

func atoiZeroIfEmpty(str string) (int, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		if str == "" {
			return 0, nil
		}
		return 0, err
	}
	return i, err
}
