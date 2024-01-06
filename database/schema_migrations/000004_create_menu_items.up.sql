CREATE TABLE IF NOT EXISTS menu_items (
  id TEXT PRIMARY KEY,
  'menu_id' TEXT NOT NULL,
  'item_id' TEXT NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('menu_id') REFERENCES 'menus'('id')
  FOREIGN KEY('item_id') REFERENCES 'items'('id')
);
