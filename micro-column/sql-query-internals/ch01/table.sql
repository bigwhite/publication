-- 1. 建表
CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(50),
    city VARCHAR(50)
);

CREATE TABLE orders (
    id INT PRIMARY KEY,
    uid INT,
    amount INT
);

-- 2. 插入我们在推演中使用的 3 个用户
INSERT INTO users VALUES
(1, 'Alice', 'Beijing'),
(2, 'Bob', 'Shanghai'),
(3, 'Charlie', 'Beijing');

-- 3. 插入我们在推演中使用的 5 个订单
INSERT INTO orders VALUES
(101, 1, 50),   -- Alice 的小额订单，会被 WHERE 丢弃
(102, 1, 150),  -- Alice 的大额订单
(103, 2, 500),  -- Bob 的超大额订单
(104, 3, 100),  -- Charlie 的订单 1
(105, 3, 120);  -- Charlie 的订单 2
