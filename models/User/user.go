package User

import (
	"flash/config"
	"flash/entities"
	"log"
)

func GetAll() []entities.Users {
	rows, err := config.DB.Query("SELECT id,username,password,bio,level,badge FROM users")
	if err != nil {
		log.Println("Query Error :", err)
		return nil
	}

	defer rows.Close()
	var users []entities.Users
	for rows.Next() {
		var user entities.Users
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Bio, &user.Level, &user.Badge); err != nil {
			log.Println("Query Error :", err)
			return nil
		}
		users = append(users, user)
	}
	return users
}
func AddData(uses entities.Users) bool {
	res, err := config.DB.Exec(`
  INSERT INTO users (username,password,bio,level,badge) VALUES(?,?,?,?,?)`,
		uses.Username, uses.Password, uses.Bio, uses.Level, uses.Badge,
	)
	if err != nil {
		log.Println("Gagal Menyimpan")
		return false
	}
	inser, erro := res.LastInsertId()
	if erro != nil {
		log.Println("Gagal Menyimpan")
	}
	return inser > 0
}
func Regis(used entities.Users) bool {
	res, err := config.DB.Exec("INSERT INTO users (username,password,bio,level,badge,email) VALUES(?,?,?,?,?,?)", used.Username, used.Password, used.Bio, used.Level, used.Badge, used.Email)
	if err != nil {
		log.Println("Error Query :", err)
		return false
	}
	sert, erno := res.LastInsertId()
	if erno != nil {
		log.Println("Payah:", err)
	}
	return sert > 0
}
func Login(username, password string) (entities.Users, bool) {
	var nilai entities.Users

	err := config.DB.QueryRow("SELECT id,username,password FROM users WHERE username = ?", username).Scan(&nilai.ID, &nilai.Username, &nilai.Password)
	if err != nil {
		log.Println("Query Bobrok")
	}
	if password != nilai.Password {
		log.Println("Password Tidak Identik")
		return nilai, false
	}
	nilai.Password = ""
	return nilai, true
}
func Update(id int, sue entities.Users) bool {
	succ, err := config.DB.Exec("UPDATE users SET username = ?,password = ? WHERE id = ? ", sue.Username, sue.Password, id)
	if err != nil {
		log.Println("Error")
	}
	affect, erno := succ.RowsAffected()
	if erno != nil {
		log.Println("Error/Gagal Update", erno)
	}
	return affect > 0
}
func Delete(id int) bool {
	dika, _ := config.DB.Exec("DELETE FROM users WHERE id = ?", id)
	rows, err := dika.RowsAffected()
	if err != nil {
		log.Println("Gagal Delete")
	}
	return rows > 0
}
