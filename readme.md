This is a Go program with HTTP server which fetches currency rates from a public API and writes them to a MySQL database. It also provides an endpoint for users to query the database for currency rates.

The `main()` function is the entry point of the program. It initializes the connection to the database and the HTTP client object that fetches data from the API. The initalization function also parses a JSON configuration file for various program settings.

The `downloadCurr()` function handles HTTP requests that fetch new currency rates from the API, in this case the National Bank of Kazakhstan. The function parses the HTTP request for a date parameter and converts it to a time.Time object. It then fetches currency rates from the API for that date and writes the rates to the database.

The `uploadCurr()` function handles HTTP requests that query the database for currency rates. It parses the HTTP request for a date and currency code parameter and converts the date to a time.Time object. It then looks up the currency rate for that date and currency code in the database and returns it to the user.

The `commonMiddleware()` function is a custom middleware that sets common HTTP headers for all HTTP responses. It sets the content type to JSON and allows CORS requests from any origin.

To run this program, make sure that MySQL is installed and running on the local machine, and that a table with the correct schema is created in the database. Then the program can be run from the command-line interface, navigating to the program directory and running the command `go run .`. The program will start a HTTP server on the port specified in the configuration file.