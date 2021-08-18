package service

import (
	"regexp"

	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/error"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo        repository.UserRepository
	userProfileRepo repository.UserProfileRepository
	userGroupRepo   repository.UserGroupRepository
}

func NewUserService(userRepo repository.UserRepository, userProfileRepo repository.UserProfileRepository, userGroupRepo repository.UserGroupRepository) *UserService {
	return &UserService{
		userRepo:        userRepo,
		userProfileRepo: userProfileRepo,
		userGroupRepo:   userGroupRepo,
	}
}

func (s *UserService) validateUserName(name string) bool {
	matched, err := regexp.MatchString("[A-Za-z0-9_]{1,32}", name)
	return err != nil && matched
}

func (s *UserService) Create(name, password string, userGroupID uuid.UUID, invitationCode string) (*entity.User, error) {
	if !s.validateUserName(name) {
		return nil, error.NewBadRequestError("invalid name")
	}

	userGroup, err := s.userGroupRepo.FindByID(userGroupID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userGroup.InvitationCodeDigest), []byte(invitationCode)); err != nil {
		return nil, err
	}

	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.userRepo.Create(&entity.User{
		Name:           name,
		DisplayName:    name,
		PasswordDigest: string(digest),
		UserGroupID:    userGroupID,
	})
}

func (s *UserService) FindByID(id uuid.UUID) (*entity.User, error) {
	return s.userRepo.FindByID(id, true)
}

func (s *UserService) Update(userID uuid.UUID, displayName, twitterID, githubID, facebookID, selfIntroduction string) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID, false)
	if err != nil {
		return nil, err
	}

	user.DisplayName = displayName
	if _, err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	userProfile, err := s.userProfileRepo.FindByUserID(userID)
	//Create
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return nil, err
		}
		if _, err := s.userProfileRepo.Create(&entity.UserProfile{
			UserID:           userID,
			TwitterID:        twitterID,
			GithubID:         githubID,
			FacebookID:       facebookID,
			SelfIntroduction: selfIntroduction,
		}); err != nil {
			return nil, err
		}

		return s.userRepo.FindByID(userID, true)
	}
	//Update
	userProfile.TwitterID = twitterID
	userProfile.GithubID = githubID
	userProfile.FacebookID = facebookID
	userProfile.SelfIntroduction = selfIntroduction
	if _, err := s.userProfileRepo.Update(userProfile); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(userID, true)
}
