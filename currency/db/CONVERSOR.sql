CREATE DATABASE converter_db;

USE converter_db;

CREATE TABLE IF NOT EXISTS converter (
    id INTEGER PRIMARY KEY,
    amount REAL,
    from_currency TEXT,
    to_currency TEXT,
    rate REAL,
    result REAL,
    creation_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);