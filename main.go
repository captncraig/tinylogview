package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join("/var/log", r.URL.Path)
		info, err := os.Stat(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if info.IsDir() {
			infos, err := ioutil.ReadDir(path)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Add("Content-Type", "text/html")
			fmt.Println(listTpl.Execute(w, map[string]interface{}{"Path": path, "Infos": infos}))
		} else {
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte("<pre>"))
			s := string(dat)
			s = strings.Replace(s, "<", "&lt;", -1)
			s = strings.Replace(s, ">", "&gt;", -1)
			w.Write([]byte(s))
			w.Write([]byte("</pre><script>window.scrollTo(0,document.body.scrollHeight);</script>"))
		}
	})
	http.ListenAndServe(":4568", nil)
}

var listTpl = template.Must(template.New("").Parse(`
	<h3>Index of {{.Path}}</h3>
	<table>
	{{range .Infos}}
	<tr><td>{{if .IsDir}}D{{end}}</td><td><a href="./{{.Name}}{{if .IsDir}}/{{end}}">{{.Name}}{{if .IsDir}}/{{end}}</a></td></tr>
	{{end}}
	</table>
`))
