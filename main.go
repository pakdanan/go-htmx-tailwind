// main.go
package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	items = []Item{{ID: 1, Name: "Sample Item"},
		{ID: 2, Name: "Another Item A"},
		{ID: 3, Name: "Another Item B"},
		{ID: 4, Name: "Another Item C"},
	}
	tmpl   = template.Must(template.ParseGlob("templates/*.html"))
	nextID = 3
)

func main() {
	e := echo.New()
	e.Static("/static", "./static")
	e.GET("/", func(c echo.Context) error {
		return tmpl.ExecuteTemplate(c.Response(), "index.html", nil) //c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/items", listItems)
	e.GET("/add", addItemPage)
	e.POST("/add", addItem)
	e.GET("/edit/:id", editItemPage)
	e.POST("/edit/:id", editItem)
	e.POST("/delete/:id", deleteItem)
	e.Logger.Fatal(e.Start(":8080"))
}

func listItems(c echo.Context) error {
	search := c.QueryParam("search")
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = len(items)
	}

	var filtered []Item
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(search)) {
			filtered = append(filtered, item)
		}
		if len(filtered) >= limit {
			break
		}
	}

	data := struct {
		Items  []Item
		Search string
		Limit  int
	}{filtered, search, limit}

	return tmpl.ExecuteTemplate(c.Response(), "list.html", data)
}

func addItemPage(c echo.Context) error {
	return tmpl.ExecuteTemplate(c.Response(), "add.html", nil)
}

func addItem(c echo.Context) error {
	name := c.FormValue("name")
	items = append(items, Item{ID: nextID, Name: name})
	nextID++
	return c.Redirect(http.StatusSeeOther, "/")
}

func editItemPage(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for _, item := range items {
		if item.ID == id {
			return tmpl.ExecuteTemplate(c.Response(), "edit.html", item)
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func editItem(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.FormValue("name")
	for i, item := range items {
		if item.ID == id {
			items[i].Name = name
			break
		}
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func deleteItem(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			break
		}
	}
	return c.Redirect(http.StatusSeeOther, "/")
}
