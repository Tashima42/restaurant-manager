CREATE TABLE IF NOT EXISTS carts (
  id INTEGER PRIMARY KEY,
  'table_id' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('table_id') REFERENCES 'tables'('id')
);
