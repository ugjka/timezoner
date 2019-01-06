package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/ugjka/go-tz.v2/tz"
)

var t = template.New("root")

func main() {
	port := flag.Uint("port", 40001, "listener port")
	flag.Parse()
	if *port > 65535 {
		os.Stderr.WriteString("Error: invalid port number\n")
		return
	}
	log.SetPrefix("timezoner ")

	funcMap := template.FuncMap{
		"FormatDate": func(value time.Time) string {
			return value.Format("2006-01-02 15:04:05")
		},
	}
	t = t.Funcs(funcMap)
	var err error
	t, err = t.Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", root)
	mux.HandleFunc("/api/json", api)
	log.Printf("listening on :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func root(w http.ResponseWriter, r *http.Request) {
	var err error
	var lon, lat float64
	values := r.URL.Query()

	switch {
	case len(values) < 2:
		lon, lat = 24.1051846, 56.9493977
	case len(values) >= 2:
		lon, lat, err = getValues(values)
		if err != nil {
			err = fmt.Errorf("invalid input")
			break
		}
	}

	var info *Info
	if err == nil {
		info, err = getInfo(&tz.Point{Lon: lon, Lat: lat})
	}

	res := struct {
		Err  error
		Data *Info
	}{
		Err:  err,
		Data: info,
	}
	t.Execute(w, res)
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	values := r.URL.Query()
	var lon, lat float64
	_, err := fmt.Sscanf(values.Get("lon"), "%f", &lon)
	if err != nil {
		v := Info{Status: "invalid lon", StatusCode: http.StatusBadRequest}
		w.WriteHeader(http.StatusBadRequest)
		encode(w, v)
		return
	}
	_, err = fmt.Sscanf(values.Get("lat"), "%f", &lat)
	if err != nil {
		v := Info{Status: "invalid lat", StatusCode: http.StatusBadRequest}
		w.WriteHeader(http.StatusBadRequest)
		encode(w, v)
		return
	}
	zones, err := tz.GetZone(tz.Point{Lon: lon, Lat: lat})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		v := Info{Status: fmt.Sprintf("%v", err), StatusCode: http.StatusBadRequest}
		encode(w, v)
		return
	}
	res := Info{
		Status:     "ok",
		StatusCode: http.StatusOK,
		Lon:        lon,
		Lat:        lat,
	}

	for _, v := range zones {
		loc, err := time.LoadLocation(v)
		if err != nil {
			log.Printf("could not load time info: %s", v)
			v := Info{Status: "internal error", StatusCode: http.StatusInternalServerError}
			w.WriteHeader(http.StatusInternalServerError)
			encode(w, v)
			return
		}
		localTime := time.Now().In(loc)
		name, offset := localTime.Zone()
		res.Info = append(res.Info, Item{
			TZID:   v,
			Time:   localTime,
			Offset: offset,
			Name:   name},
		)
	}
	w.WriteHeader(http.StatusOK)
	encode(w, res)
}

func encode(w io.Writer, v interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	err := enc.Encode(v)
	if err != nil {
		log.Printf("could not encode data :%v", err)
	}
}

// Info struct
type Info struct {
	Status     string
	StatusCode int
	Lon        float64
	Lat        float64
	Info       []Item `json:",omitempty"`
}

// Item is a tzid item
type Item struct {
	Time   time.Time
	TZID   string
	Offset int
	Name   string
}

func getValues(values url.Values) (lon, lat float64, err error) {
	_, err = fmt.Sscanf(values.Get("lon"), "%f", &lon)
	if err != nil {
		return
	}
	_, err = fmt.Sscanf(values.Get("lat"), "%f", &lat)
	if err != nil {
		return
	}
	return
}

func getInfo(p *tz.Point) (info *Info, err error) {
	zones, err := tz.GetZone(*p)
	if err != nil {
		return nil, err
	}
	info = &Info{
		Lon: p.Lon,
		Lat: p.Lat,
	}
	for _, v := range zones {
		loc, err := time.LoadLocation(v)
		if err != nil {
			log.Printf("could not load time info: %s", v)
			return nil, err
		}
		localTime := time.Now().In(loc)
		name, offset := localTime.Zone()
		info.Info = append(info.Info, Item{
			TZID:   v,
			Time:   localTime,
			Offset: offset,
			Name:   name},
		)
	}
	return info, nil
}
