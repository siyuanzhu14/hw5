package app

import (
	"encoding/json"
	"html/template"
	"net/http"
    "fmt"
    
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", handlePata)
	http.HandleFunc("/norikae", handleNorikae)
}

// このディレクトリーに入っているすべての「.html」終わるファイルをtemplateとして読み込む。
var tmpl = template.Must(template.ParseGlob("*.html"))

// Templateに渡す内容を分かりやすくするためのtypeを定義しておきます。
// （「Page」という名前などは重要ではありません）。
type Page struct {
	A    string
	B    string
	Pata string
}


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
        }
        for i:=len(Arune); i<len(Brune); i++{
            x = x + string(Brune[i])
        }
    }else { for i:=0; i<len(Brune); i++  {
            x = x + string(Arune[i]) + string(Brune[i])
        }  
        for i:=len(Brune); i<len(Arune); i++{
            x = x + string(Arune[i])
        }
    }     
    content.Pata = x
    fmt.Printf(x)
	// example.htmlというtemplateをcontentの内容を使って、{{.A}}などのとこ
	// ろを実行して、内容を埋めて、wに書き込む。
	tmpl.ExecuteTemplate(w, "example.html", content)
}

// LineはJSONに入ってくる線路の情報をtypeとして定義している。このJSON
// にこの名前にこういうtypeのデータが入ってくるということを表している。
type Line struct {
	Name     string
	Stations []string
}

type Route struct {
    Start  string
    End  string
    Network TransitNetwork
    Connected map[string][]Eki
    
    //Result TransitNetwork
    Result []Eki
    //Fastest []TransitNetwork
}

// TransitNetworkは http://fantasy-transit.appspot.com/net?format=json
// の一番外側のリストのことを表しています。
type TransitNetwork []Line



type Eki struct {
    LineName    string
    EkiName   string
}

func BuildEkiUnits(network TransitNetwork) [][]Eki{
    
    var AllStations [][]Eki
    
    for _, line:= range network{
        
        ekis := line.Stations
       
        var temp []Eki
        for j:=0; j< len(ekis); j++ {
            var eki Eki
            eki.EkiName = ekis[j]
            eki.LineName = line.Name
            temp = append(temp, eki)
            
        }
        AllStations = append(AllStations, temp)
    }
    
    return AllStations
}

func BuildConnected(ekiUnits [][]Eki)(map[string][]Eki){
    
    connected := make(map[string][]Eki)
    for _, line:= range ekiUnits{
       
        for i:=0; i< len(line); i++ {
            
            eki := line[i].EkiName
            
               
                if i > 0 && i <= len(line) - 2{
                    connected[eki] = append(connected[eki], line[i-1], line[i+1])
                }else if i == 0{
                    
                    connected[eki] = append(connected[eki], line[i+1])
                    
                }else{
                    connected[eki] = append(connected[eki], line[i-1])
                }
                
                
            }
        }
    return connected
}

func makeEki(name string)(Eki){
    res := Eki{
        EkiName: name,
        LineName: "whatever",
    }
    return res
}

func EkiInRoute(a string, list []Eki) bool {
    for _, b := range list {
        if b.EkiName == a {
            return true
        }
    }
    return false
}


//a:start station. b:end station
func BFS(connected map[string][]Eki, start string, end string) []Eki{
    
    q := make([][]Eki,0)
    q = append(q, []Eki{makeEki(start)})
    route := make([]Eki,0)
    
    for len(q) != 0 {
        
        route = q[0]
        here := route[len(route)-1].EkiName  

        if here == end {
            return route
        }
          
        //a := len(q)
        
        
        newq := make([][]Eki,0)
        newroute := make([]Eki,len(route)+1)
        for _, temp := range connected[here]{
          go func(newEki Eki){
                if EkiInRoute(newEki.EkiName, route) == false{
                    newroute = append(route, newEki)
                    newq = append(newq, newroute)
                }      
            }(temp)
        }
        q = append(q, newq...)
        
        q = q[1:]   
    
    
    }
    return nil
    
}

func handleNorikae(w http.ResponseWriter, r *http.Request) {
	// Appengineの「Context」を通してAppengineのAPIを利用する。
	ctx := appengine.NewContext(r)

	// clientはAppengine用のHTTPクライエントで、他のウェブページを読み込
	// むことができる。
	client := urlfetch.Client(ctx)

	// JSONとしての路線グラフ内容を読み込む
	resp, err := client.Get("http://fantasy-transit.appspot.com/net?format=json")
	if err != nil {
		panic(err)
	}

	// 読み込んだJSONをパースするJSONのDecoderを作る。
	decoder := json.NewDecoder(resp.Body)

	// JSONをパースして、「network」に保存する。
	var network TransitNetwork
	if err := decoder.Decode(&network); err != nil {
		panic(err)
	}

	// handleExampleと同じようにtemplateにテンプレートを埋めて、出力する。
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    

    content := Route{
		Start: r.FormValue("a"),
		End: r.FormValue("b"),
	}
    

    content.Network = network

    ekiUnits := BuildEkiUnits(content.Network)
    connected := BuildConnected(ekiUnits)
    

    
    content.Connected = connected
    result := BFS(content.Connected, content.Start, content.End) 
    content.Result  = result 
	
    tmpl.ExecuteTemplate(w, "norikae.html", content)
}
