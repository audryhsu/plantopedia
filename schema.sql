CREATE DATABASE plantopedia;
\c plantopedia

CREATE TABLE plants (
  id SERIAL PRIMARY KEY,
  species VARCHAR(256),
  description VARCHAR(1024)
);

INSERT INTO plants (species , description) VALUES 
('golden pothos', 'jungle vines')
('fiddle leaf fig', 'tropic, tempermental');