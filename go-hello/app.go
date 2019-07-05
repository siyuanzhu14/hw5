package app

import (
	"net/http"
    "html/template"
)

type Page struct {
	A    string
	B    string
	Pata string
}

var tmpl = template.Must(template.ParseGlob("*.html"))

func handlePata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// templateに埋める内容をrequestのFormValueから用意する。
	content := Page{
		A: r.FormValue("a"),
		B: r.FormValue("b"),
	}

	// とりあえずPataを簡単な操作で設定しますけど、すこし工夫をすれば
	// パタトクカシーーができます。
	//content.Pata = content.A + content.B
	x := ""
    Arune := []rune(content.A)
    Brune := []rune(content.B)
    if len(Arune) <= len(Brune){
        for i:=0; i<len(Arune); i++ {
            x = x + string(Arune[i]) + string(Brune[i])
            //x = x + string(content.A[i]) + string(content.B[i])
        }
        for i:=len(Arune); i<len(Brune); i++{
            x = x + string(Brune[i])
        }
    }else { for i:=0; i<len(Brune); i++  {
            x = x + string(Arune[i]) + string(Brune[i])
            //x = x + string(content.A[i]) + string(content.B[i])
        }  
        for i:=len(Brune); i<len(Arune); i++{
            x = x + string(Arune[i])
        }
    }     
    content.Pata = x
	// example.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	tmpl.ExecuteTemplate(w, "example.html", content)
}


func init() {
	http.HandleFunc("/", handlePata)
}


func handlePata_original(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "Hello world!\n")
}
