-- Create Logins table
CREATE TABLE Logins (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER UNIQUE NOT NULL,
  password TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES Users(id)
);
