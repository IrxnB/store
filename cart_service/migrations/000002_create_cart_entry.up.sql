CREATE TABLE cart_entry (
    user_id UUID,
    product_id UUID,
    PRIMARY KEY (product_id, user_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
            REFERENCES product_cache(id)
)