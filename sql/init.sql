DO $$ 
BEGIN 
  IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mydatabase') THEN 
    CREATE DATABASE mydatabase; 
  END IF; 
END $$; 

\c mydatabase;

CREATE TABLE redirects (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL 
    count INT DEFAULT 0
);
