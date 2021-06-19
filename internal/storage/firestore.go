package storage

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	FirestoreStorage struct {
		collection string
		client     *firestore.Client
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
		collection: collection,
		client:     client,
	}, nil
}

func (s *FirestoreStorage) ref(key string) *firestore.DocumentRef {
	return s.client.Doc(fmt.Sprintf("%s/%s", s.collection, key))
}

func (s *FirestoreStorage) Get(ctx context.Context, key string) (string, error) {
	snapshot, err := s.ref(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "", fmt.Errorf("%w: %v", ErrorNotFound, key)
		}
		return "", fmt.Errorf("failed to get document: %w", err)
	}

	data := &firestoreDocument{}
	if err := snapshot.DataTo(&data); err != nil {
		return "", fmt.Errorf("failed to read document: %w", err)
	}
	return data.Address, nil
}

func (s *FirestoreStorage) Set(ctx context.Context, key, value string) error {
	if _, err := s.ref(key).Set(ctx, &firestoreDocument{key, value}); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

func (s *FirestoreStorage) Unset(ctx context.Context, key string) error {
	if _, err := s.ref(key).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
