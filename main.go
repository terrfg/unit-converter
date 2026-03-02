package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

//var tmpl1 = template.Must(template.New("page").Parse(`
//<!DOCTYPE html>
//<html>
//<head>
//<meta charset="UTF-8">
//<title>Unit Converter</title>
//<style>
//body {
//	font-family: Arial, sans-serif;
//	background: #f2f2f2;
//}
//
//.wrapper {
//	display: flex;
//	gap: 40px;
//	margin: 40px;
//}
//
//.card {
//	background: white;
//	border: 3px solid #222;
//	border-radius: 15px;
//	padding: 30px;
//	width: 400px;
//}
//
//.tabs span {
//	margin-right: 20px;
//	font-weight: bold;
//	cursor: pointer;
//}
//
//.active {
//	color: #4169E1;
//	border-bottom: 3px solid #4169E1;
//}
//
//input, select {
//	width: 100%;
//	padding: 10px;
//	margin: 10px 0;
//	border-radius: 8px;
//	border: 2px solid #333;
//}
//
//button {
//	padding: 10px 20px;
//	border-radius: 8px;
//	border: 2px solid #333;
//	background: white;
//	cursor: pointer;
//}
//
//.result-text {
//	font-size: 28px;
//	margin: 20px 0;
//	font-weight: bold;
//}
//</style>
//</head>
//
//<body>
//
//<div class="wrapper">
//
//<!-- Left Card -->
//<div class="card">
//<h2>Unit Converter</h2>
//
//<form method="POST">
//<div class="tabs">
//	<span class="{{if eq .Category "length"}}active{{end}}">Length</span>
//	<span class="{{if eq .Category "weight"}}active{{end}}">Weight</span>
//	<span class="{{if eq .Category "temperature"}}active{{end}}">Temperature</span>
//</div>
//
//<input type="hidden" name="category" value="{{.Category}}">
//
//<label>Enter the value</label>
//<input type="number" step="any" name="value" required>
//
//<label>Unit to convert from</label>
//<input type="text" name="from" placeholder="meter">
//
//<label>Unit to convert to</label>
//<input type="text" name="to" placeholder="kilometer">
//
//<button type="submit">Convert</button>
//</form>
//</div>
//
//<!-- Right Card -->
//<div class="card">
//<h2>Unit Converter</h2>
//
//<div class="tabs">
//	<span class="{{if eq .Category "length"}}active{{end}}">Length</span>
//	<span class="{{if eq .Category "weight"}}active{{end}}">Weight</span>
//	<span class="{{if eq .Category "temperature"}}active{{end}}">Temperature</span>
//</div>
//
//<h3>Result of your calculation</h3>
//
//{{if .Result}}
//<div class="result-text">{{.Result}}</div>
//<form method="GET">
//<button>Reset</button>
//</form>
//{{end}}
//
//</div>
//
//</div>
//
//</body>
//</html>
//`))

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Unit Converter</title>
<style>
body { font-family: Arial; margin: 40px; }
input, select, button { margin: 5px 0; padding: 5px; }
</style>
</head>
<body>
<h2>Unit Converter</h2>

<form method="POST">
<label>Category:</label><br>
<select name="category">
	<option value="length">Length</option>
	<option value="weight">Weight</option>
	<option value="temperature">Temperature</option>
</select><br>

<label>Value:</label><br>
<input type="number" step="any" name="value" required><br>

<label>From:</label><br>
<input type="text" name="from" placeholder="e.g. meter"><br>

<label>To:</label><br>
<input type="text" name="to" placeholder="e.g. kilometer"><br>

<button type="submit">Convert</button>
</form>

{{if .Result}}
<h3>Result: {{.Result}}</h3>
{{end}}

</body>
</html>
`))

var lengthUnits = map[string]float64{
	"millimeter": 0.001,
	"centimeter": 0.01,
	"meter":      1,
	"kilometer":  1000,
	"inch":       0.0254,
	"foot":       0.3048,
	"yard":       0.9144,
	"mile":       1609.34,
}

var weightUnits = map[string]float64{
	"milligram": 0.001,
	"gram":      1,
	"kilogram":  1000,
	"ounce":     28.3495,
	"pound":     453.592,
}

func convertTemperature(value float64, from, to string) float64 {
	var res float64
	switch from {
	case "Celsius":
		res = value
	case "Fahrenheit":
		res = (value - 32) * 5 / 9
	case "Kelvin":
		res = value - 273.15
	}
	switch to {
	case "Celsius":
		return res
	case "Fahrenheit":
		return res*9/5 + 32
	case "Kelvin":
		return res + 273.15
	}
	return 0
}

func handle(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Result string
	}{}
	if r.Method == http.MethodPost {
		r.ParseForm()
		category := r.Form.Get("category")
		valueStr := r.Form.Get("value")
		from := r.Form.Get("from")
		to := r.Form.Get("to")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err == nil {
			var result float64
			switch category {
			case "length":
				base := value * lengthUnits[from]
				result = base / lengthUnits[to]
			case "weight":
				base := value * weightUnits[from]
				result = base / weightUnits[to]
			case "temperature":
				result = convertTemperature(value, from, to)
			}
			data.Result = fmt.Sprintf("%.4f %s = %.4f %s", value, from, result, to)
		}
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handle)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
