CREATE TABLE product (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    price NUMERIC NOT NULL,
    seller_id UUID NOT NULL,
    CONSTRAINT fk_seller
        FOREIGN KEY (seller_id)
            REFERENCES seller(id)
);