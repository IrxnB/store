package repo

import (
	"context"
	"math/big"
	"product_service/internal/model"
	"product_service/internal/model/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) (productRepo *Product, err error) {
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Product{db: db}, nil
}

func (p *Product) GetById(ctx context.Context, id uuid.UUID) (product *model.Product, err error) {
	rows, err := p.db.Query(ctx, `
		SELECT p.id, p.name, p.desc, p.price, 
		s.id as seller_id, s.name as seller_name, 
		c.id as category_id, c.name as category_name
		FROM product p
		LEFT JOIN product_category pc
		ON pc.product_id = p.id
		JOIN categoty c 
		ON pc.category_id = c.id
		JOIN seller s
		ON s.id = p.seller_id
		WHERE p.id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	fromDb := model.Product{Id: uuid.Nil, Categories: make([]model.Category, 0, 1)}

	for rows.Next() {
		var priceStr string
		var cur ProductWithCategoryRow
		err = rows.Scan(&cur.Id, &cur.Name, &cur.Description, &priceStr,
			&cur.Seller.OwnerId, &cur.Seller.Name,
			&cur.Category.Id, &cur.Category.Name)

		if err != nil {
			return nil, err
		}
		if fromDb.Id == uuid.Nil {
			fromDb.Id = cur.Id
			fromDb.Name = cur.Name
			fromDb.Description = cur.Description
			fromDb.Seller = cur.Seller
			price := big.Float{}
			price.SetString(priceStr)
			fromDb.Price = &price
		}
		if fromDb.Id != uuid.Nil {
			fromDb.Categories = append(fromDb.Categories, cur.Category)
		}
	}

	return &fromDb, nil
}

func (p *Product) GetBatch(ctx context.Context, ids []uuid.UUID) (product []dto.GetBatchItem, err error) {
	rows, err := p.db.Query(ctx, `
		SELECT p.id, p.name, p.price, p.Descriprion, s.name
		FROM product p
		JOIN seller s
		ON s.id = p.seller_id
		WHERE p.id = ANY($1)
	`, ids)

	if err != nil {
		return nil, err
	}

	batch := make([]dto.GetBatchItem, 0, len(ids))

	for rows.Next() {
		var cur dto.GetBatchItem
		var priceStr string
		err = rows.Scan(&cur.Id, &cur.Name, &priceStr, &cur.Description,
			&cur.SellerName)

		if err != nil {
			return nil, err
		}

		price := big.Float{}
		price.SetString(priceStr)

		cur.Price = &price

		batch = append(batch, cur)
	}

	return batch, nil
}

func (p *Product) Save(ctx context.Context, createProduct dto.CreateProduct, sellerId uuid.UUID) (err error) {
	trx, err := p.db.Begin(ctx)

	if err != nil {
		return err
	}
	defer trx.Commit(ctx)

	id := uuid.New()

	_, err = trx.Exec(ctx, `
		INSERT INTO product(id, name, description, price, seller_id)
		VALUES ($1, $2, $3, $4, $5)
	`, id, createProduct.Name, createProduct.Description, createProduct.Price.String(), sellerId)

	if err != nil {
		trx.Rollback(ctx)
		return err
	}

	_, err = trx.Exec(ctx, `
		INSERT INTO product_category(product_id, category_id)
		SELECT $1, id FROM category WHERE id = ANY($2)
	`, id, createProduct.CategoryIds)
	if err != nil {
		trx.Rollback(ctx)
		return err
	}
	return nil
}

func (p *Product) Update(ctx context.Context, product model.Product) (err error) {
	panic("not implemented")
}

func (p *Product) GetPage(ctx context.Context, page, limit int) (products []model.Product, err error) {
	offset := (page - 1) * limit

	rows, err := p.db.Query(ctx, `
		SELECT p.id, p.name, p.description, p.price, 
		s.id as seller_id, s.name as seller_name
		FROM product p 
		JOIN seller s
		ON s.id = p.seller_id
		ORDER BY p.Id
		LIMIT $1
		OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, err
	}

	products = make([]model.Product, 0, limit)

	for rows.Next() {
		var cur model.Product
		var priceStr string
		err = rows.Scan(&cur.Id, &cur.Name, &cur.Description, &priceStr,
			&cur.Seller.OwnerId, &cur.Seller.Name)

		if err != nil {
			return nil, err
		}

		price := big.Float{}
		price.SetString(priceStr)

		cur.Price = &price

		products = append(products, cur)
	}

	return products, nil
}

func (p *Product) GetPageBySellerId(ctx context.Context, page, limit int, id uuid.UUID) (products []model.Product, err error) {
	offset := (page - 1) * limit

	rows, err := p.db.Query(ctx, `
		SELECT p.id, p.name, p.description, p.price, 
		s.id as seller_id, s.name as seller_name
		FROM product p 
		JOIN seller s
		ON s.id = p.seller_id
		WHERE s.id = $1
		ORDER BY p.Id
		LIMIT $2
		OFFSET $3
	`, id, limit, offset)

	if err != nil {
		return nil, err
	}

	products = make([]model.Product, 0, limit)

	for rows.Next() {
		var cur model.Product
		var priceStr string
		err = rows.Scan(&cur.Id, &cur.Name, &cur.Description, &priceStr,
			&cur.Seller.OwnerId, &cur.Seller.Name)

		if err != nil {
			return nil, err
		}

		price := big.Float{}
		price.SetString(priceStr)

		cur.Price = &price

		products = append(products, cur)
	}

	return products, nil
}

type ProductWithCategoryRow struct {
	Id          uuid.UUID
	Name        string
	Description string
	Price       *big.Float
	Seller      model.Seller
	Category    model.Category
}
