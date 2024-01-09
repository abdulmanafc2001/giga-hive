package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gin-gonic/gin"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

// Home displays the home page with some sample data
func Home(c *gin.Context) {
	err := renderPage(c.Writer, "home.jet", nil)
	if err != nil {
		_, _ = fmt.Fprint(c.Writer, "Error executing template:", err)
	}
}

// SendAlertToConnectedUsers displays the page for sending data via websockets
func SendAlertToConnectedUsers(c *gin.Context) {
	err := renderPage(c.Writer, "send.jet", nil)
	if err != nil {
		_, _ = fmt.Fprint(c.Writer, "Error executing template:", err)
	}
}

// renderPage renders the page using Jet templates
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println("Unexpected template err:", err.Error())
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	return nil
}
