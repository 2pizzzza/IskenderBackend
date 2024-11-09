CREATE TABLE IF NOT EXISTS Language (
                                        code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100)
    );

CREATE TABLE IF NOT EXISTS Category (
                                        id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS CategoryTranslation (
                                                   category_id INT REFERENCES Category(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255),
    PRIMARY KEY (category_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Collection (
    id SERIAL PRIMARY KEY,
    price DECIMAL,
    isProducer BOOLEAN,
    isPainted BOOLEAN
    );

CREATE TABLE IF NOT EXISTS CollectionTranslation (
                                                     collection_id INT REFERENCES Collection(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255),
    description TEXT,
    PRIMARY KEY (collection_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Item (
                                    id SERIAL PRIMARY KEY,
                                    category_id INT REFERENCES Category(id),
    collection_id INT REFERENCES Collection(id),
    size VARCHAR(50),
    price DECIMAL,
    isProducer BOOLEAN,
    isPainted BOOLEAN
    );

CREATE TABLE IF NOT EXISTS ItemTranslation (
                                               item_id INT REFERENCES Item(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255),
    description TEXT,
    PRIMARY KEY (item_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Photo (
                                     id SERIAL PRIMARY KEY,
                                     url VARCHAR(255),
    isMain BOOLEAN
    );

CREATE TABLE IF NOT EXISTS Color (
                                     id SERIAL PRIMARY KEY,
                                     hash_color VARCHAR(7) -- Hex цвет
    );

CREATE TABLE IF NOT EXISTS CollectionPhoto (
                                               collection_id INT REFERENCES Collection(id),
    photo_id INT REFERENCES Photo(id),
    PRIMARY KEY (collection_id, photo_id)
    );

CREATE TABLE IF NOT EXISTS CollectionColor (
                                               collection_id INT REFERENCES Collection(id),
    color_id INT REFERENCES Color(id),
    PRIMARY KEY (collection_id, color_id)
    );

CREATE TABLE IF NOT EXISTS ItemPhoto (
                                         item_id INT REFERENCES Item(id),
    photo_id INT REFERENCES Photo(id),
    PRIMARY KEY (item_id, photo_id)
    );

CREATE TABLE IF NOT EXISTS ItemColor (
                                         item_id INT REFERENCES Item(id),
    color_id INT REFERENCES Color(id),
    PRIMARY KEY (item_id, color_id)
    );
