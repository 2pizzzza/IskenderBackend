CREATE TABLE IF NOT EXISTS Catalogs (
    id SERIAL PRIMARY KEY,
    price DECIMAL(10, 2) NOT NULL
    );

CREATE TABLE IF NOT EXISTS Catalogs_Localization (
    id SERIAL PRIMARY KEY,
    catalog_id INT NOT NULL REFERENCES Catalogs(id) ON DELETE CASCADE,
    language_code VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    UNIQUE(catalog_id, language_code)
    );

-- CREATE TABLE IF NOT EXISTS Category (
--                           id SERIAL PRIMARY KEY,
--                           catalog_id INT REFERENCES Catalog(id) ON DELETE CASCADE,
--                           name VARCHAR(255) NOT NULL,
--                           description TEXT
-- );
--
-- CREATE TABLE IF NOT EXISTS Collection (
--                             id SERIAL PRIMARY KEY,
--                             catalog_id INT REFERENCES Catalog(id) ON DELETE SET NULL,
--                             name VARCHAR(255) NOT NULL,
--                             description TEXT,
--                             isPainted BOOLEAN DEFAULT FALSE,
--                             isProducer BOOLEAN DEFAULT FALSE,
--                             price DECIMAL(10, 2)
-- );
--
-- CREATE TABLE IF NOT EXISTS Item (
--                       id SERIAL PRIMARY KEY,
--                       name VARCHAR(255) NOT NULL,
--                       isProducer BOOLEAN DEFAULT FALSE,
--                       description TEXT,
--                       price DECIMAL(10, 2),
--                       size VARCHAR(50)
-- );

CREATE TABLE IF NOT EXISTS Color (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL,
                       hash_color CHAR(7) NOT NULL
);

-- CREATE TABLE IF NOT EXISTS Item_color (
--                             item_id INT REFERENCES Item(id) ON DELETE CASCADE,
--                             color_id INT REFERENCES Color(id) ON DELETE CASCADE,
--                             PRIMARY KEY (item_id, color_id)
-- );

CREATE TABLE IF NOT EXISTS Catalog_Color (
                               catalog_id INT REFERENCES CatalogS(id) ON DELETE CASCADE,
                               color_id INT REFERENCES Color(id) ON DELETE CASCADE,
                               PRIMARY KEY (catalog_id, color_id)
);

-- CREATE TABLE IF NOT EXISTS Collection_Color (
--                                   collection_id INT REFERENCES Collection(id) ON DELETE CASCADE,
--                                   color_id INT REFERENCES Color(id) ON DELETE CASCADE,
--                                   PRIMARY KEY (collection_id, color_id)
-- );
--
-- CREATE TABLE IF NOT EXISTS Photo (
--                        id SERIAL PRIMARY KEY,
--                        item_id INT REFERENCES Item(id) ON DELETE CASCADE,
--                        url VARCHAR(255) NOT NULL,
--                        isMain BOOLEAN DEFAULT FALSE
-- );
--
-- CREATE TABLE IF NOT EXISTS Category_Item (
--                                category_id INT REFERENCES Category(id) ON DELETE CASCADE,
--                                item_id INT REFERENCES Item(id) ON DELETE CASCADE,
--                                PRIMARY KEY (category_id, item_id)
-- );
--
-- CREATE TABLE IF NOT EXISTS Collection_Item (
--                                  collection_id INT REFERENCES Collection(id) ON DELETE CASCADE,
--                                  item_id INT REFERENCES Item(id) ON DELETE CASCADE,
--                                  PRIMARY KEY (collection_id, item_id)
-- );
