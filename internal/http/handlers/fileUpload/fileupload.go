package fileupload

import (
	// "image"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	// "strconv"
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/roshanansy/file-upload/internal/database"
	"github.com/roshanansy/file-upload/internal/filter/images"
	"github.com/roshanansy/file-upload/internal/http/middleware/authorized"
	"github.com/roshanansy/file-upload/internal/model"
	"github.com/roshanansy/file-upload/internal/types"
	"github.com/roshanansy/file-upload/internal/utils/response"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation of file upload handling
	userClaims, err := authorized.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "unauthorized: missing user context", http.StatusUnauthorized)
		return
	}

	userID, ok := userClaims["ID"].(string)
	if !ok || len(userID)==0 {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}
	//slog.info("File upload handler invoked")
	
	header :=r.Header.Get("Content-Type")
	if len(header)==0{
		http.Error(w, "Content Type header missing", http.StatusBadRequest)
		return
	}
	slog.Info("Content Type Header","header",header)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
	http.Error(w, "Error parsing form data", http.StatusBadRequest)
	return
	}

	files := r.MultipartForm.File["files"]

	if len(files) ==0 {
		http.Error(w, "pease select the file", http.StatusBadRequest)
		return
	}

	// rootDir, err := os.Getwd()
	// if err != nil {
	// 	http.Error(w, "Unable to get current directory", http.StatusInternalServerError)
	// 	return
	// }
	
	uploadDir := os.Getenv("UPLOAD_DIR")
	var responses []types.FileResponse

	slog.Info("this is env file locations","uploadPath",uploadDir)

	if len(uploadDir)==0 {
		http.Error(w, "unable to find the path to store the file", http.StatusInternalServerError)
		return
	}
	// ✅ Go one level up (to project root)
	// projectRoot := filepath.Dir(rootDir)

	// uploaddir := "./../uploaddir"

	// ✅ Ensure upload directory exists
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	for _,headers := range files {
		file, err := headers.Open()
		if err != nil {
			slog.Error("Failed to open file", "error", err)
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer file.Close();

		fileUUID := uuid.New().String()
		newFileName := fileUUID + "_" + headers.Filename

			// Define the path to save the file
		filePath := filepath.Join(uploadDir, newFileName)
		// Create destination file
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		absPath, _ := filepath.Abs(filePath)
		publicURL := "/uploads/" + newFileName

		responses = append(responses, types.FileResponse{
			FileName:  headers.Filename,
			UUID:      fileUUID,
			TempPath:  headers.Filename,
			FullPath:  absPath,
			AccessURL: publicURL,
		})

		fileConfig :=model.FileStore{
			FileName: headers.Filename,
			FilePath: filePath,
			FileType: headers.Header.Get("Content-Type"),
			FileSize: headers.Size,
			UserID: userID,
		}
		database.DB.Create(&fileConfig)

	}


	response.WriteJson(w, http.StatusOK, responses)
	
}



func FilterImagesHandler(w http.ResponseWriter, r *http.Request) {
	userClaims, err := authorized.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "unauthorized: missing user context", http.StatusUnauthorized)
		return
	}

	userID, ok := userClaims["ID"].(string)
	if !ok || userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "please select a file", http.StatusBadRequest)
		return
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		http.Error(w, "UPLOAD_DIR not configured", http.StatusInternalServerError)
		return
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
		return
	}

	var responses []types.FileResponse

	for _, header := range files {
		srcFile, err := header.Open()
		if err != nil {
			http.Error(w, "Failed to open file", http.StatusInternalServerError)
			return
		}
		defer srcFile.Close()

		fileUUID := uuid.New().String()
		newFileName := fileUUID + "_" + header.Filename
		filePath := filepath.Join(uploadDir, newFileName)

		dstFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		// Apply filters
		filteredPath, err := images.FilterImages(filePath)
		if err != nil {
			http.Error(w, "Failed to filter image", http.StatusInternalServerError)
			return
		}

		publicURL := "/uploads/" + filepath.Base(filteredPath)

		responses = append(responses, types.FileResponse{
			FileName:  header.Filename,
			UUID:      fileUUID,
			TempPath:  filePath,
			FullPath:  filteredPath,
			AccessURL: publicURL,
		})

		fileConfig := model.FileStore{
			FileName: header.Filename,
			FilePath: filteredPath,
			FileType: header.Header.Get("Content-Type"),
			FileSize: header.Size,
			UserID:   userID,
		}
		database.DB.Create(&fileConfig)
	}

	response.WriteJson(w, http.StatusOK, responses)
}


func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Expect query params: ?filename=<name> or ?uuid=<uuid>
	// fileName := r.URL.Query().Get("filename")
	fileUUID := r.URL.Query().Get("uuid")

	if  fileUUID == "" {
		http.Error(w, "Missing uuid parameter", http.StatusBadRequest)
		return
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		http.Error(w, "UPLOAD_DIR environment variable not set", http.StatusInternalServerError)
		return
	}

	var targetFile string

	// Search for the file in the directory
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		http.Error(w, "Unable to read upload directory", http.StatusInternalServerError)
		return
	}

	// for _, f := range files {
	// 	// Match by UUID prefix or exact filename
	// 	if (!f.IsDir()) && (f.Name() == fileName || (fileUUID != "" && filepath.HasPrefix(f.Name(), fileUUID))) {
	// 		targetFile = filepath.Join(uploadDir, f.Name())
	// 		break
	// 	}
	// }

	
	// Find file that starts with given UUID
	for _, f := range files {
		if !f.IsDir() && filepath.HasPrefix(f.Name(), fileUUID+"_") {
			targetFile = filepath.Join(uploadDir, f.Name())
			break
		}
	}

	if targetFile == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Attempt to delete the file
	if err := os.Remove(targetFile); err != nil {
		slog.Error("Failed to delete file", "error", err)
		http.Error(w, "Unable to delete file", http.StatusInternalServerError)
		return
	}

	slog.Info("File deleted successfully", "file", targetFile)

	response.WriteJson(w, http.StatusOK, map[string]string{
		"message": "File deleted successfully",
		"file":    targetFile,
	})
}


func GenerateApiKeyHandler(w http.ResponseWriter,r*http.Request){
	userClaims, err := authorized.GetUserFromContext(r)
	if err != nil {
		http.Error(w, "unauthorized: missing user context", http.StatusUnauthorized)
		return
	}

	userID, ok := userClaims["ID"].(string)
	if !ok || userID == "" {
		http.Error(w, "invalid user ID in token", http.StatusUnauthorized)
		return
	}

	var apiKey types.GenerateApiKeyRequest

	if err :=json.NewDecoder(r.Body).Decode(&apiKey); err!=nil{
		response.WriteJson(w,http.StatusBadRequest,map[string]string{
			"error":"Invalid request payload",
		})
		return
	}
	publicKey,err:=GenerateAPIKey(apiKey.KeyName,32)
	privateKey,err:=GenerateAPIKey("private_"+apiKey.KeyName,64)
	hash, err := bcrypt.GenerateFromPassword([]byte(privateKey), bcrypt.DefaultCost)
	modelApiKey:= model.GenerateApiKey{
		UserID: userID,
		KeyName: apiKey.KeyName,
		PublicKey: publicKey,
		SecretKey: string(hash),
	}

	if err:=database.DB.
			Create(&modelApiKey).Error;err!=nil{
		response.WriteJson(w,http.StatusInternalServerError,map[string]string{
			"error":"Failed to store API key",
		})
		return
	}
	response.WriteJson(w,http.StatusOK,map[string]string{
		"publicKey":publicKey,
		"secretKey":privateKey,
	})
}

 func GenerateAPIKey(prefix string, length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	key := base64.URLEncoding.EncodeToString(bytes)
	return prefix + key, nil
}


 type  TestType struct{
	Param1 int32
	name string
	email string
 }

func TestCodeFix(  abc []int32,a int32, b int32, str string ,user TestType)int32{
	
	var d int32 =20
	var arr []int32=[]int32{}

	var userSlice []TestType=[]TestType{}
	userSlice=append(userSlice,TestType{Param1:1,name:"Roshan",email:"test"})
	userSlice=append(userSlice,TestType{Param1:2,name:"Ansy",email:"ansy@test"})

	arr=append(arr, 122,34,34,45,56)

	for i:=0;i<len(abc);i++{
		slog.Info("Array values", "index", i, "value", abc[i])
	}

	for i,v:=range userSlice{
		slog.Info("User Slice values", "index", i, "name", v.name, "email", v.email)
	}

	m:=10
	var  m_str string=strconv.Itoa(m)

		value, err :=strconv.ParseInt(m_str,10,64)
	if err != nil {
		slog.Error("Error converting string to int64", "error", err)
	}

	slog.Info("This is a test function", "param1", value)
	slog.Info("This is a test function", "param2", m_str)

	
	slog.Info("This is a test function", "param3", str)
	slog.Info("User Info", "name", user.name, "email", user.email, "param1", user.Param1)
	return a+b+d;

}

