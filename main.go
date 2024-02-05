package main

import (
  "context"
  "encoding/json"
  "io/ioutil"
  "log"

  "firebase.google.com/go"
  "firebase.google.com/go/auth"
  "firebase.google.com/go/auth/hash"
  "github.com/google/uuid"
  "go_migrate-users-to-firebase/models"
  "google.golang.org/api/option"
)

func main() {
  // Cargar archivo de credenciales de Firebase
  sa := option.WithCredentialsFile("serviceAccountKey.json")
  app, err := firebase.NewApp(context.Background(), nil, sa)
  if err != nil {
    log.Fatalf("error initializing app: %v", err)
  }

  // Obtener cliente de autenticación
  ctx := context.Background()
  client, err := app.Auth(ctx)
  if err != nil {
    log.Fatalf("error getting Auth client: %v\n", err)
  }

  // Leer archivo JSON con datos de los usuarios a importar
  file, err := ioutil.ReadFile("data.json")
  if err != nil {
    log.Fatalf("error reading users file: %v", err)
  }

  // Deserializar datos de usuarios desde JSON
  var usersData []models.User
  err = json.Unmarshal(file, &usersData)
  if err != nil {
    log.Fatalf("error unmarshalling users data: %v", err)
  }

  // Convertir datos de usuarios a registros importables para Firebase
  var users []*auth.UserToImport
  for _, u := range usersData {
    // Genera un UID único para cada usuario
    uid := uuid.New().String()

    // Añade el usuario a la lista de usuarios a importar
    users = append(users, (&auth.UserToImport{}).
      UID(uid).
      Email(u.Email).
      EmailVerified(u.Verified == 1).
      PasswordHash([]byte(u.Password)).
      PasswordSalt([]byte("salt1")))
  }

  // Configurar el algoritmo de hash con el cual esta codificada
  // la contraseña del usuario en el JSON con los usuarios
  h := hash.Bcrypt{}

  // Importar usuarios a Firebase
  result, err := client.ImportUsers(ctx, users, auth.WithHash(h))
  if err != nil {
    log.Fatalln("Unrecoverable error prevented the operation from running", err)
  }

  // Verificar y reportar el resultado de la importación
  log.Printf("Successfully imported %d users\n", result.SuccessCount)
  log.Printf("Failed to import %d users\n", result.FailureCount)
}
