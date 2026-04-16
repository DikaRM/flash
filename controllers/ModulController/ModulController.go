package ModulController 
import (
  "net/http"
  "log"
  "html/template"
  "flash/models/Modul"
  "flash/entities"
  "io"
  "os"
  "path/filepath"
  "strconv"
  )
func Index(w http.ResponseWriter, r *http.Request){
  tmp,_ := template.ParseFiles("views/admin/module/index.html")
  moduls := Modul.GetAll()
  data := map[string]any{
    "moduls" : moduls,
  }
  log.Println("Template Html ada!")
  tmp.Execute(w,data)
}
func Add(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        // Parse multipart form dengan limit (max 10MB)
        err := r.ParseMultipartForm(10 << 20) // 10 MB
        if err != nil {
            log.Println("Error parsing form:", err)
            http.Error(w, "File terlalu besar", http.StatusBadRequest)
            return
        }
        
        // Proses Thumbnail
        var thumbnailPath string
        thumbnail, thumbHandler, err := r.FormFile("thumbnail")
        if err == nil {
            defer thumbnail.Close()
            
            // Buat direktori
            os.MkdirAll("public/uploads/thumbnail", 0755)
            
            // Simpan file
            thumbnailPath = filepath.Join("public/uploads/thumbnail", thumbHandler.Filename)
            dst, err := os.Create(thumbnailPath)
            if err != nil {
                log.Println("Error saving thumbnail:", err)
                http.Error(w, "Gagal upload thumbnail", http.StatusInternalServerError)
                return
            }
            defer dst.Close()
            io.Copy(dst, thumbnail)
        }
        
        // Proses File
        var filePath string
        file, fileHandler, err := r.FormFile("file")
        if err != nil {
            log.Println("Error getting file:", err)
            http.Error(w, "File wajib diupload", http.StatusBadRequest)
            return
        }
        defer file.Close()
        
        // Buat direktori
        os.MkdirAll("public/uploads/file", 0755)
        
        // Simpan file
        filePath = filepath.Join("public/uploads/file", fileHandler.Filename)
        dst, err := os.Create(filePath)
        if err != nil {
            log.Println("Error saving file:", err)
            http.Error(w, "Gagal upload file", http.StatusInternalServerError)
            return
        }
        defer dst.Close()
        io.Copy(dst, file)
        
        // Simpan ke database
        var modul entities.Modules
        modul.Nama_Modules = r.FormValue("nama_modul")
        modul.Deskripsi = r.FormValue("deskripsi")
        modul.Badge = r.FormValue("badge")
        modul.Kategori = r.FormValue("kategori")
        modul.File = filePath          // Path file yang sudah disimpan
        modul.Thumbnail = thumbnailPath // Path thumbnail yang sudah disimpan
        modul.Author = "dika"
        modul.View = 0
        
        if ok := Modul.Create(modul); !ok {
            log.Println("Gagal simpan ke database!")
            // Hapus file jika gagal
            os.Remove(thumbnailPath)
            os.Remove(filePath)
            http.Error(w, "Gagal menyimpan data", http.StatusInternalServerError)
            return
        }
        
        log.Println("Berhasil menyimpan modul!")
        http.Redirect(w, r, "/admin/modul", http.StatusSeeOther)
        return
    }
    
    // Method GET - tampilkan form
    // ... kode untuk menampilkan template
}
func Delete(w http.ResponseWriter,r *http.Request){
  if r.Method == "POST"{
    idr := r.PathValue("id")
    idn,_ := strconv.Atoi(idr)
    if ok := Modul.Delete(idn);!ok{
      log.Println("Errr")
      return
    }
      log.Println("Berhasil Hapus")
      http.Redirect(w,r,"/admin/modul",http.StatusSeeOther)
      return
  }
  
}