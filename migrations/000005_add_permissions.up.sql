--Filenmae: migrations/000005_add_permiossions.up.sql
CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY, 
    code text NOT NULL
);

--create a linking table that links users to permissions
--Many to many relationship

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY  KEY(user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES ('forum:read'), ('forum:write');