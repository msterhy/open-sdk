package libx

import (
	"github.com/gin-gonic/gin"
	"github.com/trancecho/open-sdk/cache/types"
	"gorm.io/gorm"
)

func Uid(c *gin.Context) uint {
	uid := c.MustGet("uid").(uint)
	uidInt := uid
	return uidInt
}

func GetUsername(c *gin.Context) string {
	username := c.MustGet("username").(string)
	return username
}

func GetRole(c *gin.Context) string {
	role := c.MustGet("role").(string)
	return role
}

// todo :把email也存了

// todo：怎么拿到该service的数据库和缓存
func GetService(c *gin.Context) string {
	service := c.MustGet("service").(string)
	return service
}

func GetDb(c *gin.Context) *gorm.DB {
	db := c.MustGet("db").(*gorm.DB)
	return db
}

func GetRds(c *gin.Context) types.Cache {
	rds := c.MustGet("rds").(types.Cache)
	return rds
}

//func GetCurrentMundoUser(c *gin.Context) (user_entity.MundoUser, error) {
//	uid := Uid(c)
//	var user user_entity.MundoUser
//	err := dbx.DbMundo.Where("uid = ?", uid).First(&user).Error
//	if err != nil {
//		return user, err
//	}
//	return user, nil
//}
