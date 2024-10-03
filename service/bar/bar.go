package bar

import (
	"errors"
	"fmt"
)

type repository interface {
	Fetch() error
}

// Service implements the domain layer for handling bar entities.
type Service struct {
	repo repository
}

// New instantiates a new service struct.
func New(repo repository) *Service {
	return &Service{repo: repo}
}

// Fetch would naturally perform some business logic,
// fetching the bar entity from the repository layer.
func (s *Service) Fetch() error {
	if err := s.repo.Fetch(); err != nil {
		// In this example, we don't want to return this exact repository error to the transport layer.
		// Instead, we replace it with something that better represents our use case (e.g., unavailability).
		if errors.Is(err, ErrBarNotFound) {
			return fmt.Errorf("could not fetch bar (%w): %w", err, ErrBarUnavailable)
		}

		// If we encounter an unexpected error that we're not prepared to handle,
		// we can add the context we need and safely return it,
		// knowing that it won't be mapped as one of the known errors in the error map.
		return fmt.Errorf("could not fetch bar from repo: %w", err)
	}
	return nil
}
