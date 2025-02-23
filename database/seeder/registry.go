package seeder

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type ISeedRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) ISeedRegistry {
	return &Registry{db: db}
}

func (s *Registry) Run() {
	RoleSeeder(*s.db)
	UserSeeder(*s.db)
}
