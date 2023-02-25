# Go Sample App

Go sample app is an application using standary library to build RESTful APIs

* **Key features**:
    - Using standard libraries only
    - In Memory implementation for storage

* **Run application**:
    * Run via Make: ```make run-app```
    -  Run via Docker: ```docker run -port 8000:8000 --rm -it $(docker build -q .)```
    - Run via Term: ```APP_PORT=8000 go run main.go```
