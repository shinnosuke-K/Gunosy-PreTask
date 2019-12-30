package model

import (
	"crypto/sha256"
	"fmt"

	"github.com/jinzhu/gorm"
)

type AccountInfo struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Comment  string `json:"comment"`
}

func (info *AccountInfo) Insert(db *gorm.DB) (*AccountInfo, string) {

	if info.UserId == "" || info.Password == "" {
		return nil, "required"
	}

	if len(info.UserId) < 6 || len(info.UserId) >= 20 {
		return nil, "length"
	}

	var searchedID AccountInfo
	if db.Where("user_id=?", info.UserId).Find(&searchedID); searchedID.UserId != "" {
		return nil, "duplication"

	} else {
		info.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(info.Password)))
		db.Create(&info)
	}

	return info, ""
}

func (info *AccountInfo) AccountByID(db *gorm.DB) (*AccountInfo, string) {
	var accountInfos AccountInfo
	if db.Where("user_id = ?", info.UserId).Find(&accountInfos); accountInfos.UserId == "" {
		return nil, "No User found"
	}

	return &accountInfos, ""
}

func (info *AccountInfo) Update(db *gorm.DB) (*AccountInfo, string) {
	var updatedAccount AccountInfo
	if db.Where("user_id = ?", info.UserId).Find(&updatedAccount); updatedAccount.UserId == "" {
		return nil, "No User found"
	}

	beforeInfo := AccountInfo{}
	afterInfo := beforeInfo
	db.Where("user_id = ?", info.UserId).Find(&beforeInfo)

	if info.Nickname == "" {
		afterInfo.Nickname = info.UserId
	} else {
		afterInfo.Nickname = info.Nickname
	}

	if info.Comment == "" {
		afterInfo.Comment = " "
	} else {
		afterInfo.Comment = info.Comment
	}
	db.Model(&beforeInfo).Where("user_id=?", info.UserId).Update(&afterInfo)
	return &afterInfo, ""
}

func (info *AccountInfo) Delete(db *gorm.DB) string {
	var deletedAccount AccountInfo
	db.Where("user_id=?", info.UserId).Delete(&deletedAccount)
	return "Account and user successfully removed"
}
