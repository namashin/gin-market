package main

import (
	"bytes"
	"encoding/json"
	"gin-market/dto"
	"gin-market/infra"
	"gin-market/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// Test user
var TestUser = dto.SignUpInput{
	Email:    "admin@example.com",
	Password: "admin_password",
}

func TestMain(m *testing.M) {
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func setUp() *gin.Engine {
	db := infra.SetUpDB()
	err := db.AutoMigrate(&models.Item{}, &models.User{})
	if err != nil {
		log.Fatal(err)
	}

	router := setUpRouter(db)
	return router
}

// _signUp method for Test User
func _signUp(t *testing.T, router *gin.Engine, signUpTestUser dto.SignUpInput) dto.SignUpInput {
	// sign up request body
	signUpReqBody, _ := json.Marshal(signUpTestUser)

	// http request
	signUpReq, _ := http.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(signUpReqBody))

	// http response
	w := httptest.NewRecorder()

	// do API http request
	router.ServeHTTP(w, signUpReq)

	// ? signup succeeded
	assert.Equal(t, http.StatusCreated, w.Code)

	// get response new user data
	var signUpRespBody map[string]dto.SignUpInput
	_ = json.Unmarshal(w.Body.Bytes(), &signUpRespBody)

	return signUpRespBody["user"]
}

// _login method for Test User to login
func _login(t *testing.T, router *gin.Engine, loginUser dto.LoginInput) string {
	// request body
	loginReqBody, _ := json.Marshal(loginUser)

	// http request
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginReqBody))

	// http response
	w := httptest.NewRecorder()

	// do API http Request
	router.ServeHTTP(w, req)

	// ? login succeeded
	assert.Equal(t, http.StatusOK, w.Code)

	// get response token data
	var loginRespBody map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &loginRespBody)

	return loginRespBody["token"]
}

// _createTestItems make count numbers item
func _createTestItems(t *testing.T, router *gin.Engine, token string, count int) []uint {
	var createdItemIDs []uint

	for i := 1; i <= count; i++ {
		newItem := dto.CreateItemInput{
			Name:        "test" + strconv.Itoa(i),
			Price:       uint(i * 100),
			Description: "No." + strconv.Itoa(i),
		}

		// http request body
		reqBody, _ := json.Marshal(newItem)

		// http request
		req, _ := http.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(reqBody))

		// set token in header
		req.Header.Set("Authorization", "Bearer "+token)

		// http response
		w := httptest.NewRecorder()

		// do API request
		router.ServeHTTP(w, req)

		// ? new item is created successfully
		assert.Equal(t, http.StatusCreated, w.Code)

		// save new created item ID
		var resData map[string]models.Item
		_ = json.Unmarshal(w.Body.Bytes(), &resData)

		createdItemIDs = append(createdItemIDs, resData["data"].ID)
	}

	return createdItemIDs
}

func _findMyAll(t *testing.T, router *gin.Engine, token string) []models.Item {
	// http request
	req, err := http.NewRequest(http.MethodGet, "/items/mine", nil)
	assert.Equal(t, err, nil)

	// http response
	w := httptest.NewRecorder()

	// set Authorization token in header
	req.Header.Set("Authorization", "Bearer "+token)

	// do API FindMyAll request
	router.ServeHTTP(w, req)

	// put http response result in res
	var myItems map[string][]models.Item
	err = json.Unmarshal([]byte(w.Body.String()), &myItems)
	assert.Equal(t, err, nil)

	return myItems["data"]
}

func _update(t *testing.T, router *gin.Engine, token string, updateId string, updateItem dto.UpdateItemInput) models.Item {
	// update request body
	updateReqBody, err := json.Marshal(updateItem)
	assert.Equal(t, err, nil)

	// update http request
	updateReq, err := http.NewRequest(http.MethodPut, "/items/"+updateId, bytes.NewBuffer(updateReqBody))
	assert.Equal(t, err, nil)

	// set token in header
	updateReq.Header.Set("Authorization", "Bearer "+token)

	// http response
	w := httptest.NewRecorder()

	// do API request
	router.ServeHTTP(w, updateReq)

	// ? update succeeded
	assert.Equal(t, http.StatusOK, w.Code)

	// updated item
	var updateRespBody map[string]models.Item
	err = json.Unmarshal(w.Body.Bytes(), &updateRespBody)
	assert.Equal(t, err, nil)

	return updateRespBody["data"]
}

func _delete(t *testing.T, router *gin.Engine, token string, deleteId string) {
	// http request
	req, _ := http.NewRequest(http.MethodDelete, "/items/"+deleteId, nil)
	// http response
	w := httptest.NewRecorder()
	// set Authorization token in header
	req.Header.Set("Authorization", "Bearer "+token)
	// do API FindMyAll request
	router.ServeHTTP(w, req)

	// assert status check
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestFindMyAll(t *testing.T) {
	// set up router
	router := setUp()

	/* 1. Sign up by new user */
	signUpUser := _signUp(t, router, TestUser)

	/* 2. Login */
	token := _login(t, router, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. Create 10 items by login user */
	newItemIDs := _createTestItems(t, router, token, 10)

	/* 4. Test FindMyAll method */
	myItems := _findMyAll(t, router, token)

	/* 5. assert check */
	assert.Equal(t, len(newItemIDs), len(myItems))
}

func TestCreate(t *testing.T) {
	// set up test
	router := setUp()

	/* 1. Sign up by new user */
	signUpUser := _signUp(t, router, TestUser)

	/* 2. Login */
	token := _login(t, router, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. asset check (not yet created) */
	// get all items
	currentItems := _findMyAll(t, router, token)
	// ? is empty
	assert.Equal(t, 0, len(currentItems))

	/* 4. create 10 items */
	newItemIDs := _createTestItems(t, router, token, 10)

	/* 5. assert check (created 10 items) */
	assert.Equal(t, len(newItemIDs), 10)
}

func TestDelete(t *testing.T) {
	// set up test
	router := setUp()

	/* 1. Sign up by new user */
	signUpUser := _signUp(t, router, TestUser)

	/* 2. Login */
	token := _login(t, router, dto.LoginInput{Email: signUpUser.Email, Password: signUpUser.Password})

	/* 3. Create 10 items by login user */
	newItemIDs := _createTestItems(t, router, token, 10)

	/* 4. Delete one item */
	// delete first index in new items
	deleteId := strconv.Itoa(int(newItemIDs[0]))
	_delete(t, router, token, deleteId)

	/* 5. get items */
	myItems := _findMyAll(t, router, token)

	/* 6. assert check */
	assert.Equal(t, 9, len(myItems))
}

func TestUpdate(t *testing.T) {
	/* 1. set up router */
	router := setUp()

	/* 2. Signup by test user */
	newUser := _signUp(t, router, TestUser)

	/* 3. Login */
	token := _login(t, router, dto.LoginInput{Email: newUser.Email, Password: newUser.Password})

	/* 4. Create 10 items */
	itemIDs := _createTestItems(t, router, token, 10)

	/* 5. Update one of the created items */
	updatedName := "updated name"
	updatedPrice := uint(9999)
	updatedDescription := "updated"

	updateItem := dto.UpdateItemInput{
		Name:        &updatedName,
		Price:       &updatedPrice,
		Description: &updatedDescription,
	}
	// update first index item
	updateId := strconv.Itoa(int(itemIDs[0]))
	updatedItem := _update(t, router, token, updateId, updateItem)

	/* 6. assert check */
	assert.Equal(t, updatedItem.Name, updatedName)
	assert.Equal(t, updatedItem.Price, updatedPrice)
	assert.Equal(t, updatedItem.Description, updatedDescription)
}
