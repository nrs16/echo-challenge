package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Error struct {
	Error   string
	Message string
	Status  int
}

func main() {
	e := echo.New()

	// pass middleware functions, I added some for logging
	e.Use(middleware.Logger())
	// one to handle headers (right now we are just printing them and passing one to context later we might other things
	//for headers in this function)
	e.Use(MiddlewareLogHeaders)
	// create logger for every request received witha unique id to be printed on logs and pass it to the context as well
	e.Use(MiddlewareLogger)

	// define routes
	e.POST("/routes", ListAiports)

	e.Logger.Fatal(e.Start(":8080"))
}

func ListAiports(c echo.Context) error {

	log := c.Get("log").(*slog.Logger)
	var routes [][]string

	err := c.Bind(&routes)
	if err != nil {
		log.Error(fmt.Sprintf("could not decode body: %s", err.Error()))
		resp := Error{Error: "invalid_body", Message: err.Error(), Status: http.StatusBadRequest}
		return c.JSON(http.StatusBadRequest, resp)
	}
	log.Info(fmt.Sprintf("printing body: %v", routes))

	/// Loop over array ,sanitize ,group into map so we can use it to order, find source
	var tuples = make(map[string]string)

	//// I will use 2 maps to store starts and ends of flights, to get the source of the whole trip without having to loop multiple times later to find it
	//// I will use the starts maps to store potential starts of the whole trip, and ends to store the end of EVERY trip
	//// while looping over flights if the start doesn't have an entry in end, it means it's potential for trip source so I add it to start map
	//// if the end of the flight has a entry in start, I will remove that entry from start, because we know for sure that it is not the source of the trip

	var starts = make(map[string]bool)
	var ends = make(map[string]bool)

	for _, v := range routes {
		if len(v) != 2 || v[0] == v[1] {
			continue
		}
		_, ok := tuples[v[0]]
		if ok {
			resp := Error{Error: "invalid_body", Message: "There are 2 destinations from same source", Status: http.StatusBadRequest}
			return c.JSON(http.StatusBadRequest, resp)
		}
		tuples[v[0]] = v[1]
		ends[v[1]] = true
		if !ends[v[0]] {
			starts[v[0]] = true
		}
		if starts[v[1]] {
			delete(starts, v[1])
		}
	}

	log.Info(fmt.Sprintf("printing tuples: %v\n printing start: %v", tuples, starts))

	if len(starts) != 1 {
		log.Error("could not find the start")
		resp := Error{Error: "invalid_body", Message: "There is either no start or too many options for start", Status: http.StatusBadRequest}
		return c.JSON(http.StatusBadRequest, resp)
	}
	///get the start

	var source string
	for key, _ := range starts {
		source = key
		break // Since there's only one key-value pair, exit the loop after the first iteration
	}

	var result []string
	var added = make(map[string]bool)
	for {
		// use this map to keep track of already passed airport, we cannot have any loops
		if added[source] {
			log.Error("There is a loop")
			resp := Error{Error: "invalid_body", Message: "There is a loop in the itinerary", Status: http.StatusBadRequest}
			return c.JSON(http.StatusBadRequest, resp)
		}
		result = append(result, source)
		added[source] = true
		v, ok := tuples[source]
		if ok {
			source = v
		} else {
			break
		}
	}
	return c.JSON(http.StatusOK, result)
}

func MiddlewareLogHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Read headers from the request
		headers := req.Header

		// Iterate over headers and print key-value pairs for developer visibility and easy debugginf using logs
		for key, values := range headers {
			for _, value := range values {
				/// I am passing the X-correlation-ID to context so I access it in logger. for a big project
				/// I would pass the headers as struct or map to the context but for this task I am sticking to the uuid
				/// we can also not print snensitive data like tokens
				if key == "X-Correlation-Id" {
					uuid := value
					c.Set("uuid", uuid)
				}
				slog.Info(fmt.Sprintf("%s: %s\n", key, value))
			}
		}
		return next(c)
	}
}

func MiddlewareLogger(next echo.HandlerFunc) echo.HandlerFunc {
	// create logger with specific request X-Correlation-ID for every request for good logs readability
	return func(c echo.Context) error {
		var uuid string
		val, ok := c.Get("uuid").(string)
		if ok {
			uuid = val
		}
		// we can handle error for this later if in some weird case X-Correlation-Id is not string

		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		logger = logger.With("uuid", uuid)
		/// wether logs should be printed on file or stdout should be configurable, for the sake of this assignment I chose stdout
		c.Set("log", logger)
		return next(c)
	}
}
