package dynmgrm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

type callbacksRegisterer struct{}

func (c *callbacksRegisterer) Register(db *gorm.DB, config *callbacks.Config) {
	callbacks.RegisterDefaultCallbacks(db, config)
}
