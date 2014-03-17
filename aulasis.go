package main


import (
	"html/template"
	"net/http"
	"fmt"
)


const lenPath = len("/view/")


type Exercise struct{
	Text string
}



func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[lenPath:]
/*
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
*/
	
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, &Exercise{Text:"Este es un ejercicio de prueba"})
}



func main() {
	//http.HandleFunc("/view/", viewHandler)
	//http.ListenAndServe(":9090", nil)

	sc:=LoadServerConfig("srv/meta.xml")

	for c:=range sc.Courses{
		fmt.Printf("%s   %s\n",sc.Courses[c].Name, sc.Courses[c].Desc)
	}
}

