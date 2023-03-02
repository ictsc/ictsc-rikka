package seed

import (
	"fmt"
	"log"

	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"github.com/ictsc/ictsc-rikka/pkg/service"
	"github.com/pkg/errors"
)

type SeedConfig struct {
	AdminUserGroupAndUserSeeder AdminUserGroupAndUserSeederConfig `yaml:"adminUserGroupAndUserSeeder"`
}

type seeder struct {
	userService      service.UserService
	userGroupService service.UserGroupService
	bastionService   service.BastionService
	userRepo         repository.UserRepository
	userGroupRepo    repository.UserGroupRepository
	bastionRepo      repository.BastionRepository
}

func Seed(
	config *SeedConfig,
	userRepo repository.UserRepository,
	userGroupRepo repository.UserGroupRepository,
	userService service.UserService,
	userGroupService service.UserGroupService,
	bastionService service.BastionService,
	bastionRepo repository.BastionRepository,
) error {
	s := seeder{
		userService:      userService,
		userGroupService: userGroupService,
		bastionService:   bastionService,
		userRepo:         userRepo,
		userGroupRepo:    userGroupRepo,
		bastionRepo:      bastionRepo,
	}

	seeds := []struct {
		name string
		fn   func() error
	}{
		{
			name: "adminUserGroupAndUserSeeder",
			fn: func() error {
				return s.adminUserGroupAndUserSeeder(config.AdminUserGroupAndUserSeeder)
			},
		},
	}
	for _, s := range seeds {
		log.Printf("start seeding... `%s`\n", s.name)
		if err := s.fn(); err != nil {
			if err.Error() != "Skip" {
				return errors.Wrap(err, s.name)
			}
			log.Printf("Skip `%s`\n", s.name)
		}
		log.Printf("success `%s`\n", s.name)
	}

	return nil
}

type AdminUserGroupAndUserSeederConfig struct {
	UserGroupName   string `yaml:"userGroupName"`
	Organization    string `yaml:"organization"`
	InvitationCode  string `yaml:"invitationCode"`
	UserName        string `yaml:"userName"`
	UserPassword    string `yaml:"userPassword"`
	BastionUser     string `yaml:"bastionUser"`
	BastionPassword string `yaml:"bastionPassword"`
	BastionHost     string `yaml:"bastionHost"`
	BastionPort     int    `yaml:"bastionPort"`
	TeamID          string `yaml:"teamId"`
}

func (s *seeder) adminUserGroupAndUserSeeder(config AdminUserGroupAndUserSeederConfig) error {
	userGroup, err := s.userGroupRepo.FindByName(config.UserGroupName)
	if err != nil {
		return err
	}

	if userGroup != nil {
		return fmt.Errorf("Skip")
	}

	userGroup, err = s.userGroupService.Create(
		config.UserGroupName,
		config.Organization,
		config.InvitationCode,
		true,
		config.TeamID,
	)
	if err != nil {
		return err
	}

	if _, err := s.userService.Create(
		config.UserName,
		config.UserPassword,
		userGroup.ID,
		config.InvitationCode,
	); err != nil {
		return err
	}

	if _, err := s.bastionService.Create(
		userGroup.ID,
		config.BastionUser,
		config.BastionPassword,
		config.BastionHost,
		config.BastionPort,
	); err != nil {
		return err
	}

	return nil
}
