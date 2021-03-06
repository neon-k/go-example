package persistence

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/api/domain/model"
	"github.com/api/domain/repository"
)

type UserPersistence struct{}

func NewUserPersistence() repository.UserRepository {
	return &UserPersistence{}
}

// ユーザー情報全て取得
func (up UserPersistence) GetAll(DB *gorm.DB) ([]*model.User, error) {

	var users []*model.User
	user := &model.User{}

	// ユーザー全て取得
	err := DB.Scopes(user.GetResParam).Preload("Posts").Preload("Goods").Find(&users).Error

	return users, err
}

// ログインユーザーを取得
func (up UserPersistence) GetCurrentUser(DB *gorm.DB, email string) (*model.User, error) {

	currentUser := &model.User{}

	// メールでユーザーを絞り込む
	err := DB.
		Preload("Posts").
		Preload("Goods").
		Where("email = ?", email).
		First(&currentUser).Error

	return currentUser, err
}

// IDでユーザー情報を取得
func (up UserPersistence) GetCurrentUserID(DB *gorm.DB, ID float64) (*model.User, error) {

	user := &model.User{}
	var posts []*model.PostRes
	post := model.Post{}

	// メールでユーザーを絞り込む
	if err := DB.Scopes(user.GetResParam).Preload("Goods").Where("id = ?", ID).First(&user).Error; err != nil {
		return user, err
	}

	err := DB.
		Table(post.TableName()).
		Where("posts.user_id = ?", ID).
		Scopes(post.GetGoodCount).
		Find(&posts).
		Error

	user.Posts = posts

	return user, err
}

// ユーザー情報登録
func (up UserPersistence) AddUser(DB *gorm.DB, name string, age int, icon string, password string, email string) error {

	user := model.User{
		Name:     name,
		Age:      age,
		Icon:     icon,
		Password: password,
		Email:    email,
	}

	// ユーザー情報を登録
	err := DB.Create(&user).Error

	return err
}

// ユーザー情報を更新
func (up UserPersistence) UpdateUser(DB *gorm.DB, data *model.User) error {
	err := DB.Save(&data).Error

	return err
}

// ユーザー情報を削除
func (up UserPersistence) DeleteUser(DB *gorm.DB, id int) error {
	result := DB.Where("id = ?", id).Delete(&model.User{}).RowsAffected

	var err error

	if result != 0 {
		err = nil
	} else {
		err = errors.New("delete error")
	}

	return err
}
