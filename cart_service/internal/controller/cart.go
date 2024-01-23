package controller

import (
	"net/http"
	"store/cart_service/internal/model/dto"
	"store/cart_service/internal/oauth"
	"store/cart_service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Cart struct {
	uc usecase.Cart
}

func NewCart(uc usecase.Cart) *Cart {
	return &Cart{uc: uc}
}

// GetCart godoc
// @Summary      Get content of a cart
// @Tags         cart
// @Produce      json
// @Success      200 {object} []model.CartEntry
// @Failure      400
// @Router       /cart/ [get]
func (cc *Cart) GetCart(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)
	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
		return
	}

	if !user.HasRole("user") {
		abortNotAllowed(c)
		return
	}

	entries, err := cc.uc.GetCart(c.Request.Context(), user.Id)
	if err != nil {
		abortJson(c, 404, "not found entries")
		return
	}

	c.IndentedJSON(200, &entries)
}

// AddOrUpdate godoc
// @Summary      Add to cart or update existing
// @Tags         cart
// @Accept      json
// @Param		 entries body dto.AddToCart true "cart entries"
// @Success      200
// @Failure      400
// @Router       /cart/ [post]
func (cc *Cart) AddOrUpdate(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)
	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
		return
	}

	if !user.HasRole("user") {
		abortNotAllowed(c)
		return
	}

	var addToCart dto.AddToCart

	err := c.ShouldBindJSON(&addToCart)

	if err != nil {
		abortJson(c, 400, "invalid body json")
		return
	}

	err = cc.uc.AddOrUpdate(c.Request.Context(), addToCart.Requests, user)
	if err != nil {
		abortJson(c, 400, "cannot add")
		return
	}

	jsonOK(c)
}

// RemoveFromCart godoc
// @Summary      Remove product from cart
// @Tags         cart
// @Accept       json
// @Param 		 productIds body dto.RemoveFromCart true "product ids to remove"
// @Success      200 {object} []model.CartEntry
// @Failure      400
// @Router       /cart/ [delete]
func (cc *Cart) RemoveFromCart(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)
	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
		return
	}

	if !user.HasRole("user") {
		abortNotAllowed(c)
		return
	}
	var removeFromCart dto.RemoveFromCart

	err := c.ShouldBindJSON(&removeFromCart)

	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid body json")
		return
	}

	err = cc.uc.Remove(c.Request.Context(), removeFromCart.ProductIds, user)
	if err != nil {
		abortJson(c, http.StatusBadRequest, "invalid request")
		return
	}

	jsonOK(c)
}

// GetProducts godoc
// @Summary      Get products
// @Tags         cart
// @Produce	     json
// @Success      200 {object} []model.ProductFull
// @Failure      400
// @Router       /cart/products [get]
func (cc *Cart) GetProducts(c *gin.Context) {
	user, ok := oauth.ExtractUser(c)
	if !ok {
		abortJson(c, http.StatusUnauthorized, "no valid token")
		return
	}

	if !user.HasRole("user") {
		abortNotAllowed(c)
		return
	}

	products, err := cc.uc.GetProducts(c.Request.Context(), user)

	if err != nil {
		abortJson(c, 400, "bad request")
		return
	}

	c.IndentedJSON(200, products)
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
