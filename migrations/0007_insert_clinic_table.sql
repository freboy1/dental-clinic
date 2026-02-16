-- +goose Up
INSERT INTO clinics (
    id,
    name,
    description,
    phone,
    email,
    website,
    rating,
    reviews_count,
    is_active,
    created_at
) VALUES
(gen_random_uuid(), 'Healthy Life Clinic', 'Современная клиника.', '+1-202-555-0101', 'info@healthylife.com', 'https://healthylife.com', 4.5, 124, TRUE, NOW()),
(gen_random_uuid(), 'Sunrise Medical Center', 'Семейная медицина.', '+1-202-555-0102', 'contact@sunrisemed.com', 'https://sunrisemed.com', 4.2, 89, TRUE, NOW()),
(gen_random_uuid(), 'Green Valley Clinic', 'Диагностический центр.', '+1-202-555-0103', 'hello@greenvalley.com', 'https://greenvalley.com', 4.8, 201, TRUE, NOW()),
(gen_random_uuid(), 'CityCare Hospital', 'Городская клиника с круглосуточным приемом.', '+1-202-555-0104', 'support@citycare.com', 'https://citycare.com', 4.1, 156, TRUE, NOW()),
(gen_random_uuid(), 'Wellness Point', 'Центр профилактики и восстановления здоровья.', '+1-202-555-0105', 'info@wellnesspoint.com', 'https://wellnesspoint.com', 4.6, 73, TRUE, NOW()),
(gen_random_uuid(), 'North Star Clinic', 'Медицинская помощь для всей семьи.', '+1-202-555-0106', 'admin@northstarclinic.com', 'https://northstarclinic.com', 4.0, 54, TRUE, NOW()),
(gen_random_uuid(), 'Harmony Health', 'Инновационные методы лечения.', '+1-202-555-0107', 'contact@harmonyhealth.com', 'https://harmonyhealth.com', 4.7, 142, TRUE, NOW()),
(gen_random_uuid(), 'Prime Medical', 'Клиника премиального уровня.', '+1-202-555-0108', 'info@primemedical.com', 'https://primemedical.com', 4.9, 310, TRUE, NOW()),
(gen_random_uuid(), 'Family First Clinic', 'Доступная медицина для всей семьи.', '+1-202-555-0109', 'hello@familyfirst.com', 'https://familyfirst.com', 3.9, 67, TRUE, NOW()),
(gen_random_uuid(), 'CarePlus Center', 'Современный медицинский центр.', '+1-202-555-0110', 'info@careplus.com', 'https://careplus.com', 4.3, 98, FALSE, NOW());

-- +goose Down
DELETE FROM clinics
WHERE name IN (
    'Healthy Life Clinic',
    'Sunrise Medical Center',
    'Green Valley Clinic',
    'CityCare Hospital',
    'Wellness Point',
    'North Star Clinic',
    'Harmony Health',
    'Prime Medical',
    'Family First Clinic',
    'CarePlus Center'
);
