package seeder

import (
	"fmt"
	"log"

	"github.com/Muh-Sidik/BlogSimpel-Gomux/api/models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Username: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	models.User{
		Username: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var posts = []models.Post{
	models.Post{
		Title:   "Title 1",
		Content: "Hello world 1",
	},
	models.Post{
		Title:   "Title 2",
		Content: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	var confirm string

	fmt.Println("Apakah ingin drop table sebelum mulai? yes/no")
	_, err := fmt.Scan(&confirm)
	if err != nil {
		log.Fatalln("Kesalahan Server!")
	}

	if confirm == "yes" {
		err := db.Debug().DropTableIfExists(&models.User{}, models.Post{}).Error

		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}

		err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
		if err != nil {
			log.Fatalf("cannot migrate table: %v", err)
		}

		err = db.Debug().Model(&models.Post{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}

		for i, _ := range users {
			err := db.Debug().Model(&models.User{}).Create(&users[i]).Error

			if err != nil {
				log.Fatalf("cannot seed users table: %v", err)
			}

			posts[i].AuthorID = users[i].ID

			err = db.Debug().Model(&models.Post{}).Create(&posts[i]).Error

			if err != nil {
				log.Fatalf("cannot seed posts table: %v", err)
			}
		}
	}

}
