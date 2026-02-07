CREATE EXTENSION IF NOT EXISTS "uuid_ossp";

CREATE TABLE pages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name TEXT NOT NULL,
  route TEXT UNIQUE NOT NULL,
  is_home BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE widgets (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  page_id UUDI NOT NULL REFERENCES page(id) ON DELETE CASCADE,
  type TEXT NOT NULL,
  position INT NOT NULL,
  config JSNOB
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);