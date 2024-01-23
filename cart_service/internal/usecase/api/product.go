package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"store/cart_service/internal/model"
	"store/cart_service/internal/model/dto"

	"github.com/google/uuid"
)

type Product struct {
	url string
}

func NewProduct(baseUrl string) (product *Product, err error) {
	_, err = url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &Product{url: baseUrl}, nil
}

func (p *Product) ExistBatch(ctx context.Context, productIds []uuid.UUID) (products []model.Product, err error) {
	reqUrl, err := url.JoinPath(p.url, "/products/batch")
	if err != nil {
		return nil, err
	}
	jsonBody, err := json.Marshal(&dto.GetBatchRequest{ProductIds: productIds})
	if err != nil {
		return nil, err
	}
	response, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var prods dto.CheckExistResponse

	err = json.Unmarshal(body, &prods)
	if err != nil {
		return nil, err
	}

	return prods.Products, nil
}

func (p *Product) GetBatch(ctx context.Context, productIds []uuid.UUID) (products []model.ProductFull, err error) {
	reqUrl, err := url.JoinPath(p.url, "/products/batch")
	if err != nil {
		return nil, err
	}
	jsonBody, err := json.Marshal(&dto.GetBatchRequest{ProductIds: productIds})
	if err != nil {
		return nil, err
	}
	response, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var prods dto.GetBatchResponse

	err = json.Unmarshal(body, &prods)
	if err != nil {
		return nil, err
	}

	return prods.Products, nil
}
