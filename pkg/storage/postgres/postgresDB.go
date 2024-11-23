package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Store struct {
	db *pgxpool.Pool
}

func New(s string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), s)
	if err != nil {
		return nil, err
	}
	return &Store{db}, nil
}

// Получение всех публикаций
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT id, title, context,created_at FROM posts`)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения записей: %w", err)
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var post storage.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("Ошибка при сканировании: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

// Создание новой публикации
func (s *Store) AddPost(post storage.Post) error {
	createdAt := time.Now().Unix()
	err := s.db.QueryRow(context.Background(), `
	INSERT INTO posts (author_id,title,context,created_at)
	VALUES ($1,$2,$3,$4)`,
		post.AuthorID,
		post.Title,
		post.Content,
		createdAt)
	if err != nil {
		return fmt.Errorf("Ошибка вставки новой публикации: %v", err)
	}
	return nil
}

// Обновление публикации
func (s *Store) UpdatePost(post storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
	UPDATE posts
	SET title=$1, context =$2
	WHERE author_id=$3`,
		post.Title,
		post.Content,
		post.AuthorID)
	if err != nil {
		return fmt.Errorf("Ошибка при обновлении записи по ID %d : %v", post.AuthorID, err)
	}
	return nil
}

// Удаление публикации по ID
func (s *Store) DeletePost(post storage.Post) error {
	result, err := s.db.Exec(context.Background(), `
	DELETE FROM posts WHERE id=$1`,
		post.ID)
	if err != nil {
		return fmt.Errorf("Ошибка удаления записи по ID %d : %v", post.ID, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Не найдено никаких записей по id %d : %v", post.ID, err)
	}
	return nil
}
