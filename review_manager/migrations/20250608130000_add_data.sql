-- +goose Up

-- Добавление начальных пользователей
INSERT INTO users (created_at, email) VALUES
  (NOW() - INTERVAL '10 days', 'alice.williams@example.com'),
  (NOW() - INTERVAL '9 days',  'bob.johnson@example.com'),
  (NOW() - INTERVAL '8 days',  'charlie.brown@example.com'),
  (NOW() - INTERVAL '7 days',  'diana.smith@example.com'),
  (NOW() - INTERVAL '6 days',  'eric.miller@example.com'),
  (NOW() - INTERVAL '5 days',  'fiona.davis@example.com'),
  (NOW() - INTERVAL '4 days',  'george.moore@example.com'),
  (NOW() - INTERVAL '3 days',  'hannah.taylor@example.com'),
  (NOW() - INTERVAL '2 days',  'ian.anderson@example.com'),
  (NOW() - INTERVAL '1 day',   'julia.thomas@example.com');

-- Добавление начальных кафе и магазинов
INSERT INTO objects (created_at, name) VALUES
  (NOW() - INTERVAL '10 days', 'Café Mocha'),
  (NOW() - INTERVAL '9 days',  'Sunrise Bakery'),
  (NOW() - INTERVAL '8 days',  'Green Market Grocery'),
  (NOW() - INTERVAL '7 days',  'The Coffee Corner'),
  (NOW() - INTERVAL '6 days',  'Urban Deli'),
  (NOW() - INTERVAL '5 days',  'Books & Beans Café'),
  (NOW() - INTERVAL '4 days',  'FreshMart Supermarket'),
  (NOW() - INTERVAL '3 days',  'The Tea House'),
  (NOW() - INTERVAL '2 days',  'Little Italy Pizza'),
  (NOW() - INTERVAL '1 day',   'Cozy Cup Café');