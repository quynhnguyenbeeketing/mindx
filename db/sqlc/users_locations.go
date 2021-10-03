package db

import (
	"context"
	"time"
)

const users_locations = `CREATE TABLE users_locations
(
    id          BIGINT PRIMARY KEY,
    user_id     BIGINT REFERENCES users (id),
    location_id BIGINT REFERENCES locations (id),
    event_time  TIMESTAMP
);`

type UsersLocations struct {
	Id           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	LocationId   int64     `json:"location_id"`
	HealthStatus int       `json:"health_status"`
	EventTime    time.Time `json:"event_time"`
}

func (q *Queries) CreateUsersLocations(ctx context.Context, req *UsersLocations) error {
	query := `INSERT INTO users_locations(user_id, location_id, health_status, event_time)
					VALUES ($1, $2, $3, NOW());`

	q.db.QueryRowContext(ctx, query,
		req.UserId, req.LocationId, req.HealthStatus)
	return nil
}

func (q *Queries) GetUsersLocations(ctx context.Context, req *UsersLocations) ([]*UsersLocations, error) {
	query := `SELECT * FROM users_locations WHERE user_id = $1;`

	rows, err := q.db.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*UsersLocations, 0)
	for rows.Next() {
		model := &UsersLocations{}

		if err = rows.Scan(
			&req.Id,
			&req.UserId,
			&req.LocationId,
			&req.EventTime,
		); err != nil {
			return nil, err
		}

		res = append(res, model)
	}

	return res, err
}
