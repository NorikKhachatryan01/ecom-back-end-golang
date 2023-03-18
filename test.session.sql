CREATE TABLE customers (
    ID INT NOT NULL,
    Name VARCHAR(255),
    Email VARCHAR(255),
    LastLoginTime DATETIME,
    PagesVisited TEXT,
    ItemsInCart TEXT,
    PRIMARY KEY (ID)
);

CREATE TABLE products (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE categories (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE orders (
    ID INT NOT NULL,
    CustomerID INT,
    OrderDate DATE,
    TotalAmount DECIMAL(10, 2),
    CreatedAt DATETIME,
    SessionID VARCHAR(255),
    PRIMARY KEY (ID)
);

CREATE TABLE payments (
    id INT NOT NULL AUTO_INCREMENT,
    customer_id INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES customers (id)
);
