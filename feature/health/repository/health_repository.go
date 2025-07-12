package repository

type HealthRepository interface {
	CheckSystemHealth() error
}

type healthRepository struct{}

func NewHealthRepository() HealthRepository {
	return &healthRepository{}
}

func (r *healthRepository) CheckSystemHealth() error {
	return nil
}
