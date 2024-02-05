package models

// User representa la estructura de un usuario para importar
type User struct {
  Id       int64  `json:"id"`
  Email    string `json:"email"`
  Verified int8   `json:"verified"`
  Password string `json:"password"`
}
