CREATE TABLE role_user (
    role_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (role_id, user_id),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES usr(id),
    CONSTRAINT fk_role
        FOREIGN KEY(role_id)
            REFERENCES role(id)
)