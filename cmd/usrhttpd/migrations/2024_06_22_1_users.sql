-- Creates the users table

CREATE TABLE `users` (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    email_confirmed INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
