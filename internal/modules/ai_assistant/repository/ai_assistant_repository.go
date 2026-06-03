package repository

import (
	"context"

	"dental_clinic/internal/modules/ai_assistant/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AIAssistantRepository interface {
	CreateSession(userID uuid.UUID) (*models.ChatSession, error)
	GetOrCreateSession(userID uuid.UUID) (*models.ChatSession, error)
	SaveMessage(sessionID uuid.UUID, role, content string) error
	GetRecentMessages(sessionID uuid.UUID, limit int) ([]models.ChatMessage, error)
	GetOrCreateState(userID uuid.UUID) (*models.BookingState, error)
	SaveState(state *models.BookingState) error
	ClearState(userID uuid.UUID) error
	SearchServices(query string) ([]models.ServiceOption, error)
	GetClinicOptions(serviceID string) ([]models.ClinicOption, error)
	GetDoctorOptions(serviceID, clinicAddressID string) ([]models.DoctorOption, error)
}

type aiAssistantRepo struct {
	db *pgxpool.Pool
}

func NewAIAssistantRepository(db *pgxpool.Pool) AIAssistantRepository {
	return &aiAssistantRepo{db: db}
}

func (r *aiAssistantRepo) CreateSession(userID uuid.UUID) (*models.ChatSession, error) {
	session := &models.ChatSession{
		Id:     uuid.New(),
		UserID: userID,
	}
	insertQuery := `INSERT INTO chat_sessions (id, user_id) VALUES ($1, $2)`
	_, err := r.db.Exec(context.Background(), insertQuery, session.Id, session.UserID)
	return session, err
}

func (r *aiAssistantRepo) GetOrCreateSession(userID uuid.UUID) (*models.ChatSession, error) {
	session := &models.ChatSession{}
	query := `
		SELECT id, user_id
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT 1
	`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&session.Id, &session.UserID)
	if err == nil {
		return session, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	return r.CreateSession(userID)
}

func (r *aiAssistantRepo) SaveMessage(sessionID uuid.UUID, role, content string) error {
	query := `INSERT INTO chat_messages (id, session_id, role, content) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(context.Background(), query, uuid.New(), sessionID, role, content)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(context.Background(), `UPDATE chat_sessions SET updated_at = NOW() WHERE id = $1`, sessionID)
	return err
}

func (r *aiAssistantRepo) GetRecentMessages(sessionID uuid.UUID, limit int) ([]models.ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(context.Background(), query, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]models.ChatMessage, 0)
	for rows.Next() {
		var message models.ChatMessage
		if err := rows.Scan(&message.Id, &message.SessionID, &message.Role, &message.Content); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *aiAssistantRepo) GetOrCreateState(userID uuid.UUID) (*models.BookingState, error) {
	state := &models.BookingState{}
	query := `
		SELECT
			user_id::text,
			COALESCE(doctor_id::text, ''),
			COALESCE(service_id::text, ''),
			COALESCE(clinic_address_id::text, ''),
			COALESCE(date::text, ''),
			COALESCE(time::text, ''),
			COALESCE(step, '')
		FROM ai_booking_state
		WHERE user_id = $1
	`
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
		&state.UserID,
		&state.DoctorID,
		&state.ServiceID,
		&state.ClinicAddressID,
		&state.Date,
		&state.Time,
		&state.Step,
	)
	if err == nil {
		return state, nil
	}
	if err != pgx.ErrNoRows {
		return nil, err
	}

	state.UserID = userID.String()
	state.Step = "collect_service"
	err = r.SaveState(state)
	return state, err
}

func (r *aiAssistantRepo) SaveState(state *models.BookingState) error {
	query := `
		INSERT INTO ai_booking_state (user_id, doctor_id, service_id, clinic_address_id, date, time, step, updated_at)
		VALUES ($1::uuid, NULLIF($2, '')::uuid, NULLIF($3, '')::uuid, NULLIF($4, '')::uuid, NULLIF($5, '')::date, NULLIF($6, '')::time, $7, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			doctor_id = EXCLUDED.doctor_id,
			service_id = EXCLUDED.service_id,
			clinic_address_id = EXCLUDED.clinic_address_id,
			date = EXCLUDED.date,
			time = EXCLUDED.time,
			step = EXCLUDED.step,
			updated_at = NOW()
	`
	_, err := r.db.Exec(
		context.Background(),
		query,
		state.UserID,
		state.DoctorID,
		state.ServiceID,
		state.ClinicAddressID,
		state.Date,
		state.Time,
		state.Step,
	)
	return err
}

func (r *aiAssistantRepo) ClearState(userID uuid.UUID) error {
	_, err := r.db.Exec(context.Background(), `DELETE FROM ai_booking_state WHERE user_id = $1`, userID)
	return err
}

func (r *aiAssistantRepo) SearchServices(query string) ([]models.ServiceOption, error) {
	sql := `
		SELECT id::text, name, COALESCE(name_en, ''), COALESCE(name_kaz, ''), COALESCE(description, '')
		FROM services
		WHERE LOWER(name) LIKE LOWER('%' || $1 || '%')
		   OR LOWER(COALESCE(name_en, '')) LIKE LOWER('%' || $1 || '%')
		   OR LOWER(COALESCE(name_kaz, '')) LIKE LOWER('%' || $1 || '%')
		ORDER BY name
		LIMIT 5
	`
	rows, err := r.db.Query(context.Background(), sql, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.ServiceOption, 0)
	for rows.Next() {
		var option models.ServiceOption
		if err := rows.Scan(&option.Id, &option.Name, &option.NameEn, &option.NameKaz, &option.Description); err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, rows.Err()
}

func (r *aiAssistantRepo) GetClinicOptions(serviceID string) ([]models.ClinicOption, error) {
	query := `
		SELECT
			cs.clinic_id::text,
			ca.id::text,
			c.name,
			cs.price,
			cs.duration_minutes,
			COALESCE(ROUND(AVG(cr.rating)::numeric, 2), 0)::float8 AS rating
		FROM clinic_services cs
		JOIN clinics c ON c.id = cs.clinic_id
		JOIN clinic_addresses ca ON ca.clinic_id = cs.clinic_id
		LEFT JOIN clinic_reviews cr ON cr.clinic_id = c.id
		WHERE cs.service_id = $1
			AND cs.is_active = true
			AND c.is_active = true
		GROUP BY cs.clinic_id, ca.id, c.name, cs.price, cs.duration_minutes
		ORDER BY rating DESC, c.name
	`
	rows, err := r.db.Query(context.Background(), query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.ClinicOption, 0)
	for rows.Next() {
		var option models.ClinicOption
		if err := rows.Scan(&option.ClinicID, &option.ClinicAddressID, &option.ClinicName, &option.Price, &option.Duration, &option.Rating); err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, rows.Err()
}

func (r *aiAssistantRepo) GetDoctorOptions(serviceID, clinicAddressID string) ([]models.DoctorOption, error) {
	query := `
		SELECT
			d.id::text,
			d.name,
			d.specialization,
			d.experience,
			COALESCE(ROUND(AVG(dr.rating)::numeric, 2), 0)::float8 AS rating
		FROM doctors d
		JOIN clinic_addresses ca ON ca.clinic_id = d.clinic_id
		JOIN clinic_services cs ON cs.clinic_id = d.clinic_id
		LEFT JOIN doctor_ratings dr ON dr.doctor_id = d.id
		WHERE ca.id = $1
			AND cs.service_id = $2
			AND d.is_available = true
			AND d.is_deleted = 0
		GROUP BY d.id, d.name, d.specialization, d.experience
		ORDER BY rating DESC, d.name
	`
	rows, err := r.db.Query(context.Background(), query, clinicAddressID, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	options := make([]models.DoctorOption, 0)
	for rows.Next() {
		var option models.DoctorOption
		if err := rows.Scan(&option.Id, &option.Name, &option.Specialization, &option.Experience, &option.Rating); err != nil {
			return nil, err
		}
		options = append(options, option)
	}
	return options, rows.Err()
}
