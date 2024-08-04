CREATE TABLE IF NOT EXISTS database (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    host TEXT NOT NULL,
    port INTEGER NOT NULL,
    db_user TEXT NOT NULL,
    password TEXT,
    db_name TEXT NOT NULL,
    driver TEXT NOT NULL
);
