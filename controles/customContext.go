package controles

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CustomContext is our replacement for gin.Context
type CustomContext struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	Params   map[string]string
	formData url.Values
}

// NewCustomContext creates a new custom context
func NewCustomContext(w http.ResponseWriter, r *http.Request) *CustomContext {
	return &CustomContext{
		Request: r,
		Writer:  w,

		Params:   make(map[string]string),
		formData: r.Form,
	}
}

// HTML renders an HTML template
func (c *CustomContext) HTML(code int, name string, data interface{}) {
	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(code)

	tmplFiles := []string{
		"templates/layouts/layout.html",
		"templates/components/aside.html",
		"templates/components/header.html",
		"templates/components/vlink.html",
		"templates/" + name + ".html",
	}

	// Correct way to pass multiple file names
	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		http.Error(c.Writer, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(c.Writer, data)
	if err != nil {
		http.Error(c.Writer, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

// JSON responds with JSON
func (c *CustomContext) JSON(code int, obj interface{}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	if err := json.NewEncoder(c.Writer).Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Status sets the HTTP response status code
func (c *CustomContext) Status(code int) {
	c.Writer.WriteHeader(code)
}

// Redirect performs an HTTP redirect
func (c *CustomContext) Redirect(code int, location string) {
	http.Redirect(c.Writer, c.Request, location, code)
}

// Query returns the query parameter
func (c *CustomContext) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Param returns the URL parameter
func (c *CustomContext) Param(name string) string {
	return c.Params[name]
}

// GetHeader returns the header value
func (c *CustomContext) GetHeader(key string) string {
	return c.Request.Header.Get(key)
}

// Header sets a response header
func (c *CustomContext) Header(key, value string) {
	c.Writer.Header().Set(key, value)
}

// Bind binds the request data to a struct
func (c *CustomContext) Bind(obj interface{}) error {
	var err error

	contentType := c.Request.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "application/json"):
		decoder := json.NewDecoder(c.Request.Body)
		err = decoder.Decode(obj)

	case strings.Contains(contentType, "application/x-www-form-urlencoded"),
		strings.Contains(contentType, "multipart/form-data"):
		err = c.bindForm(obj)
	}

	return err
}

// bindForm binds form data to a struct
func (c *CustomContext) bindForm(_ interface{}) error {
	// In a production implementation, this would use reflection to populate the struct
	// For now, we're just returning an error to indicate it's not fully implemented
	return fmt.Errorf("form binding not fully implemented")
}

// PostForm returns the specified form value
func (c *CustomContext) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// PostFormArray returns the specified form values as an array
func (c *CustomContext) PostFormArray(key string) []string {
	return c.Request.Form[key]
}

// FormFile returns the specified file
func (c *CustomContext) FormFile(name string) (*multipart.FileHeader, error) {
	_, header, err := c.Request.FormFile(name)
	return header, err
}

// SaveUploadedFile saves the uploaded file to the specified path
func (c *CustomContext) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// SetCookie sets a cookie
func (c *CustomContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	http.SetCookie(c.Writer, cookie)
}

// SetSameSite sets SameSite cookie attribute
func (c *CustomContext) SetSameSite(sameSite http.SameSite) {
	// Not implemented fully as this would require modifying cookie settings
}
