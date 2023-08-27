package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tvvat/project/internal/model"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

type PostgresDB struct {
	client *gorm.DB
}

func NewPostgresDB(cfg PostgresConfig) (*PostgresDB, error) {
	pgdb, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.DBName,
			cfg.Password,
		),
	)
	if err != nil {
		log.Printf("can't connect to DB: %s", err)
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pgdb,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	return &PostgresDB{client: db}, nil
}

func (db *PostgresDB) Migrate() {
	db.client.AutoMigrate(&model.User{})
	db.client.AutoMigrate(&model.Group{})
}

func (db *PostgresDB) GetAllUsers() []model.User {
	var users []model.User
	db.client.Find(&users)
	return users
}

func (db *PostgresDB) GetAllGroups() []model.Group {
	var groups []model.Group
	db.client.Find(&groups)
	return groups
}

func (db *PostgresDB) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := db.client.Where("Id = ?", id).First(&user).Error
	return user, err
}

func (db *PostgresDB) GetGroupByID(id int) (model.Group, error) {
	var group model.Group
	err := db.client.Where("Id = ?", id).First(&group).Error
	return group, err
}

func (db *PostgresDB) DeleteUserByID(id int) error {
	var user model.User
	err := db.client.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	db.client.Delete(&user)
	log.Printf("user deleted successfully")

	return nil
}

func (db *PostgresDB) DeleteGroupByID(id int) error {
	var group model.Group
	err := db.client.Where("id = ?", id).First(&group).Error
	if err != nil {
		return err
	}

	db.client.Delete(&group)
	log.Printf("group deleted successfully")

	return nil
}

func (db *PostgresDB) CreateUser(id int, login, password string, groupIDs []int) model.User {
	var groups []model.Group
	db.client.Find(&groups, groupIDs)

	user := model.User{Id: id, Login: login, Password: password, Groups: groups}
	db.client.Create(&user)
	log.Printf("user created successfully")

	return user
}

func (db *PostgresDB) CreateGroup(id int, subgroupIDs []int) model.Group {
	var groups []model.Group
	db.client.Find(&groups, subgroupIDs)

	group := model.Group{Id: id, Subgroups: groups}
	db.client.Create(&group)
	log.Printf("group created successfully")

	return group
}

func (db *PostgresDB) UpdateUserByID(id int, login, password string) (model.User, error) {
	var user model.User
	err := db.client.Where("Id = ?", id).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	db.client.Model(user).Updates(map[string]interface{}{
		"id":       id,
		"login":    login,
		"password": password,
	})
	log.Printf("user update successfully")

	return user, nil
}
