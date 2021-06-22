### Fibonacci-test

Hi there and thank you for checking out my Go program.

This program exposes three endpoints per the spec. 
The endpoints: 
- Return a fibonacci sequence number supplied by the user.
- Return the number of numbers less than the number supplied by the user.
- Deletes the fibonacci sequence.

The fibonacci sequence is stored in a Postgres DB as an array.

To start, clone the repo.

Create a Postgres DB wherever you would like, such as locally or on RDS.

Next create a .env file with the following parameters:

```dockerfile
host=localhost or rds
username=postgres
password=yourpassword
port=5432
dbname=postgres
```

Next, run the `go get` command in your repo folder.

If there are no errors, run `go run main.go` and use Postman or something like it to use the API.

Here are the routes:
```go
    localhost:8080/api/v1/fibonacci/{desiredNum} 
            // ^ creates and stores a fibonacci sequence in the DB
    localhost:8080/api/v1/fibonacci/less-than/{desiredNum} 
            // ^ tells you how many numbers are less than the one you provide
    localhost:8080/api/v1/fibonacci/delete-all
            // ^ deletes the fibonacci sequences
```