package main

import (
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var Movies []Movie

// Get all movies
func getMovies(c *gin.Context) {
	c.JSON(200, Movies)
}

// Get single movie
func getMovie(c *gin.Context) {
	id := c.Param("id")

	for _, item := range Movies {
		if item.ID == id {
			c.JSON(200, item)
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"})
}

// Create movie
func createMovie(c *gin.Context) {
	var movie Movie

	if err := c.BindJSON(&movie); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	Movies = append(Movies, movie)

	c.JSON(201, movie)
}

// Update movie
func updateMovie(c *gin.Context) {
	id := c.Param("id")

	var updatedMovie Movie
	if err := c.BindJSON(&updatedMovie); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	for index, item := range Movies {
		if item.ID == id {
			updatedMovie.ID = id
			Movies[index] = updatedMovie
			c.JSON(200, updatedMovie)
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"})
}

// Delete movie
func deleteMovie(c *gin.Context) {
	id := c.Param("id")

	for index, item := range Movies {
		if item.ID == id {
			Movies = append(Movies[:index], Movies[index+1:]...)
			c.JSON(200, gin.H{"message": "Movie deleted"})
			return
		}
	}

	c.JSON(404, gin.H{"message": "Movie not found"})
}

func main() {
	r := gin.Default()

	// Initial Data
	Movies = append(Movies, Movie{
		ID:    "1",
		Isbn:  "2434432",
		Title: "Movie-1",
		Director: &Director{
			Firstname: "john",
			Lastname:  "carter",
		},
	})

	Movies = append(Movies, Movie{
		ID:    "2",
		Isbn:  "2434436",
		Title: "Movie-2",
		Director: &Director{
			Firstname: "monie",
			Lastname:  "markel",
		},
	})

	// Routes
	r.GET("/movies", getMovies)
	r.GET("/movies/:id", getMovie)
	r.POST("/movies", createMovie)
	r.PUT("/movies/:id", updateMovie)
	r.DELETE("/movies/:id", deleteMovie)

	r.Run(":8000") // Start server
}
