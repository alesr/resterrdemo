package foo

import (
	"fmt"
)

type repository interface{ Fetch() error }

// Service implements the domain layer for handling foo entities.
type Service struct{ repo repository }

// New instantiates a new service struct.
func New(repo repository) *Service { return &Service{repo: repo} }

// Fetch would naturally perform some business logic,
// fetching the foo entity from the repository layer.
func (s *Service) Fetch() error {
	if err := s.repo.Fetch(); err != nil {
		return fmt.Errorf("could not fetch foo from repo: '%s': %w", err, ErrGetFaleid)
	}
	return nil
}
