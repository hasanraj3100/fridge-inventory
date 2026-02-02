-- +migrate UP

CREATE TABLE IF NOT EXISTS item_usage (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    quantity_used FLOAT NOT NULL, 
    reason TEXT NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (item_id) REFERENCES fridge_items(id) ON DELETE CASCADE
);

CREATE INDEX idx_item_usage_item_id ON item_usage(item_id);

