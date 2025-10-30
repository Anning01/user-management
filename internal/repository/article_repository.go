package repository

import (
	"github.com/Anning01/user-management/internal/domain"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *domain.Article) error
	FindByID(id uint) (*domain.Article, error)
	FindAll(limit, offset int) ([]domain.Article, int64, error)
	FindByAuthorID(authorID uint, limit, offset int) ([]domain.Article, int64, error)
	Update(article *domain.Article) error
	Delete(id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) Create(article *domain.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) FindByID(id uint) (*domain.Article, error) {
	var article domain.Article
	if err := r.db.Preload("Author").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindAll(limit, offset int) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	if err := r.db.Model(&domain.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Author").Limit(limit).Offset(offset).Order("created_at desc").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articleRepository) FindByAuthorID(authorID uint, limit, offset int) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	query := r.db.Model(&domain.Article{}).Where("author_id = ?", authorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articleRepository) Update(article *domain.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}
