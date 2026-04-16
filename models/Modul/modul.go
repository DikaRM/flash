package Modul
import (
  "flash/entities"
  "flash/config"
  "log"
  "database/sql"
  "fmt"
  )
func GetAll()[]entities.Modules{
  query,_ := config.DB.Query("SELECT id,nama_modules,deskripsi,file,badge,total_rating,avg_rating,author,view,thumbnail FROM modules")
  defer query.Close()
  var moduls []entities.Modules
  for query.Next(){
    var modul entities.Modules
    if err := query.Scan(&modul.ID,&modul.Nama_Modules,&modul.Deskripsi,&modul.File,&modul.Badge,&modul.TotalRating,&modul.RatingAve,&modul.Author,&modul.View,&modul.Thumbnail,);err != nil {
      log.Println("Gagal Query")
      return nil
    }
    moduls = append(moduls,modul)
  }
  return moduls
}
func Prof(nama string)(entities.Users,error){
  row:= config.DB.QueryRow("SELECT id,username,password,bio,level,badge FROM users WHERE username = ? LIMIT 1",nama)
  
  var user entities.Users
  
  err := row.Scan(&user.ID,&user.Username,&user.Password,&user.Bio,&user.Level,&user.Badge,)
  if err != nil{
    
    if err == sql.ErrNoRows{
      return user,fmt.Errorf("Error Ga ada Data")
    }
    return user,err
  }
  return user,nil
}
func Mod(username string) []entities.Modules {
    if config.DB == nil {
        log.Println("Database connection is nil")
        return []entities.Modules{}
    }
    
    rows, err := config.DB.Query("SELECT id,nama_modules,deskripsi,file,badge,total_rating,avg_rating,author,view,thumbnail,kategori,created_at FROM modules WHERE author = ?", username)
    if err != nil {
        log.Println("Gagal Query:", err)
        return []entities.Modules{}
    }
    defer rows.Close()
    
    var modu []entities.Modules
    for rows.Next() {
        var modul entities.Modules
        // Perbaiki urutan Scan sesuai dengan kolom di database
        err := rows.Scan(
            &modul.ID,
            &modul.Nama_Modules,
            &modul.Deskripsi,
            &modul.File,
            &modul.Badge,
            &modul.TotalRating,
            &modul.RatingAve,
            &modul.Author,
            &modul.View,
            &modul.Thumbnail,
            &modul.Kategori,
            &modul.Created_at,
        )
        if err != nil {
            log.Println("Gagal Scan:", err)
            continue // Skip data yang error, jangan return langsung
        }
        modu = append(modu, modul)
    }
    
    // Cek error dari rows
    if err = rows.Err(); err != nil {
        log.Println("Error after rows iteration:", err)
    }
    
    return modu
}
func Create(modul entities.Modules)bool{
  query,er:= config.DB.Exec("INSERT INTO modules (nama_modules,deskripsi,file,badge,author,view,thumbnail,kategori)VALUES(?,?,?,?,?,?,?,?)",modul.Nama_Modules,modul.Deskripsi,modul.File,modul.Badge,modul.Author,modul.View,modul.Thumbnail,modul.Kategori,)
  if er != nil{
    log.Println("Error ",er)
    return false
  }
  lat,err := query.RowsAffected()
  if err != nil{
    log.Println("Gagal Gagal")
    return false
  }
  return lat > 0
}
func AddRating(moduleId,userId ,rating int)bool{
  var exit entities.Ratings
  err := config.DB.QueryRow("SELECT id,rating FROM ratings WHERE module_id = ? AND user_id = ?",moduleId,userId).Scan(&exit.ID,&exit.Rating)
  if err == nil {
    _,err := config.DB.Exec("UPDATE ratings SET rating = ? WHERE id = ? ",rating,exit.ID)
    if err != nil{
      log.Println("Gagal Gagal gagal")
      return false
    }
  }else{
    _,err := config.DB.Exec("INSERT INTO ratings  (module_id,user_id,rating) VALUES(?,?,?)",moduleId,userId,rating)
    if err != nil{
      log.Println("Gagal Gagal Insert ",err)
      return false
    }
  }
  updateModuleRating(moduleId)
  return true
}
func updateModuleRating(moduleId int){
  var avg float64
  var total int
  err := config.DB.QueryRow("SELECT COALESCE(AVG(rating),0),COUNT(*) FROM ratings WHERE modules =? ",moduleId).Scan(&avg,&total)
  if err == nil {
    config.DB.Exec("UPDATE modules SET avg_rating = ? ,total_rating = ? WHERE id =?",avg,total,moduleId)
    
  }
}
func GetComment(moduleId int)[]entities.Comments{
  query,_ := config.DB.Query(`
    SELECT c.id, c.module_id, c.user_id, c.comment, c.created_at, u.username 
    FROM comments c 
    JOIN users u ON c.user_id = u.id 
    WHERE c.module_id = ? 
    ORDER BY c.created_at DESC
`, moduleId)
  defer query.Close()
  var comments []entities.Comments
  for query.Next(){
    var comment entities.Comments
    query.Scan(&comment.ID,&comment.ModuleId,&comment.UserId,&comment.Comment,&comment.Created_at,&comment.Username)
    comments = append(comments,comment)
  }
  return comments
}
func AddComment(ModuleId,UserId int,comment string)bool{
  
  query,_ := config.DB.Exec("INSERT INTO comments (module_id,user_id,comment)VALUES(?,?,?)",ModuleId,UserId,comment,)
  lastId,err := query.LastInsertId()
  if err != nil{
    log.Println("Gagal Menambahkan Komentar")
    return false
  }
  return lastId > 0
}
func Delete(idr int)bool{
  query,_ := config.DB.Exec("DELETE FROM modules WHERE id = ?",idr)
  affect,_ := query.RowsAffected()
  return affect > 0
}
func GetId(id int)(*entities.Modules,error){
  var m entities.Modules
  err := config.DB.QueryRow("SELECT id,nama_modules,deskripsi,file,badge,author,view,kategori,created_at FROM modules WHERE id = ? ",id).Scan(&m.ID,&m.Nama_Modules,&m.Deskripsi,&m.File,&m.Badge,&m.Author,&m.View,&m.Kategori,&m.Created_at,)
  if err != nil{
    log.Println("error Diks")
    return nil,err
  }
  return &m, nil
}
func DateView(id int)bool{
  var mod entities.Modules
  err := config.DB.QueryRow("SELECT view FROM modules WHERE id = ?",id).Scan(&mod.View,)
  if err != nil{
    log.Println("Gagal Query View")
  }
  views := mod.View + 1
  query,_ := config.DB.Exec("UPDATE modules SET view = ? WHERE id = ?",views,id)
  row,_ := query.RowsAffected()
  return row > 0
}