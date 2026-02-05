CREATE TABLE company (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_company_id INTEGER REFERENCES company(id) ON DELETE SET NULL,
    sector VARCHAR(255)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    company_id INTEGER REFERENCES company(id) ON DELETE SET NULL,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    upvotes INTEGER NOT NULL DEFAULT 0 CHECK (upvotes >= 0),
    downvotes INTEGER NOT NULL DEFAULT 0 CHECK (downvotes >= 0)
);

CREATE TABLE comment (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    post_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    upvotes INTEGER NOT NULL DEFAULT 0 CHECK (upvotes >= 0),
    downvotes INTEGER NOT NULL DEFAULT 0 CHECK (downvotes >= 0)
);

-- dummy users
INSERT INTO users (username, email, password) VALUES
('brennan', 'brennan@example.com', 'password123'),
('alice', 'alice@example.com', 'alicepwd'),
('bob', 'bob@example.com', 'bobpwd'),
('carol', 'carol@example.com', 'carolpwd'),
('dave', 'dave@example.com', 'davepwd');

-- parent companies
INSERT INTO company (name, description, parent_company_id, sector) VALUES
('Acme Corp', 'Global manufacturer of industrial goods', NULL, 'Manufacturing'),
('Globex Corporation', 'Technology and services company', NULL, 'Technology'),
('Initech', 'Enterprise software and services', NULL, 'Software'),
('Umbrella Corporation', 'Pharmaceuticals and biotech research', NULL, 'Pharmaceuticals'),
('Wayne Enterprises', 'Diversified multinational conglomerate', NULL, 'Conglomerate');

-- subsidiaries (use subqueries to link parent_company_id)
INSERT INTO company (name, description, parent_company_id, sector) VALUES
('Acme Parts', 'Parts and components division of Acme', (SELECT id FROM company WHERE name = 'Acme Corp'), 'Manufacturing'),
('Globex AI', 'AI research lab under Globex', (SELECT id FROM company WHERE name = 'Globex Corporation'), 'Artificial Intelligence'),
('Initech Payroll', 'Payroll and HR services for Initech', (SELECT id FROM company WHERE name = 'Initech'), 'HR Services'),
('Wayne R&D', 'Research subsidiary of Wayne Enterprises', (SELECT id FROM company WHERE name = 'Wayne Enterprises'), 'Research');