CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE base_table (
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_account (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phoneno BIGINT UNIQUE NOT NULL,
    storename VARCHAR(255) UNIQUE NOT NULL,
    storeaddress TEXT NOT NULL,
    pincode INTEGER NOT NULL
) INHERITS (base_table);

-- CREATE TABLE store (
--     id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
--     store_id uuid,
--     productname VARCHAR(255) NOT NULL,
--     FOREIGN KEY (store_id) REFERENCES user_account (id),
--     FOREIGN KEY (store_name) REFERENCES user_account (storename))
-- ) INHERITS (base_table);

CREATE TABLE product (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    product_id uuid,
    product_name VARCHAR(255) NOT NULL,
    descript TEXT NOT NULL,
    product_location VARCHAR(255),
    FOREIGN KEY (product_id) REFERENCES user_account (id) ON DELETE CASCADE
) INHERITS (base_table);


CREATE TABLE category (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v1(),
    cat_id uuid,
    cat_name VARCHAR(255) NOT NULL,
    FOREIGN KEY (cat_id) REFERENCES stores (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
) INHERITS (base_table);