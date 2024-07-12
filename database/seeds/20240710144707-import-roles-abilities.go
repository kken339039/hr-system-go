package seeds

import (
	"hr-system-go/internal/auth/constants"
	"hr-system-go/internal/auth/models"

	"gorm.io/gorm"
)

func init() {
	Seeds = append(Seeds, Seed{
		Name: "20240710144707-import-roles-abilities",
		Exec: Exec_20240710144707,
	})
}

func Exec_20240710144707(db *gorm.DB) error {
	abilities := []models.Ability{
		{Name: constants.ABILITY_READ_USER},
		{Name: constants.ABILITY_READ_WRITE_USER},
		{Name: constants.ABILITY_DELETE_USER},
		{Name: constants.ABILITY_READ_LEAVE},
		{Name: constants.ABILITY_READ_WRITE_LEAVE},
		{Name: constants.ABILITY_DELETE_LEAVE},
		{Name: constants.ABILITY_READ_CLOCK_RECORD},
		{Name: constants.ABILITY_READ_WRITE_CLOCK_RECORD},
		{Name: constants.ABILITY_DELETE_CLOCK_RECORD},
		{Name: constants.ABILITY_READ_DEPARTMENT},
		{Name: constants.ABILITY_READ_WRITE_DEPARTMENT},
		{Name: constants.ABILITY_DELETE_DEPARTMENT},
		{Name: constants.ABILITY_ADMIN},
		{Name: constants.ABILITY_ALL_GRANTS_USER},
		{Name: constants.ABILITY_ALL_GRANTS_LEAVE},
		{Name: constants.ABILITY_ALL_GRANTS_CLOCK_RECORD},
		{Name: constants.ABILITY_ALL_GRANTS_DEPARTMENT},
	}

	for i := 0; i < len(abilities); i++ {
		var existingAbility models.Ability
		result := db.Where("name = ?", abilities[i].Name).FirstOrCreate(&existingAbility, models.Ability{Name: abilities[i].Name})
		if result.Error != nil {
			return result.Error
		}
		abilities[i] = existingAbility
	}

	roles := []struct {
		Name      string
		Abilities []string
	}{
		{
			Name: constants.ROLE_RD,
			Abilities: []string{
				constants.ABILITY_READ_USER, constants.ABILITY_READ_WRITE_USER,
				constants.ABILITY_READ_LEAVE, constants.ABILITY_READ_WRITE_LEAVE, constants.ABILITY_DELETE_LEAVE,
				constants.ABILITY_READ_CLOCK_RECORD, constants.ABILITY_READ_WRITE_CLOCK_RECORD,
				constants.ABILITY_READ_DEPARTMENT,
			},
		},
		{
			Name: constants.ROLE_RD_MANAGER,
			Abilities: []string{
				constants.ABILITY_READ_USER, constants.ABILITY_READ_WRITE_USER,
				constants.ABILITY_ALL_GRANTS_LEAVE,
				constants.ABILITY_READ_CLOCK_RECORD, constants.ABILITY_READ_WRITE_CLOCK_RECORD,
				constants.ABILITY_READ_DEPARTMENT, constants.ABILITY_DELETE_LEAVE, constants.ABILITY_READ_WRITE_DEPARTMENT,
			},
		},
		{
			Name: constants.ROLE_BD,
			Abilities: []string{
				constants.ABILITY_READ_USER, constants.ABILITY_READ_WRITE_USER,
				constants.ABILITY_READ_LEAVE, constants.ABILITY_READ_WRITE_LEAVE,
				constants.ABILITY_READ_CLOCK_RECORD, constants.ABILITY_READ_WRITE_CLOCK_RECORD,
				constants.ABILITY_READ_DEPARTMENT,
			},
		},
		{
			Name: constants.ROLE_BD_MANAGER,
			Abilities: []string{
				constants.ABILITY_READ_USER, constants.ABILITY_READ_WRITE_USER,
				constants.ABILITY_ALL_GRANTS_LEAVE,
				constants.ABILITY_READ_CLOCK_RECORD, constants.ABILITY_READ_WRITE_CLOCK_RECORD,
				constants.ABILITY_READ_DEPARTMENT, constants.ABILITY_DELETE_LEAVE, constants.ABILITY_READ_WRITE_DEPARTMENT,
			},
		},
		{
			Name: constants.ROLE_HR,
			Abilities: []string{
				constants.ABILITY_ALL_GRANTS_USER,
				constants.ABILITY_ALL_GRANTS_LEAVE,
				constants.ABILITY_ALL_GRANTS_CLOCK_RECORD,
				constants.ABILITY_ALL_GRANTS_DEPARTMENT,
			},
		},
		{
			Name:      constants.ROLE_HR_MANAGER,
			Abilities: []string{constants.ABILITY_ADMIN},
		},
		{
			Name: constants.ROLE_INTERN,
			Abilities: []string{
				constants.ABILITY_READ_USER, constants.ABILITY_READ_WRITE_USER,
				constants.ABILITY_READ_LEAVE, constants.ABILITY_READ_WRITE_LEAVE, constants.ABILITY_DELETE_LEAVE,
				constants.ABILITY_READ_CLOCK_RECORD, constants.ABILITY_READ_WRITE_CLOCK_RECORD,
				constants.ABILITY_READ_DEPARTMENT,
			},
		},
	}

	for _, roleInfo := range roles {
		var role models.Role
		result := db.Where("name = ?", roleInfo.Name).FirstOrCreate(&role, models.Role{Name: roleInfo.Name})
		if result.Error != nil {
			return result.Error
		}

		var roleAbilities []models.Ability
		for _, abilityName := range roleInfo.Abilities {
			for _, ability := range abilities {
				if ability.Name == abilityName {
					roleAbilities = append(roleAbilities, ability)
					break
				}
			}
		}

		if err := db.Model(&role).Association("Abilities").Replace(roleAbilities); err != nil {
			return err
		}
	}
	return nil
}
