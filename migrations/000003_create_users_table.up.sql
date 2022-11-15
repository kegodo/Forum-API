CREATE TABLE IF NOT EXISTS users(
    ID bigserial PRIMARY KEY,
    CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text not null,
    email citext Unique not null, 
    password_hash bytea not null,
    activated bool NOT Null,
    version integer NOT NULL DEFAULT 1
);