package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

var id = 0
type Contact struct {
		Name string
		Email string
		Id int
}
type Contacts = []Contact

func createNewContact (name, email string) Contact {
	id++
	return Contact{
		Name: name,
		Email: email,
		Id: id,	
	}
}


func (d *Data) hasEmail (email string) bool {
	for _, contact := range d.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

type Data struct {
	Contacts Contacts
}


func (d *Data) indexOf(id int) int {
	for i, contact := range d.Contacts {
		if id == contact.Id {
			return i
		}
	}

	return -1
}

func newData () Data {
	return Data{
		Contacts: []Contact{
			createNewContact("John Doe", "lol@gmila.com"),
			createNewContact("clara", "clara@gmila.com"),
		},
	}
} 

type FromData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FromData {
	return FromData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FromData
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()
	e.Use(middleware.Logger())

	page := newPage()

	e.Static("/images", "images")
	e.Static("/css", "css")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	});

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		
		if (page.Data.hasEmail(email)) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(http.StatusUnprocessableEntity, "form", formData)
		}

		contact := createNewContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	});


	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(5 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid id")
		}

		index := page.Data.indexOf(id)
		if index == -1 {
			return c.String(http.StatusNotFound, "Contact not found")
		}

		page.Data.Contacts = append(page.Data.Contacts[:index], page.Data.Contacts[index+1:]...)
		return c.NoContent(http.StatusOK)
	});



	e.Logger.Fatal(e.Start(":42069"))
}


