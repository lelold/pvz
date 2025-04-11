CREATE TABLE products (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    type TEXT NOT NULL,
    reception_id UUID REFERENCES receptions(id)
);
    