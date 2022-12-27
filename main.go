package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/datatypes"
)

var DataCore = []Core{}

func Palindrome(input string) bool {
	result := []byte{}
	for i := len(input); i > 0; i-- {
		result = append(result, input[i])
	}

	return input == string(result)
}

type Core struct {
	ID             uint     `json:"-" form:"-"`
	Language       string   `json:"language" form:"language"`
	Appeared       int      `json:"appeared" form:"appeared"`
	Created        []string `json:"created" form:"created"`
	Functional     bool     `json:"functional" form:"functional"`
	ObjectOriented bool     `json:"object-oriented" form:"object-oriented"`
	Relation       Relation `json:"relation" form:"relation"`
}

type Relation struct {
	ID           uint
	InfluencedBy []string `json:"influenced-by" form:"influenced-by"`
	Influences   []string `json:"influences" form:"influences"`
}

type Bahasa struct {
	ID             uint           `json:"-" form:"-"`
	Language       string         `json:"language" form:"language"`
	Appeared       int            `json:"appeared" form:"appeared"`
	Created        datatypes.JSON `json:"created" form:"created"`
	Functional     bool           `json:"functional" form:"functional"`
	ObjectOriented bool           `json:"object-oriented" form:"object-oriented"`
	Relation       datatypes.JSON `json:"relation" form:"relation"`
}

type T struct {
	Text string `json:"text" form:"text"`
}

func Language() echo.HandlerFunc {
	return func(c echo.Context) error {
		data := &Bahasa{
			Language:       "C",
			Appeared:       1972,
			Created:        datatypes.JSON([]byte(`["Dennis Ritchie"]`)),
			Functional:     true,
			ObjectOriented: false,
			Relation:       datatypes.JSON([]byte(`{"influenced-by": ["B", "ALGOL 68", "Assembly", "FORTRAN"], "influences": ["C++", "Objective-C", "C#", "Java", "Javascript", "PHP", "Go"]}`)),
		}

		return c.JSON(http.StatusOK, data)
	}
}

func GetPalindrome() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input T
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot bind input",
			})
		}

		Palindrome := []byte{}
		for i := len(input.Text) - 1; i >= 0; i-- {
			Palindrome = append(Palindrome, input.Text[i])
		}

		if input.Text == string(Palindrome) {
			return c.String(http.StatusOK, "Palindrome")
		} else {
			return c.String(http.StatusBadRequest, "Not Palindrome")
		}
	}
}

func PostLanguage() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input Core
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot bind input",
			})
		}

		DataCore = append(DataCore, input)

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "success insert new language",
		})
	}
}

func GetLanguage() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		cnv, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot identify id",
			})
		}

		return c.JSON(http.StatusOK, DataCore[cnv])
	}
}

func GetLanguages() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, DataCore)
	}
}

func PatchLanguage() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var input Core
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot bind input",
			})
		}

		cnv, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot identify id",
			})
		}
		DataCore[cnv] = input

		return c.JSON(http.StatusOK, DataCore[cnv])
	}
}

func DeleteLanguage() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		cnv, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "cannot identify id",
			})
		}

		DataCore = append(DataCore[:cnv], DataCore[cnv+1:]...)

		return c.JSON(http.StatusOK, DataCore)
	}
}

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Go developers")
	})

	e.POST("/palindrome", GetPalindrome())

	e.POST("/language", PostLanguage())

	e.GET("/language/:id", GetLanguage())

	e.GET("/languages", GetLanguages())

	e.PATCH("/language/:id", PatchLanguage())

	e.DELETE("/language/:id", DeleteLanguage())

	e.Start(":8000")
}
