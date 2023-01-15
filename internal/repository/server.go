package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/cmd/internal/database"
	"github.com/cmd/internal/database/datanotes"
	"github.com/cmd/internal/database/storage"
	"github.com/cmd/internal/entities"
	"github.com/cmd/internal/forms"
	"github.com/cmd/internal/repository/services"
	"github.com/cmd/internal/utils"
	"github.com/pkg/errors"
)

func PageLogin(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "page/login" {
		http.Error(w, "Invalid URL Address", http.StatusRequestURITooLong)
	}

	file, err := os.Open("./templates/auth/login.html")
	CheckError(err, "Failed to open file")

	read, err := ioutil.ReadAll(file)
	CheckError(err, "Failed to read file")

	defer file.Close()

	switch r.Method {
	case "GET":
		w.Write(read)

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm()err:%v", err)
			return
		}

		datastorage := entities.DataUser{
			Password: r.FormValue("password"),
			UserName: r.FormValue("username"),
			Email:    r.FormValue("address"),
		}

		if forms.IsEmail(datastorage.Email) && forms.IsPassword(datastorage.Password) {

			if usecase.WebsiteAccess(&datastorage) {
				http.Redirect(w, r, "/welcome/view", http.StatusSeeOther)
			} else {
				http.Error(w, "Login details are incorrect", http.StatusUnauthorized)
			}

		} else {
			fmt.Fprintf(w, "Invalid input data!")
		}

	default:
		http.Redirect(w, r, "/page/error", http.StatusNotFound)
	}

}

func PageRegistration(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "page/registration" {
		http.Error(w, "Invalid URL Address", http.StatusRequestURITooLong)
	}

	fl, err := os.Open("./templates/auth/signup.html")
	CheckError(err, "Failed to open file")

	read, err := ioutil.ReadAll(fl)
	CheckError(err, "Failed to read file")

	defer fl.Close()

	switch r.Method {
	case "GET":
		w.Write(read)

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm()err:%v", err)
			return
		}

		datastorage := entities.DataUser{
			Password: r.FormValue("password"),
			UserName: r.FormValue("username"),
			Email:    r.FormValue("address"),
		}

		if forms.IsEmail(datastorage.Email) && forms.IsPassword(datastorage.Password) &&
			forms.IsUsername(datastorage.UserName) {

			if usecase.ExistsUser(&datastorage) {
				fmt.Fprintf(w, "User with so email exists already")
			} else {
				db, err := utils.ConnectDB()
				if err != nil {
					log.Fatal("Failed to connect db ")
				}

				storage.InsertDB(db, &datastorage)
				http.Redirect(w, r, "/login/view", http.StatusFound)
			}
		} else {
			http.Error(w, "Login details are incorrect", http.StatusUnauthorized)
		}
	default:
		http.Redirect(w, r, "/page/error", http.StatusNotFound)
	}
}

func PageResetPassword(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "page/reset/password" {
		http.Error(w, "Invalid URL Address", http.StatusRequestURITooLong)
	}

	fl, err := os.Open("./templates/recovery.html")
	CheckError(err, "Failed to open file")

	defer fl.Close()

	read, err := ioutil.ReadAll(fl)
	CheckError(err, "Failed to read file")

	switch r.Method {
	case "GET":
		w.Write(read)

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm()err:%v", err)
			return
		}

		datastorage := entities.DataUser{
			Password: r.FormValue("password"), // get new password from website
			Email:    r.FormValue("addrress"),
		}

		if forms.IsEmail(datastorage.Email) && forms.IsPassword(datastorage.Password) {

			if usecase.ExistsUser(&datastorage) {
				usecase.ChangePassword(&datastorage)
			} else {
				fmt.Fprintf(w, "User with so email doesn't exist")
			}
		} else {
			http.Error(w, "Login details are incorrect", http.StatusUnauthorized)
		}

	default:
		http.Redirect(w, r, "/page/error", http.StatusNotFound)
	}
}

func PageMain(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "page/main" {
		http.Error(w, "Invalid URL Address", http.StatusRequestURITooLong)
	}

	fl, err := os.Open("./templates/home.html")
	CheckError(err, "Failed to open db")

	defer fl.Close()

	read, err := ioutil.ReadAll(fl)
	CheckError(err, "Failed to read db")

	switch r.Method {

	case "GET":
		w.Write(read)

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm()err:%v", err)
			return
		}

		note := entities.Notes{
			Note: r.FormValue("note"),
			ID:   r.FormValue("id"),
		}
		db, err := utils.ConnectDB()

		if err != nil {
			return
		}

		notesdb.InsertNoteDB(db, &note)

	default:
		http.Redirect(w, r, "/page/error", http.StatusNotFound)
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/page/error" {
		http.Error(w, "Invalid URL Address", http.StatusRequestURITooLong)
	}

	title := r.URL.Path[len("page/error"):]
	p, err := services.LoadPage(title)

	if err != nil {
		p = &entities.Page{Title: title}
	}

	services.RenderTemplate(w, "./templates/errorpage", p)
}

func CheckError(err error, msg string) {
	if err != nil {
		errors.Wrap(err, msg)
	}
}
