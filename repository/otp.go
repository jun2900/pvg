package repository

import (
	"errors"
	"pvg/model"

	"gorm.io/gorm"
)

func AddOtp(otp *model.Otp) (result *model.Otp, err error) {
	db := DB.Save(otp)
	if err = db.Error; err != nil {
		return nil, ErrInsertFailed
	}

	return otp, nil
}
func GetOtpByEmail(email string) (result *model.Otp, err error) {
	result = &model.Otp{}
	if err = DB.Where("email = ?", email).First(result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
func UpdateOtp(otpId int32, updated *model.Otp) (result *model.Otp, err error) {
	result = &model.Otp{}
	db := DB.First(result, otpId)
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
func DeleteOtp(otpId int32) (err error) {
	record := &model.Otp{}
	db := DB.First(record, otpId)
	if db.Error != nil {
		return ErrNotFound
	}

	db = db.Delete(record)
	if err = db.Error; err != nil {
		return ErrDeleteFailed
	}

	return nil
}
