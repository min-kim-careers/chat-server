CREATE TABLE IF NOT EXISTS rooms (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  item_id TEXT NOT NULL,
  client1 UUID NOT NULL,
  client2 UUID NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  CONSTRAINT uniq_item_clients UNIQUE (item_id, client1, client2)
);