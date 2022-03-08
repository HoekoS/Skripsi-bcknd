package function

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dgrijalva/jwt-go"
	"skripsi.com/backend-adm/uuid_generator"
)

var SECRET_KEY = []byte("gosecretkey")

func GetId(r *http.Request) (string, error) {
	var err error
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		log.Fatal("Url Param 'key' is missing")
		return "error", err
	}
	key := keys[0]
	return key, nil
}

func FormatTime() string {
	var now = time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return formatted
}

func Cors(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func UploadImage(uploadedFile *multipart.File, namefile string, path string, w http.ResponseWriter) (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("error getWd : ", err)
		return "", err
	}

	randId := uuid_generator.UuidGenerate()
	filename := path + randId.String() + filepath.Ext(namefile)
	fileLocation := filepath.Join(dir, filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error di tergetfile", err)
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, *uploadedFile); err != nil {
		return "", err
	}
	return filename, nil
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}
