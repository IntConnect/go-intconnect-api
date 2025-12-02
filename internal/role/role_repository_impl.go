package role

import (
	"context"
	"encoding/json"
	"go-intconnect-api/configs"
	"go-intconnect-api/internal/entity"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	redisInstance *configs.RedisInstance
	redisCacheKey string
}

func NewRepository(
	redisInstance *configs.RedisInstance,
	viperConfig *viper.Viper,
) *RepositoryImpl {
	return &RepositoryImpl{
		redisInstance: redisInstance,
		redisCacheKey: viperConfig.GetString("REDIS_ROLE_KEY"),
	}
}

func (roleRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Role, error) {
	var roleEntities []entity.Role
	err := gormTransaction.Preload("Permissions").Find(&roleEntities).Error
	return roleEntities, err
}

func (roleRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error) {
	var roleEntity entity.Role
	err := gormTransaction.Model(&entity.Role{}).Where("id = ?", roleId).Find(&roleEntity).Error

	return &roleEntity, err
}

func (roleRepositoryImpl *RepositoryImpl) Create(gormTransaction *gorm.DB, pipelineEntity *entity.Role) error {
	return gormTransaction.Model(pipelineEntity).Create(pipelineEntity).Error

}

func (roleRepositoryImpl *RepositoryImpl) Update(gormTransaction *gorm.DB, roleEntity *entity.Role) error {
	return gormTransaction.Model(roleEntity).Save(roleEntity).Error
}

func (roleRepositoryImpl *RepositoryImpl) Delete(gormTransaction *gorm.DB, id uint64) error {
	return gormTransaction.Model(&entity.Role{}).Where("id = ?", id).Delete(&entity.Role{}).Error
}

func (roleRepositoryImpl *RepositoryImpl) FindAllCache(context context.Context) ([]entity.Role, error) {
	data, err := roleRepositoryImpl.redisInstance.RedisClient.Get(context, roleRepositoryImpl.redisCacheKey).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var roles []entity.Role
	if err := json.Unmarshal([]byte(data), &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (roleRepositoryImpl *RepositoryImpl) SetAll(ctx context.Context, roles []entity.Role) error {
	data, err := json.Marshal(roles)
	if err != nil {
		return err
	}
	return roleRepositoryImpl.redisInstance.RedisClient.Set(ctx, roleRepositoryImpl.redisCacheKey, data, 0).Err()
}

func (roleRepositoryImpl *RepositoryImpl) DeleteAll(ctx context.Context) error {
	return roleRepositoryImpl.redisInstance.RedisClient.Del(ctx, roleRepositoryImpl.redisCacheKey).Err()
}
