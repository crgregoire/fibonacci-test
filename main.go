package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type FibRead struct {
	FibonacciFromDB pq.Int64Array
}

var fibRead = FibRead{}

func fibonacci(n int) []int {

	tempSlice := make([]int, n+1)
	tempSlice[0] = 0

	a := 0
	b := 1
	// Iterate until desired position in sequence.
	for i := 1; i < n+1; i++ {
		// Use temporary variable to swap values.
		temp := a
		a = b
		b = temp + a
		tempSlice[i] = a
	}
	return tempSlice
}

func insertFibonacci(fibs []int) {

	db, err := sql.Open("postgres", "host="+os.Getenv("host")+" port="+os.Getenv("port")+" "+
		"user="+os.Getenv("username")+" dbname="+os.Getenv("dbname")+
		" password="+os.Getenv("password"))

	ins := "INSERT INTO fibonacci (fibs) VALUES ($1)"

	_, err = db.Exec(ins, pq.Array(fibs))

	if err != nil {
		panic(err)
	}
}

func createFibonacciTable() error {

	db, err := sql.Open("postgres", "host="+os.Getenv("host")+" port="+os.Getenv("port")+" "+
		"user="+os.Getenv("username")+" dbname="+os.Getenv("dbname")+
		" password="+os.Getenv("password"))

	query := `CREATE TABLE IF NOT EXISTS fibonacci (id serial PRIMARY KEY, fibs integer[])`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when creating fibonacci table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return err
	}
	log.Printf("Fibonacci table creation attempt (table may already be created). Rows affected: %d", rows)
	return nil
}

func deleteFibonacciTable(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	db, err := sql.Open("postgres", "host="+os.Getenv("host")+" port="+os.Getenv("port")+" "+
		"user="+os.Getenv("username")+" dbname="+os.Getenv("dbname")+
		" password="+os.Getenv("password"))

	if err != nil {
		panic(err)
	}

	sqlStatement := `DROP TABLE IF EXISTS fibonacci CASCADE;`

	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"Fibonacci table deleted. Don't blame me it was in the spec."}`)))
	if err != nil {
		return
	}
}

//The same function but it returns an int, and gets rid of the http and response
func getFibonacciNumberForTesting(position int) int {

	var Fibonaccis []int
	fibonacciPosition := 0

	fibonacciPosition = position

	if fibonacciPosition > 46 {
		fmt.Println("Out of Integer bounds")
	}
	Fibonaccis = make([]int, fibonacciPosition)
	Fibonaccis = fibonacci(fibonacciPosition)

	return Fibonaccis[fibonacciPosition]
}

func getFibonacciNumber(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	err := createFibonacciTable()

	if err != nil {
		panic(err)
	}

	var Fibonaccis []int
	fibonacciPosition := 0

	if val, ok := pathParams["desiredNum"]; ok {
		fibonacciPosition, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"message": "Please use a number in your URL params"}`))
			if err != nil {
				panic(err)
			}
			return
		}
		if fibonacciPosition > 46 {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"message": "Please use a number less than 47"}`))
			if err != nil {
				panic(err)
			}
			return
		}
		Fibonaccis = make([]int, fibonacciPosition)
		Fibonaccis = fibonacci(fibonacciPosition)
	}

	insertFibonacci(Fibonaccis)

	_, err = w.Write([]byte(fmt.Sprintf(`{"Fib(%d)==%d"}`, fibonacciPosition, Fibonaccis[fibonacciPosition])))
	if err != nil {
		panic(err)
	}
}

func getNumbersLessThanForTesting(lessThanNum int) int {

	fibonacci := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377}
	lessThan := 0

	for i := 0; i < len(fibonacci); i++ {
		if lessThanNum > fibonacci[i] {
			lessThan++
		}
	}

	return lessThan
}

func getNumbersLessThan(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	lessThanNumber := 0
	lessThan := 0
	var err error

	if val, ok := pathParams["desiredNum"]; ok {
		lessThanNumber, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(`{"message": "Please use a number in your URL params"}`))
			if err != nil {
				return
			}
			return
		}
	}

	fibRead.FibonacciFromDB = []int64(getFib(1))

	for i := 0; i < len(fibRead.FibonacciFromDB); i++ {
		if int64(lessThanNumber) > fibRead.FibonacciFromDB[i] {
			lessThan++
		}
	}

	_, err = w.Write([]byte(fmt.Sprintf(`{"There are %d numbers less than %d in the Slice"}`, lessThan, lessThanNumber)))
	if err != nil {
		return
	}
}

func getFib(id int) (fibs pq.Int64Array) {

	db, err := sql.Open("postgres", "host="+os.Getenv("host")+" port="+os.Getenv("port")+" "+
		"user="+os.Getenv("username")+" dbname="+os.Getenv("dbname")+
		" password="+os.Getenv("password"))

	if err != nil {
		panic(err)
	}

	sel := "SELECT fibs FROM fibonacci WHERE id=$1"

	// wrap the output parameter in pq.Array for receiving into it
	if err := db.QueryRow(sel, id).Scan(&fibs); err != nil {
		log.Fatal(err)
	}

	return
}

func main() {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fmt.Println("Running...: ")

	db, err := sql.Open("postgres", "host="+os.Getenv("host")+" port="+os.Getenv("port")+" "+
		"user="+os.Getenv("username")+" dbname="+os.Getenv("dbname")+
		" password="+os.Getenv("password"))

	if err != nil {
		panic(err)
	}
	log.Printf("Connected")
	defer db.Close()

	err = createFibonacciTable()
	if err != nil {
		return
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/fibonacci/{desiredNum}", getFibonacciNumber).Methods(http.MethodGet)
	api.HandleFunc("/fibonacci/less-than/{desiredNum}", getNumbersLessThan).Methods(http.MethodGet)
	api.HandleFunc("/fibonacci/delete-all", deleteFibonacciTable).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", r))
}
