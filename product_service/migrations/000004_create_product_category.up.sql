CREATE TABLE product_category (
    category_id UUID NOT NULL,
    product_id UUID NOT NULL,
    PRIMARY KEY (product_id, category_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
            REFERENCES product(id),
    CONSTRAINT fk_category
        FOREIGN KEY(category_id)
            REFERENCES category(id)
)