package sqlite

import (
	"blog/internal/types"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) SavePost(post types.PostEntity) (int64, error) {
	const op = "persistence.sqlite.SavePost"

	stmt, err := s.db.Prepare("INSERT INTO posts(title, content, user_id) VALUES(?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(post.Title, post.Content, post.User_id)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok &&
			sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	defer stmt.Close()

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetAll() (*[]types.PostEntity, error) {
	const op = "persistence.sqlite.GetAll"

	rows, err := s.db.Query("SELECT id, title, content, likes_count, user_id FROM posts")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var posts []types.PostEntity

	for rows.Next() {
		var post types.PostEntity
		if err := rows.Scan(
			&post.Id,
			&post.Title,
			&post.Content,
			&post.Likes_count,
			&post.User_id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}

	defer rows.Close()

	return &posts, nil
}

func (s *Storage) UpdateLikes(user_id int64, like int) (bool, error) {
	const op = "persistence.sqlite.UpdateLikes"

	stmt, err := s.db.Prepare("UPDATE posts SET likes_count = likes_count + ? WHERE user_id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(like, user_id); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	return true, nil
}
