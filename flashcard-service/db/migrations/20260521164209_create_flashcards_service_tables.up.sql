CREATE TABLE categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE flashcards (
    flashcard_id SERIAL PRIMARY KEY,
    title VARCHAR(255) UNIQUE NOT NULL,
    image TEXT,
    text TEXT,
    category_id INTEGER REFERENCES categories(category_id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL
);

CREATE TABLE favorite_categories (
    user_id INTEGER NOT NULL,  
    category_id INTEGER REFERENCES categories(category_id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, category_id)
);