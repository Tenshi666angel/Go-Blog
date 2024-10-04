package sqlite

import (
	"blog/internal/types"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) SavePost(post types.PostEntity) (int64, error) {
	const op = "persistence.sqlite.SavePost"

	stmt, err := s.db.Prepare("INSERT INTO posts(title, content, user_id, app_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(post.Title, post.Content, post.UserId, post.AppID)
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

	rows, err := s.db.Query("SELECT id, title, content, likes_count, user_id, app_id FROM posts")
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
			&post.LikesCount,
			&post.UserId,
			&post.AppID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, post)
	}

	defer rows.Close()

	return &posts, nil
}

func (s *Storage) UpdateLikes(appId string, like int, username string) (bool, error) {
	const op = "persistence.sqlite.UpdateLikes"

	stmt, err := s.db.Prepare("UPDATE posts SET likes_count = likes_count + ? WHERE app_id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(like, appId); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	defer stmt.Close()

	return true, nil
}

func (s *Storage) GetLike(appId string, username string) (*types.Like, error) {
	const op = "persistence.sqlite.GetLike"

	var like types.Like

	stmt, err := s.db.Prepare("SELECT post_id, username FROM likes WHERE post_id = ? AND username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := stmt.QueryRow(appId, username).Scan(&like.PostID, &like.Username); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &like, nil
}

func (s *Storage) CreateLike(appId string, username string) (int64, error) {
	const op = "persistence.sqlite.CreateLike"

	stmt, err := s.db.Prepare("INSERT INTO likes(post_id, username) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(appId, username)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.LastInsertId()
}

func (s *Storage) DeleteLike(appId string, username string) error {
	const op = "persistence.sqlite.DeleteLike"

	stmt, err := s.db.Prepare("DELETE FROM likes WHERE post_id = ? AND username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(appId, username); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
