CREATE TYPE AuctionStatus AS ENUM (
    'PENDING',
    'ACTIVE',
    'ENDED',
    'COMPLETED',
    'CANCELED'
);

CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    lot_id INTEGER NOT NULL,
    start_price DECIMAL(10,2) NOT NULL,
    min_step DECIMAL(10,2) NOT NULL,
    current_price DECIMAL(10,2) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status AuctionStatus NOT NULL DEFAULT 'PENDING',
    winner_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    winner_bid_id INTEGER,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_lot FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE CASCADE,
    CONSTRAINT unique_lot UNIQUE (lot_id)
);

CREATE INDEX idx_auctions_lot_id ON auctions(lot_id);
