package models

import (
	"time"
	"weblog/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `json:"id" gorm:"size:12;primaryKey;unique;notnull;comment:用户id"`
	Name      string         `json:"name" gorm:"size:60;comment:用户名"`
	Email     string         `json:"email" gorm:"type:varchar(100);not null;uniqueIndex:idx_email_deleted;comment:邮箱"`
	Password  string         `json:"-"  gorm:"size:33;notnull;comment:登录密码"`
	Group     string         `json:"group" gorm:"size:9;notnull;comment:用户分组(管理员;用户)"`
	Status    int8           `json:"status" gorm:"size:1;notnull;comment:允许登录(1启用;2关闭)"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;notnull;comment:user创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:上一次修改信息时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"uniqueIndex:idx_email_deleted;precision:3"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// WithTx 【核心】WithTx 产生一个带有事务连接的新 Repo 实例
func (repo *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx} // 将原本的 db 替换为事务 tx
}

// LoginByEmail 根据邮箱获取能够登录的用户信息
func (repo *UserRepository) LoginByEmail(email string) (*User, error) {
	var user User
	err := repo.db.Model(&User{}).
		Where("email = ? AND status = 1", email).
		First(&user).Error

	if err != nil {
		// 如果找不到记录，或者数据库报错，直接向上抛出 err
		return nil, err
	}

	return &user, nil
}

// Create 注册
func (repo *UserRepository) Create(user *User) (userId uint64, rowsAffected int64, err error) {
	result := repo.db.Create(user)
	return user.ID, result.RowsAffected, result.Error
}

// UpdateProfile 编辑
func (repo *UserRepository) UpdateProfile(userId uint64, userData map[string]any) (int64, error) {
	result := repo.db.Model(&User{}).Where("id = ?", userId).Updates(userData)
	return result.RowsAffected, result.Error
}

// UpdatePassword 校验并更新用户密码
func (repo *UserRepository) UpdatePassword(userId uint64, oldPwd, newPwd string) (bool, error) {
	// 拿当前用户数据库里的密文密码
	var user User
	if err := repo.db.First(&user, userId).Error; err != nil {
		return false, err
	}

	// 验证旧密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd))
	if err != nil {
		// 旧密码对不上，返回 false，代表校验失败
		return false, nil
	}

	// 生成新密码的哈希密文
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	// 只更新密码和更新时间字段
	err = repo.db.Model(&user).Updates(map[string]any{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error

	return true, err
}

// ForceUpdatePassword 管理员强制更新指定用户的密码（无需旧密码验证）
func (repo *UserRepository) ForceUpdatePassword(targetUserId uint64, newPwd string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 根据指定的 targetUserId 强行更新密码字段
	err = repo.db.Model(&User{}).Where("id = ?", targetUserId).Updates(map[string]any{
		"password":   string(hashedPassword),
		"updated_at": time.Now(),
	}).Error

	return err
}

// Delete 删除
func (repo *UserRepository) Delete(userId uint64) int64 {
	return repo.db.Delete(&User{}, userId).RowsAffected
}

// UnregisterSelf 用户自主注销（验密后软删除）
func (repo *UserRepository) UnregisterSelf(userId uint64, password string) (bool, error) {
	var user User
	if err := repo.db.First(&user, userId).Error; err != nil {
		return false, err
	}

	// 校验输入的密码是否正确
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil
	}

	err = repo.db.Delete(&user).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// CheckEmailExist 邮箱唯一性检查
func (repo *UserRepository) CheckEmailExist(email string) (bool, error) {
	var count int64
	err := repo.db.Model(&User{}).
		Where("email != '' AND email = ?", email).
		Count(&count).Error

	return count > 0, err
}

// GetUserDetailsById 获取用户详情ById
func (repo *UserRepository) GetUserDetailsById(userId uint64) (*User, error) {
	var user User
	err := repo.db.First(&user, userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserDetailsByEmail 获取用户详情ByEmail
func (repo *UserRepository) GetUserDetailsByEmail(userEmail string) (*User, error) {
	var user User
	err := repo.db.Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetList 查列表
func (repo *UserRepository) GetList(in dto.UserListInput) (total int64, list []User, err error) {
	tx := repo.db.Model(&User{})

	// 动态拼接筛选条件
	if in.Group != "" {
		tx = tx.Where("`group` = ?", in.Group)
	}

	// 后续加上按用户名模糊搜索
	// if in.Keyword != "" {
	//     tx = tx.Where("`name` LIKE ?", "%"+in.Keyword+"%")
	// }

	// 使用 Session 隔离句柄，防止 Count 操作污染后续的 Find，这样 txCount 是独立的，tx 依然保持纯净
	txCount := tx.Session(&gorm.Session{})
	if err = txCount.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	// 如果总数本来就是 0，直接打道回府，连 Limit 都不用跑了，省一次数据库 I/O
	if total == 0 {
		return 0, make([]User, 0), nil
	}

	// 动态排序（限制前端只能传入合法的排序字段，防止 SQL 注入或报错）
	if in.OrderBy != "" {
		// 只允许按 id, created_time, updated_time 排序
		if in.OrderBy == "id" || in.OrderBy == "created_at" || in.OrderBy == "updated_at" {
			tx = tx.Order(in.OrderBy + " DESC") // 或者根据前端传的方向拼 ASC/DESC
		} else {
			tx = tx.Order("id DESC") // 恶意或非法字段，降级走默认排序
		}
	} else {
		tx = tx.Order("id DESC") // 为空时默认按 ID 降序
	}

	// 安全分页查询
	offset := (in.PageNo - 1) * in.PageSize
	err = tx.Limit(in.PageSize).Offset(offset).Find(&list).Error

	return total, list, err
}

// 事务

//func (repo *UserRepository) TxAdd(tx *gorm.DB) (userId uint64, rowsAffected int64, err error) {
//	var user User
//	result := tx.Create(&user)
//	return user.Id, result.RowsAffected, result.Error
//}
//
//func (repo *UserRepository) TxEdit(tx *gorm.DB, id uint64) (rowsAffected int64, err error) {
//	var user User
//	result := tx.Model(&User{}).Where("id = ?", id).Updates(user)
//	return result.RowsAffected, result.Error
//}

// TxCountTotalUsers   事务查询用户总数
//func (repo *UserRepository) TxCountTotalUsers(tx *gorm.DB) (count int64, err error) {
//	if err = tx.Model(&User{}).Count(&count).Error; err != nil {
//		return 0, err
//	}
//	return count, nil
//}
