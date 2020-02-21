# Go Default

A bare minimum repository setup, with SQLite & Redis

## Getting Started

Make sure you have installed `golang dep`. If not please [install dep](https://github.com/golang/dep) first.

Run dep to get vendor
```
dep ensure -v
```

Then run the app
```
go run app.go
```

Check if everything is OK with
```
curl --location --request GET 'http://127.0.0.1:2000/health/check'
```

This is the expected response
```
{
    "header": {
        "process_time": "482.674µs"
    },
    "data": {
        "DB": {
            "DB Default": "OK"
        },
        "Redis": {
            "Redis Default": "PONG"
        }
    },
    "is_success": true
}
```
Now, make the app you are inspired to do!

## References

[Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047), the repository structure is follow the article's example. Adjusted with small modification and extra utils package.