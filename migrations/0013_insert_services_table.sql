-- +goose Up
INSERT INTO services (
    id,
    name,
    description,
    price,
    duration,
    clinic_id,
    is_active,
    created_at
) VALUES
      (gen_random_uuid(), 'Teeth Whitening', 'Teeth whitening using professional bleaching gel', 60000.00, 90, 'c250dc9b-240b-48f1-b331-0a8087f7a7fb', TRUE, NOW()),
      (gen_random_uuid(), 'Root Canal Treatment', 'Root canal treatment for infected tooth pulp', 90000.00, 120, 'dd9e84e8-839a-4797-b9d1-c25f4b5149db', TRUE, NOW()),
      (gen_random_uuid(), 'Dental X-Ray', 'Full mouth X-ray and dental examination', 8000.00, 20, '63e7b11f-edf2-4896-9ad2-1d0c53af2317', TRUE, NOW()),
      (gen_random_uuid(), 'Tooth Extraction', 'Complete tooth extraction procedure', 25000.00, 30, '005e561f-2b83-4883-9e00-2ff0752f2dc6', TRUE, NOW()),
      (gen_random_uuid(), 'Dental Crown', 'Custom dental crown installation', 120000.00, 75, '30b673b8-e0d5-43e3-8287-0e9d6b5e3245', TRUE, NOW()),
      (gen_random_uuid(), 'Teeth Cleaning', 'Professional teeth cleaning and plaque removal', 15000.00, 45, '24e180e2-9c8a-4b03-8f84-c592936de3a4', TRUE, NOW()),
      (gen_random_uuid(), 'Dental Filling', 'Dental filling for cavities using composite material', 35000.00, 60, '4c074585-4452-48b5-aea1-8e76cd90a880', TRUE, NOW());

-- +goose Down
DELETE FROM services
WHERE name IN (
               'Teeth Whitening',
               'Root Canal Treatment',
               'Dental X-Ray',
               'Tooth Extraction',
               'Dental Crown',
               'Teeth Cleaning',
               'Dental Filling'
    );