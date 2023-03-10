package usecase

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cmd/internal/entities"
	"github.com/cmd/internal/utils"
	"log"
	"testing"
)

// just simple testing

func TestWebSiteAccess(t *testing.T) {
	personallyData := entities.DataUser{
		Password: "Ruslan12345",
		Email:    "ruslan@mail.ru",
		UserName: "Ruslan",
	}

	if WebsiteAccess(&personallyData) {
		fmt.Println("Login completed successfully")
	} else {
		t.Errorf("Site access denied")
	}
}

func TestExistsUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatal("Failed to connect database sqlmock")
	}

	defer db.Close()

	mock.ExpectExec("SELECT email FROM storage")

	conndb, err := utils.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect database storage")
	}

	defer conndb.Close()

	storage := entities.DataUser{
		Email: "ruslan@example.ru",
	}

	if !ExistsUser(&storage) {
		t.Errorf("User with so email address doesn't exist in database storage")
	}
}

func TestChangePassword(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatal("Failed to connect database sqlmock")
	}

	defer db.Close()

	mock.ExpectExec("SELECT password,email FROM storage")

	conndb, err := utils.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect database storage")
	}

	defer conndb.Close()

	storage := entities.DataUser{
		Password: "Ruslan12345",
		Email:    "ruslan@example.com",
	}

	if err = ChangePassword(&storage); err != nil {
		t.Errorf("Failed to change password,tests failed")
	}

}

func TestGetDataDB(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatal("Failed to connect sqlmock database", err)
	}
	defer db.Close()

	conndb, err := utils.ConnectDB()

	if err != nil {
		log.Fatal("Faield to connect to database storage", err)
	}

	defer conndb.Close()

	mock.ExpectExec("SELECT username,password,email FROM storage")

	storage := GetDataDB()

	for _, rng := range storage {
		if rng.Email == "" && rng.Password == "" && rng.UserName == "" {
			t.Errorf("Failed to get some data from database storage")
		}
	}
}

func TestGetByUserName(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		log.Fatal("Failed to connect database sqlmock")
	}

	defer db.Close()

	storagedb, err := utils.ConnectDB()

	if err != nil {
		log.Fatal("Failed to connect database storage")
	}

	defer storagedb.Close()

	mock.ExpectExec("SELECT username FROM storage")
	
	if _, err := GetByUserName(); err != nil {
		t.Errorf("Faile to get username:%v", err)
	}
}
