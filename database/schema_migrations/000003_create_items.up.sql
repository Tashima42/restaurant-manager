CREATE TABLE IF NOT EXISTS items (
  id INTEGER PRIMARY KEY,
  'name' TEXT NOT NULL,
  'description' TEXT NOT NULL,
  'picture' TEXT NOT NULL,
  'price' INTEGER NOT NULL,
  'created_at' DATE NOT NULL,
  'updated_at' DATE NOT NULL
);
