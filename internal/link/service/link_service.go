package service

import (
	"context"
	"errors"
	"log"

	"github.com/suryaadi44/linkify/internal/link/dto"
	linkRepository "github.com/suryaadi44/linkify/internal/link/repository"
	userService "github.com/suryaadi44/linkify/internal/user/service"
)

type LinkService struct {
	linkRepository linkRepository.LinkRepository
	userService    userService.UserService
}

var ErrLinkExists = errors.New("Link already exists")

func NewLinkService(linkRepository linkRepository.LinkRepository, userService userService.UserService) *LinkService {
	return &LinkService{
		linkRepository: linkRepository,
		userService:    userService,
	}
}

func (l LinkService) GetLink(ctx context.Context, username string) (*dto.LinksResponse, error) {
	links, err := l.linkRepository.GetLinkByUsername(ctx, username)
	if err != nil {
		log.Println("[Link] Error fetching links: ", err)
		return nil, err
	}

	picture, err := l.userService.GetUserPictureByUsername(ctx, username)
	if err != nil {
		log.Println("[Link] Error fetching picture: ", err)
		return nil, err
	}

	return dto.NewLinksResponse(*links, picture), nil
}

func (l LinkService) AddLink(ctx context.Context, username string, link dto.Link) error {
	// Check duplicate link title
	if l.linkRepository.IsLinkExists(ctx, username, link.Title) {
		return ErrLinkExists
	}

	// Add link to document links field (array)
	linkEntity := dto.NewLinkEntity(link)
	err := l.linkRepository.AddLink(ctx, username, *linkEntity)
	if err != nil {
		log.Println("[Link] Error adding link: ", err)
		return err
	}

	return nil
}
