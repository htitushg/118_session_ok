package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"118_session/assets"
	common "118_session/common"

	_ "github.com/go-sql-driver/mysql"
)

var Templates = template.New("Template")

func ReloadTemplates(templateFolder string) {
	newTemplates, err := LoadTemplates(templateFolder)
	if err != nil {
		log.Printf("Unable to load templates: %s", err)
		return
	}
	// override existing templates variable
	Templates = newTemplates
}
func LoadTemplates(folder string) (*template.Template, error) {
	template := template.New("Template")
	walkError := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			_, parseError := template.ParseFiles(path)
			if parseError != nil {
				return parseError
			}
		}
		return nil
	})
	return template, walkError
}

func main() {
	//database.InitDB()
	// On relie le fichier css et le favicon au nom static
	fmt.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	common.Routes()
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur\n", assets.Port)
	http.ListenAndServe(assets.Port, nil)
}
