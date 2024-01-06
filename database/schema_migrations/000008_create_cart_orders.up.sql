CREATE TABLE IF NOT EXISTS cart_orders (
  id INTEGER PRIMARY KEY,
  'cart_id' INTEGER NOT NULL,
  'order_id' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('cart_id') REFERENCES 'carts'('id')
  FOREIGN KEY('order_id') REFERENCES 'orders'('id')
);
