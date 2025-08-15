ALTER TABLE subscriptions ADD CONSTRAINT sub_price_check CHECK (price >= 0);


