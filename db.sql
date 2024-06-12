-- Table: users
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    deposit DECIMAL(10, 2) DEFAULT 0 CHECK(deposit>=0),
    last_login_date TIMESTAMP,
    jwt_token TEXT
);

-- Table: UserDetail
CREATE TABLE IF NOT EXISTS user_details (
    user_detail_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) UNIQUE NOT NULL,
    fname VARCHAR(255),
    lname VARCHAR(255),
    address TEXT,
    age INT DEFAULT 0 CHECK(age>=0),
    phone_number VARCHAR(20)
);

-- Table: books
CREATE TABLE IF NOT EXISTS books (
    book_id SERIAL PRIMARY KEY,
    book_name VARCHAR(255) NOT NULL,
    stock INT NOT NULL DEFAULT 0 CHECK(stock>=0),
    rental_cost DECIMAL(10, 2) DEFAULT 0 CHECK(rental_cost>=0),
    category VARCHAR(255),
    description TEXT,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255),
    CONSTRAINT unique_book_publisher UNIQUE (book_name, publisher)
);

-- Table: Rents
CREATE TABLE IF NOT EXISTS rents (
    rent_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) NOT NULL,
    book_id INT REFERENCES books(book_id) NOT NULL,
    total_price DECIMAL(10, 2),
    rent_status VARCHAR(50) NOT NULL,
    rent_at TIMESTAMP NOT NULL,
    deadline TIMESTAMP,
    returned_at TIMESTAMP
);

-- Table: Payments
CREATE TABLE IF NOT EXISTS Payments (
    payment_id SERIAL PRIMARY KEY,
    rent_id INT REFERENCES rents(rent_id) NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    payment_amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL
);


-- DML 
-- Seed data for users table
INSERT INTO users (username, email, password, role, deposit, last_login_date, jwt_token) VALUES
('john_doe', 'user1', 'admin', 'user', 50.00, '2023-06-10 12:34:56', 'token123'),
('jane_doe', 'jane@example.com', 'password456', 'user', 75.00, '2023-06-11 14:22:31', 'token456'),
('alice', 'alice@example.com', 'alicepass', 'user', 30.00, '2023-06-09 16:17:43', 'token789'),
('bob', 'bob@example.com', 'bobpass', 'user', 20.00, '2023-06-12 18:50:29', 'token101'),
('charlie', 'charlie@example.com', 'charliepass', 'user', 25.00, '2023-06-08 19:25:15', 'token102'),
('yoland', 'admin', 'admin', 'admin', 50.00, '2023-06-10 12:34:56', 'token123');

-- Seed data for user_details table
INSERT INTO user_details (user_id, fname, lname, address, age, phone_number) VALUES
(1, 'John', 'Doe', '123 Maple St', 28, '555-1234'),
(2, 'Jane', 'Doe', '456 Oak St', 32, '555-5678'),
(3, 'Alice', 'Smith', '789 Pine St', 24, '555-8765'),
(4, 'Bob', 'Brown', '101 Cedar St', 30, '555-4321'),
(5, 'Charlie', 'Davis', '202 Birch St', 27, '555-7890');

-- Seed data for books table with real-world books
INSERT INTO books (book_name, stock, rental_cost, category, description, author, publisher) VALUES
('To Kill a Mockingbird', 5, 1.99, 'Fiction', 'A novel about the serious issues of rape and racial inequality.', 'Harper Lee', 'J.B. Lippincott & Co.'),
('1984', 10, 2.99, 'Dystopian', 'A novel presenting a terrifying vision of a totalitarian future society.', 'George Orwell', 'Secker & Warburg'),
('The Great Gatsby', 2, 1.49, 'Fiction', 'A novel about the American dream and the roaring twenties.', 'F. Scott Fitzgerald', 'Charles Scribner''s Sons'),
('The Catcher in the Rye', 8, 2.49, 'Fiction', 'A story about adolescent angst and alienation.', 'J.D. Salinger', 'Little, Brown and Company'),
('Pride and Prejudice', 4, 3.99, 'Romance', 'A romantic novel of manners.', 'Jane Austen', 'T. Egerton, Whitehall'),
('The Hobbit', 12, 1.29, 'Fantasy', 'A fantasy novel about the journey of Bilbo Baggins.', 'J.R.R. Tolkien', 'George Allen & Unwin'),
('Moby Dick', 7, 2.19, 'Adventure', 'A novel about the voyage of the whaling ship Pequod.', 'Herman Melville', 'Harper & Brothers'),
('War and Peace', 3, 1.79, 'Historical Fiction', 'A novel that chronicles the French invasion of Russia.', 'Leo Tolstoy', 'The Russian Messenger'),
('The Lord of the Rings', 6, 2.69, 'Fantasy', 'An epic fantasy trilogy.', 'J.R.R. Tolkien', 'George Allen & Unwin'),
('Harry Potter and the Sorcerer''s Stone', 9, 2.99, 'Fantasy', 'A young wizard''s journey begins.', 'J.K. Rowling', 'Bloomsbury (UK), Scholastic (US)');


