# What cook

[![](https://images.unsplash.com/photo-1504754524776-8f4f37790ca0?ixlib=rb-1.2.1&ixid=MXwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHw%3D&auto=format&fit=crop&w=1500&q=80)](https://unsplash.com/photos/hrlvr2ZlUNk)

This is simple go app this is a simple app for finding foods 
that can be made from available ingredients

**INSTALL**
```shell script
go get github.com/go-kit/kit

go get gorm.io/gorm

go get gorm.io/driver/sqlite

go get gopkg.in/validator.v2
```
**PROJECT STRUCTURE**
```shell script
├── README.md
├── cmd
│   ├── main.go
│   └── wiring_test.go
├── domain
│   ├── crud_repository.go
│   ├── food.go
│   └── ingredient.go
├── food
│   ├── endpoint.go
│   ├── service.go
│   ├── service_test.go
│   └── trasnsport.go
├── gorm
│   ├── db.go
│   ├── repository.go
│   └── repository_test.go
├── helper
│   └── utils.go
├── ingredient
│   ├── endpoint.go
│   ├── service.go
│   ├── service_test.go
│   ├── transport.go
│   └── what_cook_test.db
└── what_cook.db
```

**EXAMPLE**

```http request
GET localhost:8080/food/byIngredients/
Accept: application/json

{"ingredients" :[
        "pasta","bacon"
    ]
}
```
request for "pasta","bacon" return carbonara(absent chicken) and omelet(has bacon but missing egg):

```json
    {
        "Foods": [
            {
                "Food": {
                    "ID": 1,
                    "CreatedAt": "0001-01-01T00:00:00Z",
                    "UpdatedAt": "0001-01-01T00:00:00Z",
                    "DeletedAt": null,
                    "Name": "carbonara",
                    "Description": "",
                    "IngredientWeights": [...],
                },
                "HasIngredients": [
                    {
                        "ID": 1,
                        "CreatedAt": "0001-01-01T00:00:00Z",
                        "UpdatedAt": "0001-01-01T00:00:00Z",
                        "DeletedAt": null,
                        "Name": "bacon",
                        "Calories": 100
                    },
                    {
                        "ID": 2,
                        "CreatedAt": "0001-01-01T00:00:00Z",
                        "UpdatedAt": "0001-01-01T00:00:00Z",
                        "DeletedAt": null,
                        "Name": "pasta",
                        "Calories": 50
                    }
                ],
                "AbsentIngredients": [
                    {
                        "ID": 3,
                        "CreatedAt": "0001-01-01T00:00:00Z",
                        "UpdatedAt": "0001-01-01T00:00:00Z",
                        "DeletedAt": null,
                        "Name": "chicken",
                        "Calories": 80
                    },
                ]
            },
            {
                "Food": {
                    "ID": 2,
                    "CreatedAt": "0001-01-01T00:00:00Z",
                    "UpdatedAt": "0001-01-01T00:00:00Z",
                    "DeletedAt": null,
                    "Name": "omelet",
                    "Description": "",
                    "IngredientWeights": [...],
                },
                "HasIngredients": [
                    {
                        "ID": 1,
                        "CreatedAt": "0001-01-01T00:00:00Z",
                        "UpdatedAt": "0001-01-01T00:00:00Z",
                        "DeletedAt": null,
                        "Name": "bacon",
                        "Calories": 100
                    }
                ],
                "AbsentIngredients": [
                    {
                        "ID": 4,
                        "CreatedAt": "0001-01-01T00:00:00Z",
                        "UpdatedAt": "0001-01-01T00:00:00Z",
                        "DeletedAt": null,
                        "Name": "egg",
                        "Calories": 70
                    }
                ]
            }
        ],
        "Err": null
    }
```
## bon appetit!

