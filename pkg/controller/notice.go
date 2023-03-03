package controller

import (
	"github.com/ictsc/ictsc-rikka/pkg/entity"
	"github.com/ictsc/ictsc-rikka/pkg/service"
)

type NoticeController struct {
	noticeService *service.NoticeService
}

func NewNoticeController(noticeService *service.NoticeService) *NoticeController {
	return &NoticeController{
		noticeService: noticeService,
	}
}

func (c *NoticeController) GetAll() ([]*entity.Notice, error) {
	return c.noticeService.GetAll()
}
