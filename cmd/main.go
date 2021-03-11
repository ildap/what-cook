package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"what_cook/domain"
	"what_cook/food"
	gormdep "what_cook/gorm"
	"what_cook/ingredient"
)

func main() {
	listen := flag.String("listen", ":8080", "HTTP listen address")

	logger := log.NewLogfmtLogger(os.Stderr)

	httpLogger := log.With(logger, "component", "http")

	var (
		db                   *gorm.DB
		ingredientRepository domain.IngredientRepository
		ingredientService    domain.IngredientService
		foodRepository       domain.FoodRepository
		foodService          domain.FoodService
	)

	db = gormdep.SqliteDbSession(gormdep.DSN_SQLITE)
	ingredientRepository = gormdep.NewIngredientRepository(db)
	ingredientService = ingredient.NewService(ingredientRepository)

	foodRepository = gormdep.NewFoodRepository(db)
	foodService = food.NewFoodService(foodRepository)

	mux := http.NewServeMux()
	mux.Handle("/ingredient/", ingredient.MakeHandler(ingredientService, httpLogger))
	mux.Handle("/food/", food.MakeHandler(foodService, httpLogger))
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)

	go func() {
		logger.Log("transport", "http", "address", *listen, "msg", "listening")
		errs <- http.ListenAndServe(*listen, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
