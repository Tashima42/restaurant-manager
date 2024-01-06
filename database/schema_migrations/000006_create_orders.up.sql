CREATE TABLE IF NOT EXISTS orders (
  id INTEGER PRIMARY KEY,
  'item_id' INTEGER NOT NULL,
  'quantity' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('item_id') REFERENCES 'items'('id')
);
