package repository

import (
	"darulabror/internal/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ArticleRepo interface {
	CreateArticle(title string, content datatypes.JSON, author string) error
	GetAllArticles() (interface{}, error)
	GetArticleByID(id uint) (interface{}, error)
	UpdateArticle(id uint, title string, content datatypes.JSON, status string) error
	DeleteArticle(id uint) error
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) ArticleRepo {
	return &articleRepo{db: db}
}

func (a *articleRepo) CreateArticle(title string, content datatypes.JSON, author string) error {
	article := models.Article{
		Title:   title,
		Content: content,
		Author:  author,
	}

	err := a.db.Create(&article).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *articleRepo) GetAllArticles() (interface{}, error) {
	var articles []models.Article
	err := a.db.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *articleRepo) GetArticleByID(id uint) (interface{}, error) {
	var article models.Article
	err := a.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (a *articleRepo) UpdateArticle(id uint, title string, content datatypes.JSON, status string) error {
	var article models.Article
	err := a.db.First(&article, id).Error
	if err != nil {
		return err
	}

	article.Title = title
	article.Content = content
	article.Status = status

	err = a.db.Save(&article).Error
	if err != nil {
		return err
	}

	return nil
}

func (a *articleRepo) DeleteArticle(id uint) error {
	var article models.Article
	err := a.db.First(&article, id).Error
	if err != nil {
		return err
	}

	err = a.db.Delete(&article).Error
	if err != nil {
		return err
	}

	return nil
}
