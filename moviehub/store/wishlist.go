package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

type WishlistItem struct {
	Movie
	Note    *string   `json:"note"`
	AddedAt time.Time `json:"addedAt"`
}

func (s *Store) GetWishlist(ctx context.Context, memberID int64) ([]WishlistItem, error) {
	rows, err := s.Pool.Query(ctx,
		`SELECT m.id, m.title, m.title_th, coalesce(g.name_th, g.name_en),
                extract(year FROM m.release_date)::int,
                round(m.vote_average, 1)::float8,
                'https://image.tmdb.org/t/p/w500' || m.poster_path,
                w.note, w.added_at
         FROM wishlist_items w
         JOIN movies m ON m.id = w.movie_id
         LEFT JOIN movie_genres mg ON mg.movie_id = m.id AND mg.position = 0
         LEFT JOIN genres g ON g.id = mg.genre_id
         WHERE w.member_id = $1
         ORDER BY w.added_at DESC`, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []WishlistItem{}
	for rows.Next() {
		var item WishlistItem
		if err := rows.Scan(&item.ID, &item.Title, &item.TitleTh, &item.Genre,
			&item.Year, &item.Rating, &item.Poster, &item.Note, &item.AddedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) PutWishlist(ctx context.Context, memberID, movieID int64, note *string) error {
	_, err := s.Pool.Exec(ctx,
		`INSERT INTO wishlist_items (member_id, movie_id, note) VALUES ($1, $2, $3)
         ON CONFLICT (member_id, movie_id) DO UPDATE SET note = EXCLUDED.note`,
		memberID, movieID, note)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return ErrNoMovie
	}
	return err
}

func (s *Store) DeleteWishlist(ctx context.Context, memberID, movieID int64) error {
	_, err := s.Pool.Exec(ctx,
		`DELETE FROM wishlist_items WHERE member_id = $1 AND movie_id = $2`, memberID, movieID)
	return err
}
