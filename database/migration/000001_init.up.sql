CREATE TABLE IF NOT EXISTS Language (
    code VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100)
    );

CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   username VARCHAR(255) UNIQUE NOT NULL,
   password VARCHAR(255) NOT NULL,
   created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS Category (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS CategoryTranslation (
    category_id INT REFERENCES Category(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255) NOT NULL,
    PRIMARY KEY (category_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Collection (
    id SERIAL PRIMARY KEY,
    price DECIMAL DEFAULT 0,
    isProducer BOOLEAN DEFAULT false,
    isPainted BOOLEAN DEFAULT false,
    isPopular BOOLEAN DEFAULT false,
    isNew BOOLEAN DEFAULT false
    );

CREATE TABLE IF NOT EXISTS CollectionTranslation (
    collection_id INT REFERENCES Collection(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    PRIMARY KEY (collection_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Vacancy(
    id SERIAL PRIMARY KEY,
    isActive BOOLEAN DEFAULT true,
    salary DECIMAL NOT NULL
);

CREATE TABLE IF NOT EXISTS VacancyTranslation (
    vacancy_id INT REFERENCES Vacancy(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    title VARCHAR(255) NOT NULL,
    requirements TEXT[] NOT NULL,
    responsibilities TEXT[] NOT NULL,
    conditions TEXT[] NOT NULL,
    information TEXT[] NOT NULL
    );


CREATE TABLE IF NOT EXISTS Item (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL REFERENCES Category(id),
    collection_id INT REFERENCES Collection(id),
    size VARCHAR(50) NOT NULL,
    price DECIMAL DEFAULT 0,
    isProducer BOOLEAN DEFAULT false,
    isPainted BOOLEAN DEFAULT false,
    isPopular BOOLEAN DEFAULT false,
    isNew BOOLEAN DEFAULT false
    );


CREATE TABLE IF NOT EXISTS ItemTranslation (
    item_id INT REFERENCES Item(id),
    language_code VARCHAR(10) REFERENCES Language(code),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    PRIMARY KEY (item_id, language_code)
    );

CREATE TABLE IF NOT EXISTS Brand(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL
    );

CREATE TABLE IF NOT EXISTS Review (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    text TEXT,
    isShow BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS Discount (
    id SERIAL PRIMARY KEY,
    discount_type VARCHAR(50) NOT NULL CHECK (discount_type IN ('collection', 'item')),
    target_id INT NOT NULL,
    discount_percentage DECIMAL(5, 2) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
    );


CREATE TABLE IF NOT EXISTS Photo (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255),
    isMain BOOLEAN,
    hash_color VARCHAR(7)
    );

CREATE TABLE IF NOT EXISTS CollectionPhoto (
    collection_id INT REFERENCES Collection(id),
    photo_id INT REFERENCES Photo(id),
    PRIMARY KEY (collection_id, photo_id)
    );


CREATE TABLE IF NOT EXISTS ItemPhoto (
    item_id INT REFERENCES Item(id),
    photo_id INT REFERENCES Photo(id),
    PRIMARY KEY (item_id, photo_id)
    );


CREATE OR REPLACE FUNCTION validate_discount_type()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.discount_type = 'collection' THEN
        IF NOT EXISTS (SELECT 1 FROM Collection WHERE id = NEW.target_id) THEN
            RAISE EXCEPTION 'Collection with id % does not exist', NEW.target_id;
END IF;
    ELSIF NEW.discount_type = 'item' THEN
        IF NOT EXISTS (SELECT 1 FROM Item WHERE id = NEW.target_id) THEN
            RAISE EXCEPTION 'Item with id % does not exist', NEW.target_id;
END IF;
ELSE
        RAISE EXCEPTION 'Invalid discount_type: %', NEW.discount_type;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER discount_target_validation
    BEFORE INSERT ON Discount
    FOR EACH ROW
    EXECUTE FUNCTION validate_discount_type();

