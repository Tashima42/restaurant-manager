CREATE TABLE IF NOT EXISTS menu_items (
  id INTEGER PRIMARY KEY,
  'menu_id' INTEGER NOT NULL,
  'item_id' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('menu_id') REFERENCES 'menus'('id')
  FOREIGN KEY('item_id') REFERENCES 'items'('id')
);
