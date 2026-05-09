package main

import (
	"flash/config"
	"flash/controllers/AuthGoogleController"
	"flash/controllers/ModulController"
	"flash/controllers/UserController"
	"flash/entities"
	"flash/models/Modul"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	config.Init()
	cfg := config.LoadConfig()
	database := config.Conn()
	fmt.Println("Flash Here we Go")
	http.Handle("/public/",
		http.StripPrefix("/public/",
			http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", UserController.Land)

	http.HandleFunc("/admin/", UserController.AdminInd)
	http.HandleFunc("/admin/modul", ModulController.Index)
	http.HandleFunc("/admin/modul-create", ModulController.Add)
	CtrlAuth := AuthGoogleController.NewAuthController(cfg, database)

	http.HandleFunc("/login", UserController.Login)
	http.HandleFunc("/auth/google/login", CtrlAuth.GoogleLogin)
	http.HandleFunc("/auth/google/callback", CtrlAuth.GoogleCallback)
	http.HandleFunc("/admin/users/update/{id}", UserController.Update)
	http.HandleFunc("/admin/users/delete/{id}", UserController.Delete)
	http.HandleFunc("/admin/modul/{id}", ModulController.Delete)
	http.HandleFunc("/content/{id}", cont)
	http.HandleFunc("/dashboard", Midle(dashboard))
	http.HandleFunc("/profile", Midle(makes))
	http.HandleFunc("/profile/mod-create", SaveMod)
	http.HandleFunc("/register", UserController.Register)
	http.HandleFunc("/admin/add-users", UserController.Add)
	http.HandleFunc("/logout", UserController.Logout)
	http.ListenAndServe(":8000", nil)
}
func makes(w http.ResponseWriter, r *http.Request) {
	// Tambahkan debugging
	log.Println("=== START makes handler ===")

	// 1. Cek session
	session, err := config.GetSession(r)
	if err != nil {
		log.Println("Session error:", err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}
	log.Println("Session OK")

	// 2. Dapatkan username
	var username string
	val, ok := session.Values["username"].(string)
	if !ok {
		log.Println("Username not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	username = val
	log.Println("Username:", username)

	// 3. Dapatkan data user
	users, err := Modul.Prof(username)
	if err != nil {
		log.Println("Error getting profile:", err)
		// Jangan return, tetap lanjut dengan users kosong
		users = entities.Users{} // Gunakan struct kosong jika ada error
	}
	log.Printf("User data: %+v", users)

	// 4. Dapatkan data moduls
	moduls := Modul.Mod(username)
	log.Printf("Moduls count: %d", len(moduls))
	if moduls == nil {
		moduls = []entities.Modules{} // Gunakan slice kosong jika nil
	}

	// 5. Siapkan data untuk template
	data := map[string]interface{}{
		"user":   users,
		"moduls": moduls,
	}
	log.Printf("Data prepared: user=%+v, moduls count=%d", users, len(moduls))

	// 6. Parse template dengan pengecekan path absolut
	templatePath := "views/profil.html"

	// Cek apakah file benar-benar ada
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Template file not found at: %s", templatePath)

		// Coba cari di current directory
		dir, _ := os.Getwd()
		log.Printf("Current working directory: %s", dir)

		// Coba cek path alternatif
		altPath := filepath.Join(dir, "views", "profil.html")
		log.Printf("Checking alternative path: %s", altPath)

		if _, err := os.Stat(altPath); err == nil {
			templatePath = altPath
			log.Println("Found template at alternative path")
		} else {
			http.Error(w, "Template not found", http.StatusInternalServerError)
			return
		}
	}

	// Parse template
	tmp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("ERROR parsing template: %v", err)
		log.Printf("Template path: %s", templatePath)
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("Template parsed successfully")

	// 7. Execute template
	err = tmp.Execute(w, data)
	if err != nil {
		log.Printf("ERROR executing template: %v", err)
		http.Error(w, fmt.Sprintf("Execution error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("=== END makes handler ===")
}
func Midle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := config.GetSession(r)
		if err != nil {
			log.Println("error Session", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		loged, ok := session.Values["logged_in"].(bool)
		if !ok || !loged {
			session.AddFlash("Kudu Login", "error")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	moduls := Modul.GetAll()
	session, err := config.GetSession(r)
	if err != nil {
		log.Println("Gagal Dapat Session")
	}
	val := session.Values["username"]
	username, ok := val.(string)
	if !ok {
		username = "Guest"
	}
	data := map[string]any{
		"moduls":   moduls,
		"username": username,
	}
	tmpl, _ := template.ParseFiles("views/dashboard.html")
	tmpl.Execute(w, data)
}
func SaveMod(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse form dengan limit size (misal 10MB)
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "File terlalu besar", http.StatusBadRequest)
			return
		}

		// Proses Thumbnail (opsional)
		var dstpath string
		file, handler, err := r.FormFile("thumbnail")
		if err == nil {
			defer file.Close()

			// Buat direktori jika belum ada
			uploadDir := "public/uploads/thumbnail"
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				log.Println("Error creating directory:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			dstpath = filepath.Join(uploadDir, handler.Filename)
			dst, err := os.Create(dstpath)
			if err != nil {
				log.Println("Error creating file:", err)
				http.Error(w, "Gagal upload thumbnail", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				log.Println("Error copying file:", err)
				http.Error(w, "Gagal upload thumbnail", http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("Tidak Ada Thumbnail:", err)
			// Jika thumbnail wajib, return error
			// http.Error(w, "Thumbnail wajib diupload", http.StatusBadRequest)
			// return
		}

		// Proses File (wajib)
		files, dler, err := r.FormFile("file")
		if err != nil {
			log.Println("Tidak Ada files:", err)
			http.Error(w, "File wajib diupload", http.StatusBadRequest)
			// Hapus thumbnail jika sudah terupload
			if dstpath != "" {
				os.Remove(dstpath)
			}
			return
		}
		defer files.Close()

		// Buat direktori untuk file
		fileUploadDir := "public/uploads/file"
		if err := os.MkdirAll(fileUploadDir, 0755); err != nil {
			log.Println("Error creating directory:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		stp := filepath.Join(fileUploadDir, dler.Filename)
		st, err := os.Create(stp)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(w, "Gagal upload file", http.StatusInternalServerError)
			if dstpath != "" {
				os.Remove(dstpath)
			}
			return
		}
		defer st.Close()

		_, err = io.Copy(st, files)
		if err != nil {
			log.Println("Error copying file:", err)
			http.Error(w, "Gagal upload file", http.StatusInternalServerError)
			if dstpath != "" {
				os.Remove(dstpath)
			}
			os.Remove(stp)
			return
		}

		// Simpan ke database
		var modul entities.Modules
		modul.Nama_Modules = r.FormValue("nama_module")
		modul.Deskripsi = r.FormValue("deskripsi")
		modul.Badge = r.FormValue("badge")
		modul.Kategori = r.FormValue("kategori")
		modul.File = stp // Simpan path relatif
		modul.Thumbnail = dstpath
		modul.Author = r.FormValue("author")
		modul.View = 0

		// Panggil fungsi Create dengan context
		if ok := Modul.Create(modul); !ok {
			log.Println("Gagal menyimpan ke database!")
			// Hapus file jika gagal simpan ke database
			if dstpath != "" {
				os.Remove(dstpath)
			}
			os.Remove(stp)
			http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
			return
		}

		log.Println("Berhasil menyimpan modul!")

		// Set header dan redirect
		w.Header().Set("Content-Type", "application/json")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// Jika method bukan POST
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
func cont(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

	}
	if r.Method == "GET" {
		idstr := r.PathValue("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		modul, err := Modul.GetId(id)
		if err != nil {
			http.Error(w, "Module not found", http.StatusNotFound)
			return
		}

		ext := filepath.Ext(modul.File)
		pdfURL := modul.File

		// Sesuaikan dengan static server
		// Jika static server di /public/, dan file ada di public/uploads/file/...

		// Debug
		log.Printf("File from DB: %s", modul.File)
		log.Printf("PDF URL: %s", pdfURL)
		// Tambah slash di depan jika tidak ada

		if !strings.HasPrefix(modul.File, "/") {
			pdfURL = "/" + modul.File
		} else {
			pdfURL = modul.File
		}
		data := map[string]any{
			"modul":     modul,
			"extension": ext,
			"pdfURL":    pdfURL,
		}

		tmp, err := template.ParseFiles("views/content.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			log.Printf("Template error: %v", err)
			return
		}

		err = tmp.Execute(w, data)
		if err != nil {
			http.Error(w, "Render error", http.StatusInternalServerError)
			log.Printf("Render error: %v", err)
			return
		}
	}
	if r.Method == "POST" {
		idstr := r.PathValue("id")
		id, _ := strconv.Atoi(idstr)
		if ok := Modul.DateView(id); !ok {
			log.Println("Data Gagal Query Viww")
		}
		http.Redirect(w, r, "/content/"+idstr, http.StatusSeeOther)
		return
	}
}
