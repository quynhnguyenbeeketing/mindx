package db

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const users = `CREATE TABLE users
(
    id                BIGINT PRIMARY KEY,
    name              VARCHAR(255),
    permanent_address VARCHAR(255),  -- full address in text
    current_address   VARCHAR(255),  -- location in detail
    location_id       BIGINT REFERENCES locations (id),
    health_status     INT DEFAULT 0, -- 1 - 3 means F0 - F2, other means unknown

    created_at        TIMESTAMP,
    updated_at        TIMESTAMP
);`

type Users struct {
	Id               int64     `json:"id"`
	Name             string    `json:"name"`
	PermanentAddress string    `json:"permanent_address"`
	CurrentAddress   string    `json:"current_address"`
	LocationId       int64     `json:"location_id"`
	HealthStatus     int64     `json:"health_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (q *Queries) GetUsers(ctx context.Context, req *Users) ([]*Users, error) {
	query := `SELECT * FROM users ORDER BY id LIMIT 100;`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*Users, 0)
	for rows.Next() {
		model := &Users{}

		if err = rows.Scan(
			&model.Id,
			&model.Name,
			&model.PermanentAddress,
			&model.CurrentAddress,
			&model.LocationId,
			&model.HealthStatus,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, err
		}

		res = append(res, model)
	}

	return res, err
}

func (q *Queries) CreateUsers(ctx context.Context, req *Users) (*Users, error) {
	query := `INSERT INTO users(name, permanent_address, current_address, location_id, health_status, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
					RETURNING id;`

	row := q.db.QueryRowContext(ctx, query,
		req.Name,
		req.PermanentAddress,
		req.CurrentAddress,
		req.LocationId,
		req.HealthStatus)

	err := row.Scan(&req.Id)

	query = `INSERT INTO users_locations(user_id, location_id, health_status, event_time)
					VALUES ($1, $2, $3, NOW());`

	q.db.QueryRowContext(ctx, query,
		req.Id, req.LocationId, req.HealthStatus)
	return req, err
}

func (q *Queries) CreateListUsers(ctx context.Context, req []*Users) error {
	params := ""
	args := make([]interface{}, 0)

	for i, user := range req {
		sub := ""
		for j := 1; j <= 5; j++ {
			sub += fmt.Sprintf(`$%v, `, i*5+j)
		}

		sub += "NOW(), NOW()"
		params += fmt.Sprintf("(%v), ", sub)

		args = append(args, user.Name, user.PermanentAddress,
			user.CurrentAddress, user.LocationId, user.HealthStatus)
	}

	query := fmt.Sprintf(`INSERT INTO users(name, permanent_address, current_address, location_id, health_status,
					created_at, updated_at) VALUES %v`, strings.TrimSuffix(params, ", "))
	q.db.QueryRowContext(ctx, query, args...)
	return nil
}

func (q *Queries) UpdateUsers(ctx context.Context, req *Users) error {
	query := `UPDATE users
				SET name              = $1,
    				permanent_address = $2,
    				current_address   = $3,
    				location_id       = $4,
    				health_status     = $5,
					updated_at		  = NOW()
				WHERE id = $6;`

	q.db.QueryRowContext(ctx, query,
		req.Name,
		req.PermanentAddress,
		req.CurrentAddress,
		req.LocationId,
		req.HealthStatus,
		req.Id)

	query = `INSERT INTO users_locations(user_id, location_id, health_status, event_time)
					VALUES ($1, $2, $3, NOW());`

	q.db.QueryRowContext(ctx, query,
		req.Id, req.LocationId, req.HealthStatus)
	return nil
}
