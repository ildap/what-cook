package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"what_cook/domain"
	"what_cook/food"
	"what_cook/gorm"
	"what_cook/ingredient"
)

type requestResponseTest struct {
	method        string
	url           string
	body          string
	testResponses []testResponse
}

type testResponse func(t *testing.T, httpCode int, responseBody io.Reader)

func responseBodyContains(value string) testResponse {
	return func(t *testing.T, httpCode int, responseBody io.Reader) {
		b, _ := ioutil.ReadAll(responseBody)
		body := string(b)
		if !strings.Contains(body, value) {
			t.Errorf("respone is not contains \"%s\"", value)
		}
	}
}

func responseStatusIs(code int) testResponse {
	return func(t *testing.T, httpCode int, responseBody io.Reader) {
		if httpCode != code {
			t.Errorf("http status is not %d (%d)", code, httpCode)
		}
	}
}

func checkIngredient(t *testing.T, httpCode int, responseBody io.Reader) {
	var body struct {
		Ingredient domain.Ingredient `json:"Ingredient"`
	}
	err := json.NewDecoder(responseBody).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if !body.Ingredient.Equal(&testIngredient) {
		t.Error("ingredient not equal test ingredient")
	}
}

func checkFood(t *testing.T, httpCode int, responseBody io.Reader) {
	var body struct {
		Food domain.Food `json:"Food"`
	}
	err := json.NewDecoder(responseBody).Decode(&body)
	if err != nil {
		t.Error(err)
	}
	if !body.Food.Equal(&testFood) {
		t.Error("food not equal test food")
	}
}

var (
	logger            log.Logger
	testIngredient    domain.Ingredient
	testFood          domain.Food
	testFoods         []domain.Food
	foodService       domain.FoodService
	ingredientService domain.IngredientService
	baseUrl           string
)

func TestMain(m *testing.M) {
	// setup
	logger = log.NewLogfmtLogger(os.Stderr)
	db := gorm.SqliteDbSession(gorm.DSN_SQLITE_TEST)
	gorm.ClearData(db)
	testIngredient = gorm.CreateRandomIngredient(db)
	testFood = gorm.CreateRandomFood(db)
	testFoods = []domain.Food{
		gorm.CreateRandomFood(db),
		gorm.CreateRandomFood(db),
		gorm.CreateRandomFood(db),
	}
	ingredientRepository := gorm.NewIngredientRepository(db)
	ingredientService = ingredient.NewService(ingredientRepository)
	foodRepository := gorm.NewFoodRepository(db)
	foodService = food.NewFoodService(foodRepository)
	// server
	mux := http.NewServeMux()
	mux.Handle("/ingredient/", ingredient.MakeHandler(ingredientService, logger))
	mux.Handle("/food/", food.MakeHandler(foodService, logger))
	http.Handle("/", accessControl(mux))
	srv := httptest.NewServer(mux)
	defer srv.Close()
	baseUrl = srv.URL
	// run tests
	m.Run()
}

func TestWiring(t *testing.T) {
	var requestResponseTestData = []requestResponseTest{
		// check read
		{
			method: "GET",
			url:    "/ingredient/1",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				checkIngredient,
			},
		},
		{
			method:        "GET",
			url:           "/ingredient/0",
			testResponses: []testResponse{responseStatusIs(http.StatusNotFound)},
		},
		// check create
		{
			method: "POST",
			url:    "/ingredient/",
			body:   "{\"ingredient\" : {\"name\" : \"chocolate\",\"calories\": 546}}",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				responseBodyContains("\"IngredientID\""),
			},
		},
		// check update
		{
			method: "PUT",
			url:    "/ingredient/2",
			body:   "{\"ingredient\" : {\"name\" : \"chocolate\",\"calories\" : 10}}",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
			},
		},
		//check delete
		{
			method: "DELETE",
			url:    "/ingredient/2",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
			},
		},
		// FOOD
		// check read
		{
			method:        "GET",
			url:           "/food/0",
			testResponses: []testResponse{responseStatusIs(http.StatusNotFound)},
		},
		{
			method: "GET",
			url:    "/food/1",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				checkFood,
			},
		},
		// check create
		{
			method: "POST",
			url:    "/food/",
			body:   "{\"food\":{\"name\":\"pasta\"}}",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				responseBodyContains("\"FoodId\""),
			},
		},
		// check update
		{
			method: "PUT",
			url:    "/food/1",
			body:   "{\"food\":{\"name\":\"carbonara\"}}",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
			},
		},
		// check delete
		{
			method: "DELETE",
			url:    "/food/1",
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				func(t *testing.T, httpCode int, responseBody io.Reader) {
					_, err := foodService.Get(1)
					if err != domain.ModelNotFoundError {
						t.Error("food is not deleted")
					}
				},
			},
		},
		// check query by ingredients
		{
			method: "GET",
			url:    "/food/byIngredients/",
			body: func() string {
				ingredients := make([]string, len(testFoods[0].IngredientWeights))
				for i, ingredientWeight := range testFoods[0].IngredientWeights {
					ingredients[i] = fmt.Sprintf("\"%s\"", ingredientWeight.Ingredient.Name)
				}
				return "{\"ingredients\":[" + strings.Join(ingredients, ",") + "]}"
			}(),
			testResponses: []testResponse{
				responseStatusIs(http.StatusOK),
				responseBodyContains(testFoods[0].Name),
			},
		},
	}

	for _, testcase := range requestResponseTestData {
		req, _ := http.NewRequest(testcase.method, baseUrl+testcase.url, strings.NewReader(testcase.body))
		resp, _ := http.DefaultClient.Do(req)

		logger.Log("url", req.URL, "method", req.Method)
		for _, testResponse := range testcase.testResponses {
			testResponse(t, resp.StatusCode, resp.Body)
		}
	}
}
