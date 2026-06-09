package models

import (
	"errors"
	"time"
	"weblog/dto"

	"gorm.io/gorm"
)

type SystemConfig struct {
	Id        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Key       string    `json:"key" gorm:"type:varchar(60);not null;uniqueIndex;comment:配置项唯一键名(如: site_title)"`
	Value     string    `json:"value" gorm:"type:text;not null;comment:配置项参数值"`
	Remark    string    `json:"remark" gorm:"type:varchar(250);not null;default:'';comment:配置项备注介绍"`
	Status    int8      `json:"status" gorm:"type:tinyint;not null;default:1;comment:状态: 1启用 2禁用"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:编辑时间"`
}

type SystemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository 初始化时注入 DB 实例
func NewSystemConfigRepository(db *gorm.DB) *SystemConfigRepository {
	return &SystemConfigRepository{db: db}
}

// WithTx  产生一个带有事务连接的新 Repo 实例
func (repo *SystemConfigRepository) WithTx(tx *gorm.DB) *SystemConfigRepository {
	return &SystemConfigRepository{db: tx} // 将原本的 db 替换为事务 tx
}

// CreateConfig 新增一条系统配置
func (repo *SystemConfigRepository) CreateConfig(config *SystemConfig) error {
	return repo.db.Create(config).Error
}

// DeleteConfigs 批量删除配置项（安全版）
func (repo *SystemConfigRepository) DeleteConfigs(ids []uint64) error {
	// GORM 的 IN 查询会自动处理切片
	tx := repo.db.Where("id IN ?", ids).Delete(&SystemConfig{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// GetConfigList 纯粹且安全的配置项分页查询
func (repo *SystemConfigRepository) GetConfigList(in dto.ConfigListInput) ([]SystemConfig, int64, error) {
	var list []SystemConfig
	var total int64

	// 构建初始查询对象
	tx := repo.db.Model(&SystemConfig{})

	// 动态拼接筛选条件（按 Key 或 备注 模糊搜索）
	if in.Keyword != "" {
		likeKeyword := "%" + in.Keyword + "%"
		tx = tx.Where("`key` LIKE ? OR `remark` LIKE ?", likeKeyword, likeKeyword)
	}

	// 使用 Session 隔离 Count 句柄，彻底杜绝链式污染
	txCount := tx.Session(&gorm.Session{})
	if err := txCount.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接秒回空数组，免去后面的 Limit 数据库开销
	if total == 0 {
		return make([]SystemConfig, 0), 0, nil
	}

	// 执行安全分页，默认按 ID 降序排列
	offset := (in.PageNo - 1) * in.PageSize
	err := tx.Order("id DESC").Limit(in.PageSize).Offset(offset).Find(&list).Error

	return list, total, err
}

// GetValueByKey 根据 Key 快速获取配置值（有 uniqueIndex，这里速度极快）
func (repo *SystemConfigRepository) GetValueByKey(key string) (string, error) {
	var config SystemConfig
	// 性能优化：只查 value 字段，不查 id, remark, time，减少 MySQL 网络 I/O
	err := repo.db.Select("value").Where("`key` = ?", key).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}
	return config.Value, nil
}

// GetConfigById 通过主键ID获取配置完整详情
func (repo *SystemConfigRepository) GetConfigById(id uint64) (*SystemConfig, error) {
	var config SystemConfig
	// 通过主键查询，GORM 会自动优化为：SELECT * FROM `system_configs` WHERE `id` = id LIMIT 1
	err := repo.db.First(&config, id).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateValueByKey 根据 Key 安全更新配置的 Value 和 Remark
func (repo *SystemConfigRepository) UpdateValueByKey(key string, value string, remark string) error {
	// 🌟 性能与安全双优解：
	// Where("`key` = ?", key) 精准锁定唯一行
	// Select("value", "remark") 锁死更新白名单，防止结构体里的其他字段（如 id 或 key 本身）被意外修改
	// models 里有 UpdatedAt，GORM 会自动将 updated_at 刷新为当前系统时间（time.Time）
	err := repo.db.Model(&SystemConfig{}).
		Where("`key` = ?", key).
		Select("value", "remark").
		Updates(&SystemConfig{
			Value:  value,
			Remark: remark,
		}).Error

	return err
}
