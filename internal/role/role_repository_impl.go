package role

import (
	"context"
	"encoding/json"
	"fmt"
	"go-intconnect-api/configs"
	"go-intconnect-api/internal/entity"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	redisInstance      *configs.RedisInstance
	redisRolesCacheKey string
	redisRoleCacheKey  string
}

func NewRepository(
	redisInstance *configs.RedisInstance,
	viperConfig *viper.Viper,
) *RepositoryImpl {
	return &RepositoryImpl{
		redisInstance:      redisInstance,
		redisRolesCacheKey: viperConfig.GetString("REDIS_ROLES_KEY"),
		redisRoleCacheKey:  viperConfig.GetString("REDIS_ROLE_KEY"),
	}
}

func (roleRepositoryImpl *RepositoryImpl) FindAll(gormTransaction *gorm.DB) ([]entity.Role, error) {
	var roleEntities []entity.Role
	err := gormTransaction.Preload("Permissions").Find(&roleEntities).Error
	return roleEntities, err
}

func (roleRepositoryImpl *RepositoryImpl) FindById(gormTransaction *gorm.DB, roleId uint64) (*entity.Role, error) {
	var roleEntity entity.Role
	err := gormTransaction.
		Model(&entity.Role{}).
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "code")
		}).
		Where("id = ?", roleId).
		First(&roleEntity).Error
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
	data, err := roleRepositoryImpl.redisInstance.RedisClient.Get(context, roleRepositoryImpl.redisRolesCacheKey).Result()
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

func (roleRepositoryImpl *RepositoryImpl) SetAllCache(ctx context.Context, roles []entity.Role) error {
	data, err := json.Marshal(roles)
	if err != nil {
		return err
	}
	return roleRepositoryImpl.redisInstance.RedisClient.Set(ctx, roleRepositoryImpl.redisRolesCacheKey, data, 0).Err()
}

func (roleRepositoryImpl *RepositoryImpl) DeleteAllCache(ctx context.Context) error {
	return roleRepositoryImpl.redisInstance.RedisClient.Del(ctx, roleRepositoryImpl.redisRolesCacheKey).Err()
}

func (roleRepositoryImpl *RepositoryImpl) FindRoleCacheById(ctx context.Context, roleId uint64) (*entity.Role, error) {
	redisKey := fmt.Sprintf("%s%d", roleRepositoryImpl.redisRoleCacheKey, roleId)

	data, err := roleRepositoryImpl.redisInstance.RedisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var role entity.Role
	if err := json.Unmarshal([]byte(data), &role); err != nil {
		return nil, err
	}

	return &role, nil
}

func (roleRepositoryImpl *RepositoryImpl) SetByIdCache(ctx context.Context, roleId uint64, roleEntity *entity.Role) error {
	redisKey := fmt.Sprintf("%s%d", roleRepositoryImpl.redisRoleCacheKey, roleId)

	roleEntityPayload, err := json.Marshal(roleEntity)
	if err != nil {
		return err
	}

	return roleRepositoryImpl.redisInstance.RedisClient.Set(ctx, redisKey, roleEntityPayload, 0).Err()
}
func (roleRepositoryImpl *RepositoryImpl) DeleteByIdCache(ctx context.Context, roleId uint64) error {
	redisKey := fmt.Sprintf("%s%d", roleRepositoryImpl.redisRoleCacheKey, roleId)
	return roleRepositoryImpl.redisInstance.RedisClient.Del(ctx, redisKey).Err()
}
