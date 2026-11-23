package repository

type UserRepository struct {

}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByID(id string) (interface{}, error) {
	return nil, nil
}
