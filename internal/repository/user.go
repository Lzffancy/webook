package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var ErrDuplicateEmail = dao.ErrDuplicateEmail
var ErrUserNotFind = dao.ErrRecordNotFind

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	//S

	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	//S
	u, err := repo.dao.FindByEmail(ctx, email)
	switch err {
	case nil:
		return repo.toDomain(u), nil
	case ErrUserNotFind:
		return domain.User{}, ErrUserNotFind
	default:
		return domain.User{}, err
	}
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repo *UserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *UserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), err
}

// 为什么？从repo 里面返回Entity user 而不是domain user
// 为什么叫UpdateNonZeroFields？
func (repo *UserRepository) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}
