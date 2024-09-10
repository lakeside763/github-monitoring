CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE repositories (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id INT UNIQUE,
  name VARCHAR(255) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  description TEXT,
  url VARCHAR(255) NOT NULL,
  language VARCHAR(50),
  forks_count INT,
  stars_count INT,
  open_issues_count INT,
  watchers_count INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS commits (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  repo_id INT NOT NULL,
  sha VARCHAR(40) NOT NULL,
  message TEXT NOT NULL,
  author VARCHAR(255) NOT NULL,
  date TIMESTAMP NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (repo_id) REFERENCES repositories(repo_id)
);