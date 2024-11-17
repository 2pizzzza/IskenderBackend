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
    price DECIMAL NOT NULL,
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
    price DECIMAL NOT NULL,
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS Discount (
    id SERIAL PRIMARY KEY,
    discount_type VARCHAR(50) NOT NULL,
    target_id INT NOT NULL,
    discount_percentage DECIMAL(5, 2) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    CONSTRAINT fk_target_collection FOREIGN KEY (target_id) REFERENCES Collection(id) ON DELETE CASCADE,
    CONSTRAINT fk_target_item FOREIGN KEY (target_id) REFERENCES Item(id) ON DELETE CASCADE
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

