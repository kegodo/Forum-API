--File: migrations/000001_create_forum_table.up.sql
CREATE TABLE IF NOT EXISTS forums(
    ID bigserial PRIMARY KEY,
    CreatedAt timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    Title text not null,
    Category text not null, 
    Publisher text not null,
    Description text,
    ReleaseDate text not null,
    version integer NOT NULL DEFAULT 1
);