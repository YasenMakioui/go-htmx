package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Todo struct {
	Id          int
	Title       string
	Description string
	Completed   bool
	Expires     string
}

type TodoList struct {
	Todos []*Todo
}

func (t *TodoList) Delete(id int) {
	for i, todo := range t.Todos {
		if todo.Id == id {
			copy(t.Todos[i:], t.Todos[i+1:])
			t.Todos[len(t.Todos)-1] = &Todo{}
			t.Todos = t.Todos[:len(t.Todos)-1]
		}
	}
}

func NewTodoList() *TodoList {
	return &TodoList{}
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	id := 0

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()

	e.Renderer = t

	todolist := NewTodoList()

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", todolist)
	})

	e.POST("/todo", func(c echo.Context) error {

		var isCompleted bool

		title := c.FormValue("title")
		description := c.FormValue("description")
		expires := c.FormValue("expires")
		completed := c.FormValue("completed")

		if completed == "on" {
			isCompleted = true
		} else {
			isCompleted = false
		}

		todo := &Todo{
			Id:          id,
			Title:       title,
			Description: description,
			Expires:     expires,
			Completed:   isCompleted,
		}

		id++

		todolist.Todos = append(todolist.Todos, todo)

		return c.Render(http.StatusOK, "todo", todo)
	})

	e.DELETE("/todo/:id", func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		todolist.Delete(id)

		return c.NoContent(200)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
