
CREATE TYPE order_status AS ENUM ('open', 'close', 'rejected');

CREATE TABLE inventory (
    ingredient_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    quantity FLOAT NOT NULL,
    price INTEGER ,
    unit VARCHAR(20) NOT NULL
);

CREATE TABLE menu_items (
    product_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50),
    allergens TEXT[]
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name VARCHAR(100) NOT NULL,
    status order_status NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT REFERENCES menu_items(product_id) ON DELETE CASCADE,
    quantity INT NOT NULL,
    item_details JSONB
);

CREATE TABLE menu_item_ingredients (
    product_id INT REFERENCES menu_items(product_id) ON DELETE CASCADE,
    ingredient_id INT REFERENCES inventory(ingredient_id) ON DELETE CASCADE,
    quantity FLOAT NOT NULL,
    PRIMARY KEY (product_id, ingredient_id)
);

CREATE TABLE price_history (
    history_id SERIAL PRIMARY KEY,
    product_id INT REFERENCES menu_items(product_id) ON DELETE CASCADE,
    old_price DECIMAL(10, 2) NOT NULL,
    new_price DECIMAL(10, 2) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_status_history (
    history_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    status order_status NOT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    ingredient_id INT REFERENCES inventory(ingredient_id) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
--Автоматическое создание записи в price_history при обновлении цены товара в таблице menu_items.
CREATE OR REPLACE FUNCTION log_price_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.price <> OLD.price THEN
        INSERT INTO price_history (product_id, old_price, new_price, updated_at)
        VALUES (OLD.product_id, OLD.price, NEW.price, CURRENT_TIMESTAMP);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER price_change_trigger
AFTER UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION log_price_change();

--Автоматическое создание записи в order_status_history при изменении статуса заказа.
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status <> OLD.status THEN
        INSERT INTO order_status_history (order_id, status, changed_at)
        VALUES (OLD.order_id, NEW.status, CURRENT_TIMESTAMP);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER order_status_change_trigger
AFTER UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_status_change();

--Если в inventory вводится отрицательное значение для quantity, операция будет отклонена.
CREATE OR REPLACE FUNCTION check_inventory_quantity()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity < 0 THEN
        RAISE EXCEPTION 'Quantity in inventory cannot be negative';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_negative_inventory
BEFORE INSERT OR UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION check_inventory_quantity();


--Проверка на отрицательные цены в menu_items
CREATE OR REPLACE FUNCTION check_menu_item_price()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.price <= 0 THEN
        RAISE EXCEPTION 'Price of menu item must be greater than zero';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_negative_price
BEFORE INSERT OR UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION check_menu_item_price();

--Проверка на нулевое количество в order_items

CREATE OR REPLACE FUNCTION check_order_item_quantity()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity <= 0 THEN
        RAISE EXCEPTION 'Order item quantity must be greater than zero';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prevent_zero_quantity
BEFORE INSERT OR UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION check_order_item_quantity();


CREATE INDEX idx_orders_id ON orders(order_id);

CREATE INDEX idx_orders_status ON orders(status);