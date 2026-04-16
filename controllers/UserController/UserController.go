package UserController

import (
	"flash/config"
	"flash/entities"
	"flash/models/User"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func Land(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("views/index.html")
	tmp.Execute(w, nil)
}
func AdminInd(w http.ResponseWriter, r *http.Request) {
	users := User.GetAll()
	data := map[string]any{
		"users": users,
	}
	tmp, err := template.ParseFiles("views/admin/users/index.html")
	if err != nil {
		http.Error(w, "Template ga ada"+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmp.Execute(w, data)
	if err != nil {
		http.Error(w, "Data Error ga ada"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var uses entities.Users
		uses.Username = r.FormValue("username")
		uses.Password = r.FormValue("password")
		uses.Bio = "i am programmer"
		uses.Level = 1
		uses.Badge = "stator"
		if ok := User.AddData(uses); !ok {
			log.Println("Eror")
			http.Error(w, "Gagal maning", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmp, err := template.ParseFiles("views/login.html")
		if err != nil {
			log.Println("Error Template :", err)
		}
		tmp.Execute(w, nil)
	}
	if r.Method == "POST" {
		nama := r.FormValue("nama")
		passwd := r.FormValue("password")

		if _, okk := User.Login(nama, passwd); !okk {
			session, _ := config.GetSession(r)
			session.AddFlash("Usernam Gagal")
			session.Save(r, w)
			log.Println("Gagal")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		session, err := config.GetSession(r)
		if err != nil {
			log.Println("Gagal Session !")
			return
		}
		session.Values["username"] = nama
		session.Values["logged_in"] = true
		session.AddFlash("Username Sudah Ditambahkan Ke session")
		err = session.Save(r, w)
		if err != nil {
			log.Println("Gagal Session !")
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmp, _ := template.ParseFiles("views/register.html")

		tmp.Execute(w, nil)
	}
	if r.Method == "POST" {
		var used entities.Users
		used.Username = r.FormValue("username")
		used.Password = r.FormValue("password")
		used.Email = r.FormValue("email")
		used.Bio = "I am Programmer"
		used.Level = 1
		used.Badge = "stator"
		if right := User.Regis(used); !right {
			log.Println("Nih errornya")
			http.Error(w, "Errror Yeuh", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var sue entities.Users
		idstr := r.PathValue("id")
		Idn, _ := strconv.Atoi(idstr)
		sue.Username = r.FormValue("username")
		sue.Password = r.FormValue("Password")
		if nows := User.Update(Idn, sue); !nows {
			log.Println("gagal")
			http.Error(w, "Gagal Total", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

}
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idstr := r.PathValue("id")
		ids, err := strconv.Atoi(idstr)
		if err != nil {
			log.Println("Errror", err)
		}
		if des := User.Delete(ids); !des {
			log.Println("gagal")
			http.Error(w, "Gagal Total", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session, _ := config.GetSession(r)
		session.Values = make(map[interface{}]interface{})
		session.Options.MaxAge = -1
		session.AddFlash("Anda Telah Logout")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
