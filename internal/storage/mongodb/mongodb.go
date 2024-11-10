package mongodb

import (
	"log/slog"
	"sync"
	"time"

	"github.com/dngrudin/notes-service/internal/config"
	"github.com/dngrudin/notes-service/internal/model"
	"github.com/dngrudin/notes-service/internal/storage"
	"github.com/dngrudin/notes-service/pkg/logger"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
)

type mongoNoteFields struct {
	ID    uuid.UUID `bson:"_id"`
	Title string    `bson:"title"`
	Text  string    `bson:"text"`
}

type MongoDBStorage struct {
	client *mongo.Client
	db     *mongo.Database
	mu     sync.RWMutex
}

func New(config config.StorageConfig) (storage.Storage, error) {
	log := slog.Default().With(slog.String("storage", "New"))

	clientOptions := options.Client().ApplyURI(config.URI)

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("MongoDB connection error", logger.Err(err))
		return nil, storage.ErrStorage
	}

	db := client.Database(config.DbName)

	return &MongoDBStorage{client: client, db: db}, nil
}

func (s *MongoDBStorage) GetNotes() ([]model.Note, error) {
	log := slog.Default().With(slog.String("storage", "GetNotes"))

	s.mu.RLock()
	defer s.mu.RUnlock()

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.getCollection()
	cur, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Error("Find notes error", logger.Err(err))
		return nil, storage.ErrStorage
	}

	defer cur.Close(ctx)

	var notes []model.Note

	for cur.Next(ctx) {
		var elem mongoNoteFields
		err := cur.Decode(&elem)
		if err != nil {
			log.Error("Decode note error", logger.Err(err))
			return nil, storage.ErrStorage
		}

		notes = append(notes, model.NewNote(elem.ID, elem.Title, elem.Text))
	}

	if err := cur.Err(); err != nil {
		log.Error("Decode note error", logger.Err(err))
		return nil, storage.ErrStorage
	}

	return notes, nil
}

func (s *MongoDBStorage) GetNoteByID(id uuid.UUID) (model.Note, error) {
	log := slog.Default().With(slog.String("storage", "GetNoteByID"))

	s.mu.RLock()
	defer s.mu.RUnlock()

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.getCollection()

	var elem mongoNoteFields
	err := collection.FindOne(ctx, bson.M{"_id": id}).
		Decode(&elem)

	if err != nil {
		log.Error("Find note error", logger.Err(err))
		return model.Note{}, storage.ErrStorage
	}

	return model.NewNote(elem.ID, elem.Title, elem.Text), nil
}

func (s *MongoDBStorage) CreateNote(note model.Note) error {
	log := slog.Default().With(slog.String("storage", "CreateNote"))

	s.mu.Lock()
	defer s.mu.Unlock()

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.getCollection()
	_, err := collection.InsertOne(ctx, mongoNoteFields{note.ID, note.Title, note.Text})

	if err != nil {
		log.Error("Insert note error", logger.Err(err))
		return storage.ErrStorage
	}

	return nil
}

func (s *MongoDBStorage) UpdateNote(note model.Note) error {
	log := slog.Default().With(slog.String("storage", "UpdateNote"))

	s.mu.Lock()
	defer s.mu.Unlock()

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.getCollection()
	updateResult, err := collection.UpdateOne(ctx, bson.M{"_id": note.ID}, bson.M{"$set": bson.M{"title": note.Title, "text": note.Text}})

	if err != nil {
		log.Error("Update note error", logger.Err(err))
		return storage.ErrStorage
	}

	if updateResult.MatchedCount > 0 {
		return nil
	}

	return storage.ErrNotFound
}

func (s *MongoDBStorage) DeleteNoteByID(id uuid.UUID) error {
	log := slog.Default().With(slog.String("storage", "DeleteNoteByID"))

	s.mu.Lock()
	defer s.mu.Unlock()

	var ctx, cancel = context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := s.getCollection()
	deleteResult, err := collection.DeleteMany(ctx, bson.M{"_id": id})
	if err != nil {
		log.Error("Delete note error", logger.Err(err))
		return storage.ErrStorage
	}

	if deleteResult.DeletedCount > 0 {
		return nil
	}

	return storage.ErrNotFound
}

func (s *MongoDBStorage) getCollection() *mongo.Collection {
	return s.db.Collection("notes")
}
