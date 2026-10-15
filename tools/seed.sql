INSERT INTO events (id, title, total_tickets, available_tickets, price_cents)
VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Exit Festival 2026', 1000, 1000, 7500),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Coldplay Live in Belgrade', 500, 500, 12000),
    ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Tech Conference Niš', 200, 200, 3000)
ON CONFLICT (id) DO NOTHING;