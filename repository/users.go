package repository

import (
	"errors"
	"pvg/model"

	"gorm.io/gorm"
)

func GetSpecificUser(userId int32) (result *model.User, err error) {
	result = &model.User{}
	if err = DB.First(result, userId).Error; err != nil {
		err = ErrNotFound
		return result, err
	}

	return result, nil
}

func GetAllUsers() (results []*model.User, err error) {
	resultOrm := DB.Model(&model.User{})
	if err = resultOrm.Find(&results).Error; err != nil {
		err = ErrNotFound
		return nil, err
	}
	return results, nil
}

func CreateUser(user *model.User) (result *model.User, err error) {
	db := DB.Save(user)
	if err = db.Error; err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(userId int32, updated *model.User) (result *model.User, err error) {
	result = &model.User{}
	db := DB.First(result, userId)
	if err = db.Error; err != nil {
		return nil, ErrNotFound
	}

	if err = Copy(result, updated); err != nil {
		return nil, ErrUpdateFailed
	}

	db = db.Save(result)
	if err = db.Error; err != nil {
		return nil, ErrUpdateFailed
	}

	return result, nil
}

func DeleteUser(userId int32) (err error) {
	record := &model.User{}
	db := DB.First(record, userId)
	if db.Error != nil {
		return ErrNotFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return ErrDeleteFailed
	}

	return nil
}

func GetUserByUsernameOrEmailOrPhoneNumber(username, email, phoneNumber string) (result *model.User, err error) {
	result = &model.User{}
	if err = DB.Where("username = ? OR email = ? OR phone_number = ?", username, email, phoneNumber).First(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func GetUserByEmail(email string) (result *model.User, err error) {
	result = &model.User{}
	if err = DB.Where("email = ?", email).First(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
