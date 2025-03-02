package builder

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
	"runtime"
)

// MyGroup 方便扩展和打印group
type MyGroup struct {
	name        string
	fatherGroup *MyGroup
	middlewares []string

	G *gin.RouterGroup
}

func GetMyGroupDetail(group *MyGroup) {

	fmt.Printf("Group Name: %s\n", group.name)

	if group.fatherGroup.name != "" {
		fmt.Printf("Father Group: %s\n", group.fatherGroup.name)
	} else {
		fmt.Println("Father Group: None")
	}

	fmt.Println("Middlewares:")
	for _, middleware := range group.middlewares {
		fmt.Printf("  - %s\n", middleware)
	}
	fmt.Printf("\n\n")
}

type GroupBuilder struct {
	name           string
	newMiddlewares []gin.HandlerFunc
	path           string //一个添加的路径
	group          *MyGroup
	routes         func(*gin.RouterGroup)
}

func NewGroupBuilder() *GroupBuilder {
	return &GroupBuilder{}
}

func (b *GroupBuilder) SetFatherGroup(entity *MyGroup) *GroupBuilder {
	b.group = entity
	//b.path = entity.G.BasePath()
	return b
}

func (b *GroupBuilder) SetName(name string) *GroupBuilder {
	b.name = name
	return b
}

func (b *GroupBuilder) SetRoutes(routerFunc func(*gin.RouterGroup)) *GroupBuilder {
	b.routes = routerFunc
	return b
}

type Routes interface {
	Router(*gin.RouterGroup)
}

func (b *GroupBuilder) AddMiddleware(middleware ...gin.HandlerFunc) *GroupBuilder {
	if b.group == nil || b.group.G == nil {
	}
	b.newMiddlewares = append(b.newMiddlewares, middleware...)
	return b
}

func (b *GroupBuilder) AddRoute(path string) *GroupBuilder {
	b.path = path
	return b
}

func (b *GroupBuilder) Build() *MyGroup {
	if b.group == nil {
		log.Panic("fatherGroup is nil")
	}
	if b.routes == nil {
		log.Panic("routes function is nil")
	}

	group := b.group.G.Group(b.path)

	var middlewaresAll []string
	middlewaresAll = append(middlewaresAll, b.group.middlewares...)
	for _, middleware := range b.newMiddlewares {
		group.Use(middleware)

		middlewareName := runtime.FuncForPC(reflect.ValueOf(middleware).Pointer()).Name()
		middlewaresAll = append(middlewaresAll, middlewareName)
	}
	b.routes(group)
	return &MyGroup{
		name:        b.name,
		fatherGroup: b.group,
		middlewares: middlewaresAll,
		G:           group,
	}
}
