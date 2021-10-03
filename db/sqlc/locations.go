package db

import (
	"context"
	"time"
)

const locations = `CREATE TABLE locations
(
    id         BIGINT PRIMARY KEY,
    address    VARCHAR(255),
    city_id    BIGINT,
    country_id BIGINT,

    created_at TIMESTAMP,
    updated_at TIMESTAMP
);`

type Locations struct {
	Id        int64     `json:"id"`
	Address   string    `json:"address"`
	CityId    int64     `json:"city_id"`
	CountryId int64     `json:"country_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) CreateLocations(ctx context.Context, req *Locations) (*Locations, error) {
	query := `INSERT INTO locations(address, city_id, country_id)
					VALUES ($1, $2, $3);`

	row := q.db.QueryRowContext(ctx, query,
		req.Address,
		req.CityId,
		req.CountryId)

	err := row.Scan()
	return req, err
}

func (q *Queries) GetLocations(ctx context.Context, req *Locations) ([]*Locations, error) {
	query := `SELECT * FROM locations;`

	rows, err := q.db.QueryContext(ctx, query, req.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*Locations, 0)
	for rows.Next() {
		model := &Locations{}

		if err = rows.Scan(
			&req.Id,
			&req.Address,
			&req.CityId,
			&req.CountryId,
			&req.CreatedAt,
			&req.UpdatedAt,
		); err != nil {
			return nil, err
		}

		res = append(res, model)
	}

	return res, err
}
