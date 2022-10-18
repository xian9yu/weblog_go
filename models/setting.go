package models

type Setting struct {
	Id          uint64 `json:"id" gorm:"primaryKey;comment:id"`
	Name        string `json:"name" gorm:"not null;comment:名称"`
	Value       string `json:"value" gorm:"not null;comment:参数"`
	Remark      string `json:"remark" gorm:"comment:备注"`
	CreatedTime uint64 `json:"created_time" gorm:"not null;comment:创建时间"`
	UpdatedTime uint64 `json:"updated_time" gorm:"not null;comment:编辑时间"`
}

func (setting Setting) Add() (id uint64, rowsAffected int64, err error) {
	result := db.Create(&setting)
	return setting.Id, result.RowsAffected, result.Error
}

func (setting Setting) Edit() (rowsAffected int64, err error) {
	result := db.Model(&Setting{}).Where("id = ?", setting.Id).Updates(setting)
	return result.RowsAffected, result.Error
}

func (Setting) Delete(id uint64) int64 {
	return db.Delete(&Setting{}, id).RowsAffected
}

func (Setting) List() (total int64, setting []Setting, err error) {
	db.Model(&Setting{}).Count(&total)
	err = db.Find(&setting).Error
	return total, setting, err
}

func GetSettingDetailsByAny[T getDetailsGenerics](key string, value T) (setting Setting, err error) {
	err = db.Where(key+" = ?", value).Find(&setting).Error
	return setting, err
}

func CountSettingByAny[T getDetailsGenerics](key string, value T) (count int64) {
	db.Model(&Setting{}).Where(key+" = ?", value).Count(&count)
	return count
}
