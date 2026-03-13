-- +goose Up
INSERT INTO doctors (
    id,
    specialization,
    experience,
    clinic_id,
    bio,
    is_available
) VALUES
      (gen_random_uuid(), 'Pediatric Dentist', 3, '4c074585-4452-48b5-aea1-8e76cd90a880', 'Expert in pediatric dentistry and child care', TRUE),
      (gen_random_uuid(), 'Orthodontist', 5, '005e561f-2b83-4883-9e00-2ff0752f2dc6', 'Specialist in orthodontics and teeth alignment', TRUE),
      (gen_random_uuid(), 'Endodontist', 7, 'dd9e84e8-839a-4797-b9d1-c25f4b5149db', 'Experienced in root canal treatments and endodontic surgery', TRUE),
      (gen_random_uuid(), 'Cosmetic Dentist', 2, '30b673b8-e0d5-43e3-8287-0e9d6b5e3245', 'Focused on cosmetic dentistry and smile makeovers', TRUE),
      (gen_random_uuid(), 'Dental Surgeon', 8, '24e180e2-9c8a-4b03-8f84-c592936de3a4', 'Experienced dental surgeon with focus on implantology', TRUE),
      (gen_random_uuid(), 'Periodontist', 10, '63e7b11f-edf2-4896-9ad2-1d0c53af2317', 'Specialist in gum disease treatment and prevention', FALSE);

-- +goose Down
DELETE FROM doctors
WHERE specialization IN (
                         'Pediatric Dentist',
                         'Orthodontist',
                         'Endodontist',
                         'Cosmetic Dentist',
                         'Dental Surgeon',
                         'Periodontist'
    );