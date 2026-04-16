package entities

type Users struct {
	ID         int64
	Username   string
	Password   string
	Bio        string
	Level      int
	Badge      string
	Email      string
	GoogleID   *string
	Avatar     string
	Created_at string
}
type Modules struct {
	ID           int
	Nama_Modules string
	Deskripsi    string
	File         string
	Badge        string
	TotalRating  int
	RatingAve    float64
	Author       string
	View         int
	Thumbnail    string
	Kategori     string
	Created_at   string
}
type Ratings struct {
	ID int
	string
	UserId   int
	ModuleId int
	Rating   int
}
type Comments struct {
	ID         int
	UserId     int
	ModuleId   int
	Comment    string
	Username   string
	Created_at string
}
