CREATE TABLE session (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_client
        FOREIGN KEY(client_id)
            REFERENCES client(id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES usr(id)
)