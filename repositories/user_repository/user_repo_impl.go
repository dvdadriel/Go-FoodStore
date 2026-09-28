package userrepository

import (
	"gorm.io/gorm"

	"go-food-store/models"
)

type UserRepoImpl struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) UserRepo {
	return &UserRepoImpl{DB: DB}
}

func (u *UserRepoImpl) Create(user models.User) (models.User, error) {
	err := u.DB.Create(&user).Error
	return user, err
}

// FindByUsername mengembalikan gorm.ErrRecordNotFound kalau tidak ada.
func (u *UserRepoImpl) FindByUsername(username string) (models.User, error) {
	var user models.User
	err := u.DB.Where("username = ?", username).First(&user).Error
	return user, err
}

func (u *UserRepoImpl) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := u.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
