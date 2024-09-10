package utils

import "gorm.io/gorm"

func HandleGormError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}