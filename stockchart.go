//  package main
//
//  import (
//    "fmt"
//    "github.com/pdevty/stockchart"
//    "io/ioutil"
//    "os"
//  )
//
//  func main() {
//    // params json format
//    // brand id,name from yahoo finance japan
//    params := `[
//      {"id":"7618","name":"ピーシーデポ"},
//      {"id":"6157","name":"日進工具"},
//      {"id":"7821","name":"前田工繊"},
//      {"id":"7917","name":"藤森工業"},
//      {"id":"4681","name":"リゾートトラスト"},
//      {"id":"4301","name":"アミューズ"},
//      {"id":"4290","name":"プレステージ"},
//      {"id":"2780","name":"コメ兵"},
//      {"id":"2695","name":"くらコーポレーション"},
//      {"id":"2695","name":"内外トランスライン"}
//    ]`
//    // new
//    sc, err := stockchart.New(params)
//    if err != nil {
//      panic(err)
//    }
//    // return chart html
//    fmt.Println(sc.Html())
//    // return chart csv
//    fmt.Println(sc.Csv())
//    // create chart html file
//    ioutil.WriteFile("index.html",
//      []byte(sc.Html()), os.ModePerm)
//  }
package stockchart

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Brand struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Brands []Brand

type Map struct {
	Key   string
	Value string
}

type Maps []Map

type Client struct {
	Header string
	Body   map[string]string
}

type Data struct {
	Id   string
	Name string
	Data [][]string
}

func (m Maps) Len() int {
	return len(m)
}

func (m Maps) Less(i, j int) bool {
	return m[i].Key < m[j].Key
}

func (m Maps) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func getHtml() string {
	return `<html>
	<head>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/dygraph/1.2/dygraph-combined.min.js"></script>
	</head>
	<body>
		<div id="graphdiv"></div>
		<script>
			var csv = "";
			new Dygraph(document.getElementById("graphdiv"),csv);
		</script>
	</body>
</html>
`
}

func scrape(url string) ([][]string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	dat := [][]string{}
	doc.Find("div.padT12 table tbody tr").Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			t, _ := time.Parse("2006年1月2日", s.Find("td").First().Text())
			dat = append(dat, []string{
				t.Format("2006/01/02"),
				strings.Replace(s.Find("td").Last().Text(), ",", "", 1),
			})
		}
	})
	return dat, nil
}

func getData(b Brand) (*Data, error) {
	datCh := make(chan [][]string)
	errCh := make(chan error)
	data := [][]string{}
	for i := 1; i <= 4; i++ {
		go func(i int) {
			url := "http://info.finance.yahoo.co.jp/history/?code=" +
				b.Id + "&p=" + strconv.Itoa(i)
			dat, err := scrape(url)
			if err != nil {
				errCh <- err
			} else {
				datCh <- dat
			}
		}(i)
	}
	for i := 1; i <= 4; i++ {
		select {
		case err := <-errCh:
			return nil, err
		case dat := <-datCh:
			data = append(data, dat...)
		}
	}
	return &Data{
		Id:   b.Id,
		Name: b.Name,
		Data: data,
	}, nil
}

func New(params string) (*Client, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	p := Brands{}
	json.Unmarshal([]byte(params), &p)
	header := ""
	body := map[string]string{}
	datCh := make(chan *Data)
	errCh := make(chan error)
	for _, b := range p {
		go func(b Brand) {
			dat, err := getData(b)
			if err != nil {
				errCh <- err
			} else {
				datCh <- dat
			}
		}(b)
	}
	for i := 0; i < len(p); i++ {
		select {
		case err := <-errCh:
			return nil, err
		case dat := <-datCh:
			header += "," + dat.Name
			for _, d := range dat.Data {
				body[d[0]] += "," + d[1]
			}
		}
	}
	return &Client{
		Header: header,
		Body:   body,
	}, nil
}

func (c *Client) Csv() string {
	r := "Date" + c.Header + "\\n"
	ms := Maps{}
	for k, v := range c.Body {
		ms = append(ms, Map{Key: k, Value: v})
	}
	sort.Sort(ms)
	for _, v := range ms {
		r += v.Key + v.Value + "\\n"
	}
	return r
}

func (c *Client) Html() string {
	r := strings.Replace(getHtml(),
		"csv = \"\"",
		"csv = \""+c.Csv()+"\"", 1)
	return r
}
