package models

import (
	base_model "hr-system-go/internal/base/models"
)

type Role struct {
	base_model.BaseModel
	Name      string    `gorm:"not null"`
	Status    string    `gorm:"default:'active'"`
	Abilities []Ability `gorm:"many2many:role_abilities;"`
}

func (r *Role) GetAbilityNames() []string {
	abilityNames := make([]string, 0, len(r.Abilities))
	for _, ability := range r.Abilities {
		abilityNames = append(abilityNames, ability.Name)
	}
	return abilityNames
}
