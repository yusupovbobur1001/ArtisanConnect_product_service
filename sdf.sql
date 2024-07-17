INSERT INTO order_items (order_id, product_id, total_amount, price) VALUES
('9bc8a809-920d-475c-8271-fa73b5c3e35d', 'f7999ed7-58e5-4697-b91c-2cddc6e1f1e1', 3.0, 15.99),
('f48a02ab-ea9f-4ca1-abcf-311cb8a82a9b', 'f7999ed7-58e5-4697-b91c-2cddc6e1f1e1', 2.0, 45.50),
('9bc8a809-920d-475c-8271-fa73b5c3e35d', '838eeb8d-ed26-4af9-b17c-a2444b5b75aa', 1.0, 23.75);



INSERT INTO ratings (product_id, user_id, rating, comment) VALUES
('f7999ed7-58e5-4697-b91c-2cddc6e1f1e1', '7b7f82b3-6c48-4111-9f6a-3e56d4d44d5b', 4.5, 'Great product, highly recommended!'),
('838eeb8d-ed26-4af9-b17c-a2444b5b75aa', '1a1d9cbe-3c4b-431c-8e45-4c2ef1c1b672', 3.0, 'Average quality, could be better.'),
('f7999ed7-58e5-4697-b91c-2cddc6e1f1e1', '4f6e1f9c-6c73-4114-8b35-3f67b1d2e4c6', 5.0, 'Excellent, very satisfied with the purchase!');
