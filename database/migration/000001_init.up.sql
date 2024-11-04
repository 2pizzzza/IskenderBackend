CREATE TABLE IF NOT EXISTS Category (
                          category_id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Item (
                      item_id SERIAL PRIMARY KEY,
                      category_id INT NOT NULL REFERENCES Category(category_id) ON DELETE CASCADE,
                      name VARCHAR(255) NOT NULL,
                      description TEXT,
                      price DECIMAL(10, 2) NOT NULL,
                      is_produced BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Color (
                       color_id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS Item_Color (
                            item_id INT NOT NULL REFERENCES Item(item_id) ON DELETE CASCADE,
                            color_id INT NOT NULL REFERENCES Color(color_id) ON DELETE CASCADE,
                            PRIMARY KEY (item_id, color_id)
);

CREATE TABLE IF NOT EXISTS Photo (
                       photo_id SERIAL PRIMARY KEY,
                       item_id INT NOT NULL REFERENCES Item(item_id) ON DELETE CASCADE,
                       url VARCHAR(255) NOT NULL,
                       is_main BOOLEAN DEFAULT FALSE
);
