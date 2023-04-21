CREATE TABLE users(
    username VARCHAR (255) UNIQUE NOT NULL,
    firstname VARCHAR (255) NOT NULL,
    lastname VARCHAR (255) NOT NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, 
    user_id UUID PRIMARY KEY,

    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
ALTER TABLE users ALTER COLUMN user_id SET DEFAULT uuid_generate_v4();
CREATE UNIQUE INDEX users_username_idx ON users (username);