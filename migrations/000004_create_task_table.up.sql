CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    userId INT REFERENCES users(id),
    taskName VARCHAR(100),
    startTime TIMESTAMP,
    endTime TIMESTAMP
);