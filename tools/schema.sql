CREATE TABLE IF NOT EXISTS events (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      title VARCHAR(255) NOT NULL,
                                      total_tickets INT NOT NULL CHECK (total_tickets >= 0),
                                      available_tickets INT NOT NULL CHECK (available_tickets >= 0),
                                      price_cents INT NOT NULL CHECK (price_cents >= 0),
                                      created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                      event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                      customer_email VARCHAR(255) NOT NULL,
                                      quantity INT NOT NULL CHECK (quantity > 0),
                                      total_amount_cents INT NOT NULL CHECK (total_amount_cents >= 0),
                                      status VARCHAR(50) NOT NULL DEFAULT 'COMPLETED', -- 'PENDING', 'COMPLETED', 'FAILED'
                                      created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tickets (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
                                       event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                       ticket_code VARCHAR(64) UNIQUE NOT NULL,
                                       issued_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_event_id ON orders(event_id);
CREATE INDEX idx_tickets_event_id ON tickets(event_id);