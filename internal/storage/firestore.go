package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	FirestoreStorage struct {
		collection *firestore.CollectionRef
	}

	firestoreDocument struct {
		Target  string `firestore:"target"`
		Address string `firestore:"address"`
	}
)

func NewFirestoreStorage(ctx context.Context) (Storage, error) {
	project := os.Getenv("GCP_PROJECT")
	if project == "" {
		return nil, fmt.Errorf("GCP_PROJECT is missing")
	}

	collection := os.Getenv("GCP_FIRESTORE_COLLECTION")
	if collection == "" {
		return nil, fmt.Errorf("GCP_FIRESTORE_COLLECTION is missing")
	}

	client, err := firestore.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return &FirestoreStorage{
		collection: client.Collection(collection),
	}, nil
}

func (s *FirestoreStorage) findByKey(ctx context.Context, key string) (*firestore.DocumentSnapshot, error) {
	snapshot, err := s.collection.Doc(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("%w: key=%v", ErrorUndefinedKey, key)
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return snapshot, nil
}
func (s *FirestoreStorage) findByValue(ctx context.Context, value string) (*firestore.DocumentSnapshot, error) {
	snapshots, err := s.collection.Where("address", "==", value).Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	if len(snapshots) == 0 {
		return nil, fmt.Errorf("%w: value=%v", ErrorUndefinedValue, value)
	}
	return snapshots[0], nil
}

func (s *FirestoreStorage) Get(ctx context.Context, key string) (string, error) {
	snapshot, err := s.findByKey(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to find document: %w", err)
	}

	data := &firestoreDocument{}
	if err := snapshot.DataTo(&data); err != nil {
		return "", fmt.Errorf("failed to read document: %w", err)
	}
	return data.Address, nil
}

func (s *FirestoreStorage) Set(ctx context.Context, key, value string) error {
	if _, err := s.findByKey(ctx, key); err == nil {
		return fmt.Errorf("%w: key=%v", ErrorDuplicatedKey, key)
	} else if !errors.Is(err, ErrorUndefinedKey) {
		return fmt.Errorf("error occurred while querying by key: %v", err)
	}

	if _, err := s.findByValue(ctx, value); err == nil {
		return fmt.Errorf("%w: value=%v", ErrorDuplicatedValue, value)
	} else if !errors.Is(err, ErrorUndefinedValue) {
		return fmt.Errorf("error occurred while querying by value: %v", err)
	}

	if _, err := s.collection.Doc(key).Create(ctx, &firestoreDocument{key, value}); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

func (s *FirestoreStorage) UnsetByKey(ctx context.Context, key string) error {
	snapshot, err := s.findByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}

	if _, err := snapshot.Ref.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

func (s *FirestoreStorage) UnsetByValue(ctx context.Context, value string) error {
	snapshot, err := s.findByValue(ctx, value)
	if err != nil {
		return fmt.Errorf("failed to find document: %w", err)
	}

	if _, err := snapshot.Ref.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
