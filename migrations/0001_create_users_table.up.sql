CREATE TABLE IF NOT EXISTS users(
    username VARCHAR (255) UNIQUE NOT NULL,
    firstname VARCHAR (255) NOT NULL,
    lastname VARCHAR (255) NOT NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, 
    user_id SERIAL PRIMARY KEY,

    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE UNIQUE INDEX IF NOT EXISTS users_user_id_idx ON users (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
