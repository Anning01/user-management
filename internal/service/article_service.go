package service

import (
	"errors"

	"github.com/Anning01/user-management/internal/domain"
	"github.com/Anning01/user-management/internal/repository"
)

type ArticleService interface {
	CreateArticle(article *domain.Article) error
	GetArticleByID(id uint) (*domain.Article, error)
	ListArticles(page, pageSize int) ([]domain.Article, int64, error)
	ListArticlesByAuthor(authorID uint, page, pageSize int) ([]domain.Article, int64, error)
	UpdateArticle(id, authorID uint, title, content string) error
	DeleteArticle(id, authorID uint) error
}

type articleService struct {
	articleRepo repository.ArticleRepository
	userRepo    repository.UserRepository
}

func NewArticleService(articleRepo repository.ArticleRepository, userRepo repository.UserRepository) ArticleService {
	return &articleService{
		articleRepo: articleRepo,
		userRepo:    userRepo,
	}
}

func (s *articleService) CreateArticle(article *domain.Article) error {
	// 验证作者是否存在
	_, err := s.userRepo.FindByID(article.AuthorID)
	if err != nil {
		return errors.New("author not found")
	}

	return s.articleRepo.Create(article)
}

func (s *articleService) GetArticleByID(id uint) (*domain.Article, error) {
	return s.articleRepo.FindByID(id)
}

func (s *articleService) ListArticles(page, pageSize int) ([]domain.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.articleRepo.FindAll(pageSize, offset)
}

func (s *articleService) ListArticlesByAuthor(authorID uint, page, pageSize int) ([]domain.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.articleRepo.FindByAuthorID(authorID, pageSize, offset)
}

func (s *articleService) UpdateArticle(id, authorID uint, title, content string) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	// 验证是否是文章作者
	if article.AuthorID != authorID {
		return errors.New("permission denied")
	}

	article.Title = title
	article.Content = content

	return s.articleRepo.Update(article)
}

func (s *articleService) DeleteArticle(id, authorID uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("article not found")
	}

	// 验证是否是文章作者
	if article.AuthorID != authorID {
		return errors.New("permission denied")
	}

	return s.articleRepo.Delete(id)
}
