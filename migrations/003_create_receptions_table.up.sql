CREATE TABLE receptions (
    id UUID PRIMARY KEY,
    date_time TIMESTAMP NOT NULL,
    pvz_id UUID REFERENCES pvzs(id),
    status TEXT NOT NULL
);
