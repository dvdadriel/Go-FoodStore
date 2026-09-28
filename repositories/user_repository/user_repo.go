package userrepository

import "go-food-store/models"

// UserRepo mengembalikan (data, error), bukan response.WebResponse seperti
// repository lama di repo ini. Layer data tidak seharusnya tahu status code
// HTTP; ini bentuk yang benar, dan repository lain menyusul kalau disentuh.
type UserRepo interface {
	Create(user models.User) (models.User, error)
	FindByUsername(username string) (models.User, error)
	ExistsByUsername(username string) (bool, error)
}
