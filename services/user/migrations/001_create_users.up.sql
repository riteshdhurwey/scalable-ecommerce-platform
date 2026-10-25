CREATE TABLE users (
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_name VARCHAR(255) NOT NULL,
    user_id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,

    CONSTRAINT users_pkey PRIMARY KEY (user_name),
    CONSTRAINT users_email_key UNIQUE (email),
    CONSTRAINT users_user_id_key UNIQUE (user_id)
);