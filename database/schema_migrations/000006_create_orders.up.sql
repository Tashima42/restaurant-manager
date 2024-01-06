CREATE TABLE IF NOT EXISTS orders (
  id TEXT PRIMARY KEY,
  'item_id' TEXT NOT NULL,
  'quantity' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('item_id') REFERENCES 'items'('id')
);
