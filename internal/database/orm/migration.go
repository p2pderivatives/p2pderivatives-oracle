package orm

import (
	"github.com/pkg/errors"
)

// Migrator is a structure used to perform DB migrations.
type Migrator struct {
	orm         *ORM
	initialized bool
	models      []interface{}
}

// NewMigrator creates a new Migrator class.
// ex.) NewMigrator(orm, &model.Hoge{}, &model.Fuga{})
func NewMigrator(orm *ORM, models ...interface{}) *Migrator {
	return &Migrator{
		orm:         orm,
		initialized: false,
		models:      models,
	}
}

// Initialize performs the migrations for the models handled by the Migrator.
func (m *Migrator) Initialize() error {
	for _, model := range m.models {
		if err := m.orm.GetDB().AutoMigrate(model).Error; err != nil {
			return errors.Errorf("migration failed for [%v]", model)
		}
	}

	m.initialized = true

	return nil
}

// IsInitialized returns whether the Migrator is initialized.
func (m *Migrator) IsInitialized() bool {
	return m.initialized
}
