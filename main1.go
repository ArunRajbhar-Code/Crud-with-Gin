package main // entry point of the Go program

import (
	"math/rand" // used to generate random numbers (for IDs)
	"net/http"  // provides HTTP status codes and utilities
	"strconv"   // used to convert numbers to string
	"strings"   // used for string operations like split
	"time"      // used for token expiry time

	"github.com/gin-contrib/cors" // middleware to enable CORS
	"github.com/gin-gonic/gin" // Gin web framework
	"github.com/golang-jwt/jwt/v5" // JWT library for authentication
)

var secretKey = []byte("mysecretkey") // secret key used to sign JWT tokens

// ---------------- STRUCTS ----------------

type Movie struct {
	ID       string    `json:"id"` // movie ID
	Isbn     string    `json:"isbn"` // ISBN number of movie
	Title    string    `json:"title"` // movie title
	Director *Director `json:"director"` // pointer to director struct
}

type Director struct {
	Firstname string `json:"firstname"` // director first name
	Lastname  string `json:"lastname"` // director last name
}

type User struct {
	ID       string `json:"id"` // user ID
	Username string `json:"username"` // username for login
	Password string `json:"password"` // password for login
}

type Login struct {
	Username string `json:"username"` // username provided during login
	Password string `json:"password"` // password provided during login
}

var Movies []Movie // slice to store movies in memory
var Users []User // slice to store users in memory

// ---------------- JWT TOKEN ----------------

func generateToken(username string) (string, error) { // function to create JWT token

	claims := jwt.MapClaims{ // JWT payload data
		"username": username, // store username inside token
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // token expiry time (24 hours)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // create token with HS256 algorithm

	return token.SignedString(secretKey) // sign token using secret key
}

// ---------------- REGISTER API ----------------

func register(c *gin.Context) { // handler for user registration

	var user User // create empty user object

	if err := c.BindJSON(&user); err != nil { // bind incoming JSON request body to user struct
		c.JSON(400, gin.H{"error": err.Error()}) // return error if JSON is invalid
		return
	}

	// check if user already exists
	for _, u := range Users { // loop through all existing users
		if u.Username == user.Username { // check if username already exists
			c.JSON(400, gin.H{"message": "User already exists"}) // return error
			return
		}
	}

	user.ID = strconv.Itoa(rand.Intn(100000)) // generate random user ID and convert to string

	Users = append(Users, user) // add new user to Users slice

	c.JSON(201, gin.H{
		"message": "User registered successfully", // success response
	})
}

// ---------------- LOGIN API ----------------

func login(c *gin.Context) { // handler for login

	var loginUser Login // create login struct

	if err := c.BindJSON(&loginUser); err != nil { // bind JSON request to login struct
		c.JSON(400, gin.H{"error": err.Error()}) // return error if JSON invalid
		return
	}

	for _, user := range Users { // loop through registered users

		if user.Username == loginUser.Username && user.Password == loginUser.Password { // verify credentials

			token, err := generateToken(user.Username) // generate JWT token

			if err != nil { // check if token generation fails
				c.JSON(500, gin.H{"error": "Token generation failed"})
				return
			}

			c.JSON(200, gin.H{
				"token": token, // return token to client
			})

			return
		}
	}

	c.JSON(401, gin.H{"message": "Invalid credentials"}) // return unauthorized if login fails
}

// ---------------- JWT MIDDLEWARE ----------------

func AuthMiddleware() gin.HandlerFunc { // middleware function for authentication

	return func(c *gin.Context) { // middleware handler

		authHeader := c.GetHeader("Authorization") // read Authorization header from request

		if authHeader == "" { // check if header missing
			c.JSON(401, gin.H{"error": "Authorization header required"}) // return error
			c.Abort() // stop further request processing
			return
		}

		parts := strings.Split(authHeader, " ") // split header (Bearer TOKEN)

		if len(parts) != 2 { // ensure correct format
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1] // extract token part

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { // parse and validate token
			return secretKey, nil // provide secret key for validation
		})

		if err != nil || !token.Valid { // check if token invalid or expired
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next() // allow request to proceed if token is valid
	}
}

// ---------------- MOVIE APIs ----------------

// Get all movies
func getMovies(c *gin.Context) {
	c.JSON(200, Movies) // return all movies as JSON
}

// Get single movie
func getMovie(c *gin.Context) {

	id := c.Param("id") // get movie ID from URL parameter

	for _, item := range Movies { // loop through movies
		if item.ID == id { // check if movie ID matches
			c.JSON(200, item) // return movie
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"}) // return error if movie not found
}

// Create movie
func createMovie(c *gin.Context) {

	var movie Movie // create movie struct

	if err := c.BindJSON(&movie); err != nil { // bind JSON body to movie struct
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	movie.ID = strconv.Itoa(rand.Intn(100000000)) // generate random movie ID

	Movies = append(Movies, movie) // add movie to slice

	c.JSON(201, movie) // return created movie
}

// Update movie
func updateMovie(c *gin.Context) {

	id := c.Param("id") // get movie ID from URL

	var updatedMovie Movie // create updated movie struct

	if err := c.BindJSON(&updatedMovie); err != nil { // bind JSON request body
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for index, item := range Movies { // loop through movies

		if item.ID == id { // check matching movie

			updatedMovie.ID = id // keep same ID
			Movies[index] = updatedMovie // replace movie

			c.JSON(200, updatedMovie) // return updated movie
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"}) // return error if not found
}

// Delete movie
func deleteMovie(c *gin.Context) {

	id := c.Param("id") // get movie ID from URL

	for index, item := range Movies { // loop through movies

		if item.ID == id { // check matching movie

			Movies = append(Movies[:index], Movies[index+1:]...) // remove movie from slice

			c.JSON(200, gin.H{"message": "Movie deleted"}) // return success message
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"}) // return error if movie doesn't exist
}

// ---------------- MAIN ----------------

func main() {

	r := gin.Default() // create Gin router with default middleware (logger + recovery)

	// CORS
	r.Use(cors.Default()) // enable CORS so frontend apps can access API

	// Initial Movies
	Movies = append(Movies, Movie{ // add initial movie to slice
		ID:    "1",
		Isbn:  "2434432",
		Title: "Movie-1",
		Director: &Director{ // create director object
			Firstname: "John",
			Lastname:  "Carter",
		},
	})

	// Public Routes
	r.POST("/register", register) // register endpoint
	r.POST("/login", login) // login endpoint

	// Protected Routes
	protected := r.Group("/") // create route group

	protected.Use(AuthMiddleware()) // apply JWT middleware to group

	protected.GET("/movies", getMovies) // get all movies
	protected.GET("/movies/:id", getMovie) // get movie by ID
	protected.POST("/movies", createMovie) // create movie
	protected.PUT("/movies/:id", updateMovie) // update movie
	protected.DELETE("/movies/:id", deleteMovie) // delete movie

	r.Run(":8000") // start server on port 8000
}