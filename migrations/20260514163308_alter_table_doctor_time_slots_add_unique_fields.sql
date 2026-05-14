-- +goose Up

-- удаляем дубликаты
DELETE FROM doctor_time_slots
WHERE id IN (
    SELECT id
    FROM (
             SELECT id,
                    ROW_NUMBER() OVER (
                   PARTITION BY
                       doctor_id,
                       clinic_address_id,
                       slot_start,
                       slot_end
                   ORDER BY id
               ) AS rn
             FROM doctor_time_slots
         ) t
    WHERE t.rn > 1
);

-- добавляем unique constraint
ALTER TABLE doctor_time_slots
    ADD CONSTRAINT unique_slot
        UNIQUE (
                doctor_id,
                clinic_address_id,
                slot_start,
                slot_end
            );

-- +goose Down

ALTER TABLE doctor_time_slots
DROP CONSTRAINT unique_slot;