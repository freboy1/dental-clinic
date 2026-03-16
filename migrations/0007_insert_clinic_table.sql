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
      ('a1b2c3d4-0001-4000-8000-000000000001', 'Healthy Life Clinic', 'Современная клиника.', '+1-202-555-0101', 'info@healthylife.com', 'https://healthylife.com', 4.5, 124, TRUE, NOW()),
      ('a1b2c3d4-0002-4000-8000-000000000002', 'Sunrise Medical Center', 'Семейная медицина.', '+1-202-555-0102', 'contact@sunrisemed.com', 'https://sunrisemed.com', 4.2, 89, TRUE, NOW()),
      ('a1b2c3d4-0003-4000-8000-000000000003', 'Green Valley Clinic', 'Диагностический центр.', '+1-202-555-0103', 'hello@greenvalley.com', 'https://greenvalley.com', 4.8, 201, TRUE, NOW()),
      ('a1b2c3d4-0004-4000-8000-000000000004', 'CityCare Hospital', 'Городская клиника с круглосуточным приемом.', '+1-202-555-0104', 'support@citycare.com', 'https://citycare.com', 4.1, 156, TRUE, NOW()),
      ('a1b2c3d4-0005-4000-8000-000000000005', 'Wellness Point', 'Центр профилактики и восстановления здоровья.', '+1-202-555-0105', 'info@wellnesspoint.com', 'https://wellnesspoint.com', 4.6, 73, TRUE, NOW()),
      ('a1b2c3d4-0006-4000-8000-000000000006', 'North Star Clinic', 'Медицинская помощь для всей семьи.', '+1-202-555-0106', 'admin@northstarclinic.com', 'https://northstarclinic.com', 4.0, 54, TRUE, NOW()),
      ('a1b2c3d4-0007-4000-8000-000000000007', 'Harmony Health', 'Инновационные методы лечения.', '+1-202-555-0107', 'contact@harmonyhealth.com', 'https://harmonyhealth.com', 4.7, 142, TRUE, NOW()),
      ('a1b2c3d4-0008-4000-8000-000000000008', 'Prime Medical', 'Клиника премиального уровня.', '+1-202-555-0108', 'info@primemedical.com', 'https://primemedical.com', 4.9, 310, TRUE, NOW()),
      ('a1b2c3d4-0009-4000-8000-000000000009', 'Family First Clinic', 'Доступная медицина для всей семьи.', '+1-202-555-0109', 'hello@familyfirst.com', 'https://familyfirst.com', 3.9, 67, TRUE, NOW()),
      ('a1b2c3d4-0010-4000-8000-000000000010', 'CarePlus Center', 'Современный медицинский центр.', '+1-202-555-0110', 'info@careplus.com', 'https://careplus.com', 4.3, 98, FALSE, NOW());

-- +goose Down
DELETE FROM clinics
WHERE id IN (
             'a1b2c3d4-0001-4000-8000-000000000001',
             'a1b2c3d4-0002-4000-8000-000000000002',
             'a1b2c3d4-0003-4000-8000-000000000003',
             'a1b2c3d4-0004-4000-8000-000000000004',
             'a1b2c3d4-0005-4000-8000-000000000005',
             'a1b2c3d4-0006-4000-8000-000000000006',
             'a1b2c3d4-0007-4000-8000-000000000007',
             'a1b2c3d4-0008-4000-8000-000000000008',
             'a1b2c3d4-0009-4000-8000-000000000009',
             'a1b2c3d4-0010-4000-8000-000000000010'
    );