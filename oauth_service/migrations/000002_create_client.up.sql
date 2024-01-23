CREATE TABLE client (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        secret VARCHAR(255) NOT NULL,
        domain VARCHAR(255) NOT NULL,
        type VARCHAR(255) NOT NULL
    );