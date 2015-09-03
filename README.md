# stockchart [![GoDoc](https://godoc.org/github.com/pdevty/stockchart?status.svg)](https://godoc.org/github.com/pdevty/stockchart)

stock chart from [yahoo finance japan](http://apl.morningstar.co.jp/webasp/yahoo-fund/fund/index.aspx)

![chart](https://github.com/pdevty/stockchart/blob/master/chart.png)

## Installation

    $ go get github.com/pdevty/stockchart

## Usage

```go
package main

import (
	"fmt"
	"github.com/pdevty/stockchart"
	"io/ioutil"
	"os"
)

func main() {
	// params json format
	// brand id,name from yahoo finance japan
	params := `[
		{"id":"7618","name":"ピーシーデポ"},
		{"id":"6157","name":"日進工具"},
		{"id":"7821","name":"前田工繊"},
		{"id":"7917","name":"藤森工業"},
		{"id":"4681","name":"リゾートトラスト"},
		{"id":"4301","name":"アミューズ"},
		{"id":"4290","name":"プレステージ"},
		{"id":"2780","name":"コメ兵"},
		{"id":"2695","name":"くらコーポレーション"},
		{"id":"2695","name":"内外トランスライン"}
	]`
	// new
	sc, err := stockchart.New(params)
	if err != nil {
		panic(err)
	}
	// return chart html
	fmt.Println(sc.Html())
	// return chart csv
	fmt.Println(sc.Csv())
	// create chart html file
	ioutil.WriteFile("index.html",
		[]byte(sc.Html()), os.ModePerm)
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
