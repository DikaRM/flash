package config
import (
  "os"
  "net/http"
  
  "github.com/gorilla/sessions"
)
var Store *sessions.CookieStore
func Init(){
  secretKey := os.Getenv("SESSION_SECRET")
  if secretKey == ""{
    secretKey = "dika"
  }
  Store = sessions.NewCookieStore([]byte(secretKey))
  Store.Options = &sessions.Options {
    Path : "/",
    MaxAge : 86400 *7,
    HttpOnly : true,
    Secure : false,
    SameSite : http.SameSiteStrictMode,
  }
}
func GetSession(r *http.Request)(*sessions.Session,error){
  if Store == nil{
    Init()
  }
  session,err := Store.Get(r,"sess")
  if err != nil {
    return nil,err
  }
  return session,nil
}
  