package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Store struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func New(s string, dbName string, collectionName string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(s)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	collection := client.Database(dbName).Collection(collectionName)

	return &Store{client, collection}, nil
}

// Получение всех публикаций
func (s *Store) Posts() ([]storage.Post, error) {
	filter := bson.D{{}}
	cursor, err := s.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения постов: %w", err)
	}
	defer cursor.Close(context.Background())
	var posts []storage.Post
	for cursor.Next(context.Background()) {
		var post storage.Post
		if err := cursor.Decode(&post); err != nil {
			return nil, fmt.Errorf("ошибка декодирования поста: %w", err)
		}
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по курсору: %w", err)
	}
	return posts, nil
}

// Создание новой публикации
func (s *Store) AddPost(post storage.Post) error {
	doc := bson.D{
		{"author_id", post.AuthorID},
		{"title", post.Title},
		{"content", post.Content},
		{"created_at", time.Now().Unix()},
	}

	_, err := s.collection.InsertOne(context.Background(), doc)
	if err != nil {
		return fmt.Errorf("ошибка вставки новой публикации: %w", err)
	}
	return nil
}

// Обновление публикации
func (s *Store) UpdatePost(post storage.Post) error {
	filter := bson.D{{"_id", post.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"title", post.Title},
			{"content", post.Content},
		}},
	}
	result, err := s.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении публикации: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("не найдено публикаций с ID %d", post.ID)
	}
	return nil
}

// Удаление публикации по ID
func (s *Store) DeletePost(post storage.Post) error {
	filter := bson.D{{"_id", post.ID}}
	result, err := s.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("ошибка при удалении публикации: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("не найдено публикаций с ID %d", post.ID)
	}
	return nil
}
