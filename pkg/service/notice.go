package service

import (
	"github.com/google/uuid"
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/repository"
)

type NoticeService struct {
	noticeRepo repository.NoticeRepository
}

func NewNoticeService(noticeRepo repository.NoticeRepository) *NoticeService {
	return &NoticeService{
		noticeRepo: noticeRepo,
	}
}

func (s *NoticeService) Create(notice *entity.Notice) (*entity.Notice, error) {
	//TODO implement me
	panic("implement me")
}

func (s *NoticeService) GetAll() ([]*entity.Notice, error) {
	return s.noticeRepo.GetAll()
}

func (s *NoticeService) FindByID(id uuid.UUID) (*entity.Notice, error) {
	//TODO implement me
	panic("implement me")
}

func (s *NoticeService) Update(notice *entity.Notice) (*entity.Notice, error) {
	//TODO implement me
	panic("implement me")
}
