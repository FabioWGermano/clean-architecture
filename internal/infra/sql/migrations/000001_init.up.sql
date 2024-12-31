CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    price DECIMAL(10, 2) NOT NULL,
    tax DECIMAL(10, 2) NOT NULL,
    final_price DECIMAL(10, 2) NOT NULL
);

INSERT INTO orders (id, price, tax, final_price) VALUES ('a', 100.5, 0.5, 101.0);
INSERT INTO orders (id, price, tax, final_price) VALUES ('b', 120.5, 0.15, 120.65);
INSERT INTO orders (id, price, tax, final_price) VALUES ('c', 130.5, 0.45, 130.95);