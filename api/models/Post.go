package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:100;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdateAt = time.Now()
}

func (p *Post) Validated() error {

	if p.Title == "" {
		return errors.New("Isi Titlenya!")
	}

	if p.Content == "" {
		return errors.New("Isi Contentnya!")
	}

	if p.AuthorID < 1 {
		return errors.New("Isi Authornya!")
	}

	return nil
}

func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error

	err = db.Debug().Model(&Post{}).Create(&p).Error

	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) AllPosts(db *gorm.DB) (*[]Post, error) {
	var err error

	posts := []Post{}

	err = db.Debug().Model(&Post{}).Find(&posts).Error

	if err != nil {
		return &[]Post{}, err
	}

	if len(posts) > 0 {
		for i, _ := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}

	return &posts, nil
}

func (p *Post) PostById(db *gorm.DB, pid uint64) (*Post, error) {
	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error

	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) UpdatePost(db *gorm.DB) (*Post, error) {
	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(
		Post{Title: p.Title, Content: p.Content, UpdateAt: time.Now()}).Error

	if err != nil {
		return &Post{}, err
	}

	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}

	return p, nil
}

func (p *Post) DeletePost(db *gorm.DB, uid uint32, pid uint64) (int64, error) {

	db = db.Debug().Model(&Post{}).Where("id = ? AND author_id = ?", pid, uid).Take(&Post{}).Delete(&Post{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post Tidak ada!")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}
