package models

import "gorm.io/gorm"

type Auth struct {
	UserId    uint64 `json:"user_id" gorm:"not null;unique;comment:id"`
	Token     string `json:"token" gorm:"not null;unique;comment:token"`
	Account   string `json:"account" gorm:"size:60;comment:账号"`
	Name      string `json:"name" gorm:"notnull;comment:用户名"`
	Mail      string `json:"mail" gorm:"notnull;comment:邮箱"`
	Group     string `json:"group" gorm:"notnull;comment:用户分组(管理员;用户)"`
	LoginTime uint64 `json:"login_time" gorm:"autoCreateTime;notnull;comment:登录时间"`
}

func (auth Auth) Edit() (rowsAffected int64, err error) {
	result := db.Model(&Auth{}).Where("token = ?", auth.Token).Updates(auth)
	return result.RowsAffected, result.Error
}

func (auth Auth) Delete() (rowsAffected int64) {
	return db.Delete(&Auth{}, auth.UserId).RowsAffected
}

// 获取详情
func GetAuthDetailsByAny[T getDetailsGenerics](key string, value T) (user User, err error) {
	err = db.Where(key+" = ?", value).Find(&user).Error
	return user, err
}

func (auth Auth) GetDetailsByToken() (a Auth, err error) {
	err = db.Where("token = ?", auth.Token).Find(&auth).Error
	return auth, err
}

// 事务操作 Tx

// 登录
func (auth Auth) TxAdd(tx *gorm.DB) (authId uint64, rowsAffected int64, err error) {
	result := tx.Create(&auth)
	return auth.UserId, result.RowsAffected, result.Error
}

func (auth Auth) TxDelete(tx *gorm.DB) int64 {
	return tx.Delete(&Auth{}, auth.UserId).RowsAffected
}

// 查询数量
func (auth Auth) TxCountById(tx *gorm.DB) (count int64) {
	tx.Model(&Auth{}).Where("user_id = ?", auth.UserId).Count(&count)
	return count
}

func (Auth) TxGetDetailsByToken(tx *gorm.DB, token string) (auth Auth, err error) {
	err = tx.Where("token = ?", token).Find(&auth).Error
	return auth, err
}
