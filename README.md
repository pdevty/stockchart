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
	"github.com/pdevty/trustchart"
	"io/ioutil"
	"os"
)

func main() {

	// params json format
	// term 1y (1 year ) 2y 3y ...
	//      1m (1 month) 2m 3m ...
	//      1d (1 day  ) 2d 3d ...
	// brand id,name from yahoo finance japan
	params := `{
		"term":"1y",
		"brands":[
			{"id":"89311067","name":"jrevive"},
			{"id":"29311041","name":"ﾆｯｾｲ日経225"},
			{"id":"03316042","name":"健次"},
			{"id":"2931113C","name":"ﾆｯｾｲ外国株"}
		]
	}`

	// new
	tc, err := trustchart.New(params)
	if err != nil {
		panic(err)
	}

	// return chart csv
	fmt.Println(tc.Csv())

	// return chart html
	fmt.Println(tc.Html())

	// create chart html file
	ioutil.WriteFile("index.html",
		[]byte(tc.Html()), os.ModePerm)
}
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
