CREATE TABLE usr (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        username VARCHAR(255) UNIQUE NOT NULL,
        password VARCHAR(255) NOT NULL
    );
CREATE INDEX idx_usr_username ON usr (username);
    
    