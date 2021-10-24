package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AaronRebel09/golang-deployment-pipeline/database"
	helper "github.com/AaronRebel09/golang-deployment-pipeline/helpers"
	"github.com/AaronRebel09/golang-deployment-pipeline/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func fileLog() *os.File {
	// log to custom file
	LOG_FILE := "/tmp/app_log"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	return logFile
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panicln(err)
	}
	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the password in the DB
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	return check, msg
}

// CreateUser is the api used to get a single user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

		logFile := fileLog()
		defer logFile.Close()
		// Set log out put and enjoy :)
		log.SetOutput(logFile)
		// optional: log date-time, filename, and line number
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("endpoint SignUp ...")


		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("endpoint SignUp error BindJSON(&user) ...", err.Error())
			return
		}
		log.Println("endpoint success SignUp error BindJSON(&user) ...")
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			log.Println("endpoint SignUp error StatusBadRequest(validationErr) ...", validationErr.Error())
			return
		}
		log.Println("endpoint SignUp success StatusBadRequest(validationErr) ...")
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the email"})
			log.Println("endpoint SignUp error CountDocuments(user.Email) ...", err.Error())
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the phone number"})
			log.Println("endpoint SignUp error CountDocuments(user.Phone) ...", err.Error())
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// Login is the api used to get a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password is incorrect"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name,
			*foundUser.Last_name, foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		c.JSON(http.StatusOK, foundUser)
	}
}

func uploadFile(session *session.Session, uploadFileDir string) (error,string) {
	var response string
	upFile, err := os.Open(uploadFileDir)
	if err != nil {
		return nil, ""
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	result, errS3 := s3manager.NewUploader(session).Upload(&s3manager.UploadInput{
		Bucket:               aws.String("gopher-files"),
		Key:                  aws.String(uploadFileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		//ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if errS3 != nil {
		fmt.Sprintf("Failed to upload file, %v", errS3)
	} else {
		fmt.Sprintf("File uploaded to, %s", result.Location)
		response = result.Location
	}

	return errS3,response
}

func UploadImages() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* LOG FILE BEGIN */
		logFile := fileLog()
		defer logFile.Close()
		// Set log out put and enjoy :)
		log.SetOutput(logFile)
		// optional: log date-time, filename, and line number
		log.SetFlags(log.Lshortfile | log.LstdFlags)
		log.Println("endpoint UploadImages ...")
		/* LOG FILE END */

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var foundUser models.User
		var email = c.Request.PostFormValue("email")
		fmt.Println("email: ",email)

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)

		mySlice := []models.Files{}
		for _, x := range foundUser.Files {
			mySlice = append(
				mySlice,
				models.Files{
					File_name: 	x.File_name,
					File_size: x.File_size,
				})
		}

		err1 := c.Request.ParseMultipartForm(32 << 20)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Access the name key - Third Approach
		fmt.Println(c.Request.MultipartForm.File["file"])

		session, err := session.NewSession(
			&aws.Config{
				Region: aws.String("us-east-1"),
				Credentials: credentials.NewStaticCredentials(
					"AKIASYQU2445NXASRYOI", // AccessKeyID
					"ND8ngXOB6SdtIA7ytp1MPSNzTOu8IfoVdvf9KBfA", // SecretAccessKey
				"", // a token will be created when the session it's used.
				),
			})
		if err != nil {
			log.Fatal(err)
		}
		for _, h := range c.Request.MultipartForm.File["file"] {
			file, _ := h.Open()
			tmpfile, _ := os.Create("./" + h.Filename)
			io.Copy(tmpfile, file)
			fmt.Println("Archivo: ", h.Filename)
			fmt.Println("Peso: ", h.Size)

			// Upload Files
			err, responseUploadFile := uploadFile(session, h.Filename)
			if err != nil {
				log.Fatal(err)
			}

			mySlice = append(
				mySlice,
				models.Files{
				File_name: h.Filename,
				File_size: h.Size,
				File_path: responseUploadFile,
			})
		}

		foundUser.Files  = mySlice
		filter := bson.D{{"_id", foundUser.ID}}
		update := bson.D{{"$set", bson.D{{"files", foundUser.Files}}}}

		result, updateErr := userCollection.UpdateOne(ctx,filter,update)
		if updateErr != nil {
			log.Fatal(updateErr)
		}
		if result.MatchedCount != 0 {
				fmt.Println("matched and replaced an existing document")
				return
		}
		if result.UpsertedCount != 0 {
			fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		}
		defer cancel()

		c.JSON(http.StatusOK, "OK")
	}
}

func GetDocsByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret": "super secret information!",
		})
	}
}