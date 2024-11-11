DELETE FROM bids;
DELETE FROM auctions;
DELETE FROM lots;
DELETE FROM users;

INSERT INTO users (username, email, balance) VALUES
('john_doe', 'john@example.com', 1000.00),
('jane_smith', 'jane@example.com', 2000.00),
('bob_wilson', 'bob@example.com', 1500.00),
('alice_brown', 'alice@example.com', 3000.00),
('mike_davis', 'mike@example.com', 2500.00),
('sarah_miller', 'sarah@example.com', 1800.00),
('david_taylor', 'david@example.com', 3500.00),
('emma_wilson', 'emma@example.com', 2200.00),
('james_moore', 'james@example.com', 1700.00),
('lisa_anderson', 'lisa@example.com', 2800.00);

INSERT INTO lots (title, description, start_price, creator_id) VALUES
('Vintage Rolex Watch', 'Rare collectible watch from 1950s', 5000.00, 1),
('PlayStation 5', 'Brand new gaming console with 2 controllers', 500.00, 2),
('Antique Chinese Vase', 'Beautiful Ming dynasty vase', 10000.00, 3),
('Classic Ferrari', 'Restored 1960s Ferrari in perfect condition', 150000.00, 4),
('Diamond Ring', '3 carat diamond ring with platinum band', 15000.00, 5),
('Rare Comic Book', 'First edition Superman comic', 8000.00, 6),
('Modern Art Painting', 'Original artwork by contemporary artist', 12000.00, 7),
('Vintage Guitar', '1970s Fender Stratocaster', 6000.00, 8),
('Luxury Watch', 'Limited edition Patek Philippe', 25000.00, 9),
('Ancient Coin Collection', 'Set of rare Roman coins', 20000.00, 10),
('First Edition Book', 'Signed first edition of classic novel', 3000.00, 1),
('Sports Memorabilia', 'Signed jersey by legendary player', 1500.00, 2);

INSERT INTO auctions (
    lot_id, 
    start_price,
    current_price,
    min_step,
    start_time, 
    end_time
) VALUES
(1, 5000.00, 5000.00, 100.00, NOW(), NOW() + INTERVAL '7 days'),
(2, 500.00, 500.00, 25.00, NOW(), NOW() + INTERVAL '5 days'),
(3, 10000.00, 10000.00, 500.00, NOW(), NOW() + INTERVAL '10 days'),
(4, 150000.00, 150000.00, 5000.00, NOW(), NOW() + INTERVAL '14 days'),
(5, 15000.00, 15000.00, 500.00, NOW(), NOW() + INTERVAL '6 days'),
(6, 8000.00, 8000.00, 200.00, NOW(), NOW() + INTERVAL '8 days'),
(7, 12000.00, 12000.00, 500.00, NOW(), NOW() + INTERVAL '9 days'),
(8, 6000.00, 6000.00, 200.00, NOW(), NOW() + INTERVAL '7 days'),
(9, 25000.00, 25000.00, 1000.00, NOW(), NOW() + INTERVAL '12 days'),
(10, 20000.00, 20000.00, 1000.00, NOW(), NOW() + INTERVAL '11 days');

INSERT INTO bids (auction_id, user_id, amount) VALUES
(1, 2, 5200.00),
(1, 3, 5500.00),
(1, 4, 6000.00),
(1, 5, 6500.00),

(2, 1, 550.00),
(2, 3, 600.00),
(2, 5, 650.00),
(2, 7, 700.00),

(3, 4, 10500.00),
(3, 6, 11000.00),
(3, 8, 12000.00),

(4, 5, 155000.00),
(4, 7, 160000.00),
(4, 9, 165000.00),

(5, 1, 15500.00),
(5, 3, 16000.00),
(5, 5, 16500.00),

(6, 2, 8500.00),
(6, 4, 9000.00),
(7, 3, 12500.00),
(7, 5, 13000.00),
(8, 6, 6500.00),
(8, 8, 7000.00),
(9, 7, 26000.00),
(9, 9, 27000.00),
(10, 8, 21000.00),
(10, 10, 22000.00);

UPDATE auctions a
SET current_price = (
    SELECT amount
    FROM bids b
    WHERE b.auction_id = a.id
    ORDER BY amount DESC
    LIMIT 1
)
WHERE EXISTS (
    SELECT 1
    FROM bids b
    WHERE b.auction_id = a.id
);

UPDATE auctions 
SET 
    winner_id = (
        SELECT user_id 
        FROM bids 
        WHERE auction_id = auctions.id 
        ORDER BY amount DESC 
        LIMIT 1
    ),
    winner_bid_id = (
        SELECT id 
        FROM bids 
        WHERE auction_id = auctions.id 
        ORDER BY amount DESC 
        LIMIT 1
    )
WHERE id IN (1, 2, 3);
