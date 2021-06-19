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
		client *firestore.Client
	}

	firestoreDocument struct {
		Target  string `firestore:"target"`
		Address string `firestore:"address"`
	}
)

var (
	project    = os.Getenv("GCP_PROJECT")
	collection = os.Getenv("GCP_FIRESTORE_COLLECTION")
)

func IsFirestoreStorageAvailable() error {
	if project == "" {
		return fmt.Errorf("GCP_PROJECT is missing")
	}
	if collection == "" {
		return fmt.Errorf("GCP_FIRESTORE_COLLECTION is missing")
	}
	return nil
}

func NewFirestoreStorage() (Storage, error) {
	client, err := firestore.NewClient(context.Background(), project)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}
	return &FirestoreStorage{client}, nil
}

func (s *FirestoreStorage) ref(key string) *firestore.DocumentRef {
	return s.client.Doc(fmt.Sprintf("%s/%s", collection, key))
}

func (s *FirestoreStorage) Get(key string) (string, error) {
	snapshot, err := s.ref(key).Get(context.Background())
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

func (s *FirestoreStorage) Set(key string, value string) error {
	if _, err := s.ref(key).Set(context.Background(), &firestoreDocument{key, value}); err != nil {
		return fmt.Errorf("failed to write document: %w", err)
	}
	return nil
}

func (s *FirestoreStorage) Unset(key string) error {
	if _, err := s.ref(key).Delete(context.Background()); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
