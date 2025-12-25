package database

import (
	"fmt"
	"minos/config"
	"minos/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(
		&model.PromptTemplate{},
		&model.Interview{},
		&model.Message{},
		&model.Submission{},
		&model.Evaluation{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create unique index. Don't know how GORM can't handle this.
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_prompt_templates_name_version ON prompt_templates(name, version) WHERE deleted_at IS NULL").Error
	if err != nil {
		return nil, fmt.Errorf("failed to create unique index: %w", err)
	}

	return db, nil
}
