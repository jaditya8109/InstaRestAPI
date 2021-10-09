package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
	. "./config"
	. "./dao"
	. "./models"
)

var config = Config{}
var dao = UsersDAO{}

// GET a users by its ID
func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

// Create a new user
func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, user)
}

// get user post
func (c *Client) GetUserMedia(userId, token string, count int) (*MediaPage, error) {
	url := fmt.Sprintf("%s/users/%s/media/recent?count=%d&access_token=%s", baseUrl, userId, count, token)
	return c.getMediaPage(url)
}

// Create Post
func (s *ImageUploaderService) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}

	db.SingularTable(true)

	db.AutoMigrate(&api.Image{})
	db.AutoMigrate(&api.ImageType{})
	db.AutoMigrate(&api.Meta{})

	db.Model(&api.Image{}).AddForeignKey("image_type_id", "image_type(id)", "RESTRICT", "RESTRICT")

	meta := api.Meta{
		Name:       "image_count",
		Value_int:  0,
		Created_at: int32(time.Now().Unix()),
		Updated_at: int32(time.Now().Unix()),
	}
	db.Create(&meta)

	imageType := api.ImageType{
		Name:         "order",
		Path:         "/order",
		Thumb_width:  320,
		Thumb_height: 240,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	imageType = api.ImageType{
		Name:         "user",
		Path:         "/user",
		Thumb_width:  150,
		Thumb_height: 150,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	imageType = api.ImageType{
		Name:         "advert",
		Path:         "/advert",
		Thumb_width:  320,
		Thumb_height: 240,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	imageType = api.ImageType{
		Name:         "image",
		Path:         "/image",
		Thumb_width:  320,
		Thumb_height: 240,
		Created_at:   int32(time.Now().Unix()),
		Updated_at:   int32(time.Now().Unix()),
	}
	db.Create(&imageType)

	return nil
}

// Get all posts of user
func (c *Client) getMediaPage(url string) (*MediaPage, error) {
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var page MediaPage
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, err
	}

	page.client = c
	return &page, nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", AllUsersEndPoint).Methods("GET")
	r.HandleFunc("/users", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUserEndPoint).Methods("PUT")
	r.HandleFunc("/users", DeleteUserEndPoint).Methods("DELETE")
	r.HandleFunc("/users/{id}", FindUserEndpoint).Methods("GET")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
