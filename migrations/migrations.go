package migrations

import (
	"user-management/internal/domain"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions(), []*gormigrate.Migration{
		{
			ID: "20230101000001",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&domain.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&domain.User{})
			},
		},
		{
			ID: "20230101000002",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&domain.Article{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&domain.Article{})
			},
		},
	})

	return m.Migrate()
}
