CREATE TABLE IF NOT EXISTS carts (
  id TEXT PRIMARY KEY,
  'table_id' TEXT NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL,

  FOREIGN KEY('table_id') REFERENCES 'tables'('id')
);
