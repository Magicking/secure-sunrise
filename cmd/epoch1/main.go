package main

import (
	"bytes"
	"fmt"
	"github.com/dinedal/astrotime"
	"github.com/oschwald/geoip2-golang"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const LATITUDE = float64(38.8895)
const LONGITUDE = float64(77.0352)
const entrypoint = string("http://www.insecam.org")

func GetListCamByTimeZone(tz string, page uint) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/en/bytimezone/%s/?page=%d", entrypoint, tz, page)
	req, err := http.NewRequest("GET", url, nil)
	// ...
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.96 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type Cam struct {
	URL    *url.URL
	Record *geoip2.City
}

func FilterCams(b []byte) ([]*url.URL, error) {
	thumbPrefix := []byte(`thumbnail-item__img img-responsive" src="`)
	thumbPostfix := []byte(`"`)

	var cams []*url.URL
	for i := bytes.Index(b, thumbPrefix); i > 0; i = bytes.Index(b, thumbPrefix) {
		b = b[i+len(thumbPrefix):]
		sep := bytes.Index(b, thumbPostfix)
		camUrl, err := url.Parse(html.UnescapeString(string(b[:sep])))
		if err != nil {
			return nil, err
		}
		cams = append(cams, camUrl)
		b = b[sep+len(thumbPostfix):]
	}

	return cams, nil
}

func enOrFirst(m map[string]string) string {
	ret, ok := m["en"]
	if ok {
		return ret
	}
	for _, city := range m {
		return city
	}
	return "Unknown city"
}

func NewCam(db *geoip2.Reader, url *url.URL) (*Cam, error) {
	h := strings.Split(url.Host, ":")
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(h[0])
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}
	return &Cam{
		URL:    url,
		Record: record,
	}, nil
}

func main() {
	db, err := geoip2.Open("/home/magicking/source/gocode/src/github.com/Magicking/secure-sunrise/GeoLite2-City_20170606/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	b, err := GetListCamByTimeZone("+01:00", 1)
	if err != nil {
		log.Fatal(err)
	}
	camsURLs, err := FilterCams(b)
	if err != nil {
		log.Fatal(err)
	}

	var cams []*Cam
	for i := range camsURLs {
		cam, err := NewCam(db, camsURLs[i])
		if err == nil {
			cams = append(cams, cam)
		}
	}
	for _, v := range cams {
		t := astrotime.NextSunrise(time.Now(), v.Record.Location.Latitude, v.Record.Location.Longitude)
		tzname, _ := t.Zone()
		fmt.Printf("The next sunrise at %q is %d:%02d %s on %d/%d/%d %s.\n",
			enOrFirst(v.Record.City.Names),
			t.Hour(), t.Minute(), tzname, t.Month(), t.Day(), t.Year(),
			v.URL,
		)
	}
	//fmt.Println(string(b))
}
