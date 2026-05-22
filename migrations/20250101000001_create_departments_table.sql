-- +goose Up
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id INT REFERENCES departments(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_departments_parent_id ON departments(parent_id);
ALTER TABLE departments ADD CONSTRAINT unique_name_per_parent UNIQUE (parent_id, name);

-- +goose Down
DROP TABLE departments;