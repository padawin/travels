package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

// Declare type pointer to a template
var (
	funcMap = template.FuncMap{
		"url": createURL,
	}
	randomImageTemplate   = template.Must(template.New("random.php").Funcs(funcMap).ParseFiles("templates/random.php"))
	homepageTemplate      = template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/index.html"))
	latestTemplate        = template.Must(template.New("latest.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/latest.html"))
	latestPictureTemplate = template.Must(template.New("latest-picture.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/latest-picture.html"))
	placesTemplate        = template.Must(template.New("places-list.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/places-list.html"))
	picturesTemplate      = template.Must(template.New("pictures-list.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/pictures-list.html"))
	pictureTemplate       = template.Must(template.New("picture.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/picture.html"))

	force bool
)

func createURL(components ...interface{}) string {
	var res string
	for _, component := range components {
		if c, ok := component.(string); ok {
			if c != "" {
				res += "/" + c
			}
		} else if c, ok := component.(int); ok {
			res += "/" + strconv.Itoa(c)
		} else if c, ok := component.(*int); ok {
			if c != nil {
				res += "/" + strconv.Itoa(*c)
			}
		}
	}
	return res
}

func getTitle(travel, place string) string {
	title := travel
	if place != "" {
		title += " - " + place
	}
	return title
}

type AlbumTravel struct {
	Pictures     [][]string `json:"pics"`
	Places       []string   `json:"places"`
	Title        string     `json:"title"`
	EncodedTitle string
}

type Albums struct {
	Latest  []string               `json:"latest"`
	Travels map[string]AlbumTravel `json:"travels"`
}

func (a *Albums) getTravelsOrdered() []AlbumTravel {
	var res []AlbumTravel
	for _, travel := range a.Travels {
		travel.EncodedTitle = url.PathEscape(travel.Title)
		res = append(res, travel)
	}
	sort.Slice(res, func(i, j int) bool {
		iIsDate := res[i].Title[0] == '2'
		jIsDate := res[j].Title[0] == '2'
		if iIsDate && jIsDate {
			return res[i].Title > res[j].Title
		} else if !iIsDate && !jIsDate {
			return res[i].Title < res[j].Title
		} else {
			return iIsDate && !jIsDate
		}
	})
	return res
}

type placesListData struct {
	MetaTitle string
	MetaUrl   string
	MetaImage string
	Travel    string
	Places    []string
}

type picturesListData struct {
	MetaTitle string
	MetaUrl   string
	MetaImage string
	Travel    string
	Place     string
	Pictures  []string
	Places    []string
}

type picturePageData struct {
	MetaTitle string
	MetaUrl   string
	MetaImage string
	Travel    string
	Place     string
	Picture   string
	Previous  *int
	Next      *int
	Places    []string
}

type coordinate struct {
	X int
	Y int
}

type randomPictureData struct {
	DestWidth      int
	DestHeight     int
	CountThumbs    int
	ThumbWidth     int
	ThumbHeight    int
	ThumbPositions []coordinate
}

func main() {
	var jsonFile string = os.Args[1]
	var destDir string = os.Args[2]

	force = os.Getenv("FORCE") != ""

	albums := getAlbums(jsonFile)

	// Home page
	compileTemplate(homepageTemplate, destDir+"/index.html", albums.getTravelsOrdered())

	// Latest
	processLatestPages(destDir, albums.Latest)
	prepareLatestImages(destDir, albums)
	compileRandomLatest(destDir)

	// Travels
	processTravels(destDir, albums.Travels)
}

func processLatestPages(destDir string, pictures []string) {
	latestPicturesData := picturesListData{
		MetaTitle: "Latest pictures",
		MetaUrl:   "/latest/index.html",
		MetaImage: "/pictures/512x269x1/latest/random.php",
		Pictures:  pictures,
	}
	err := os.MkdirAll(destDir+"/latest", 0755)
	if err != nil {
		log.Fatal(err)
	}
	compileTemplate(latestTemplate, destDir+"/latest/index.html", latestPicturesData)
	for i, p := range pictures {
		url := "/latest/" + strconv.Itoa(i) + ".html"
		destFile := destDir + url
		var previous, next *int
		if i < len(pictures)-1 {
			n := i + 1
			next = &n
		}
		if i > 0 {
			p := i - 1
			previous = &p
		}
		parts := strings.Split(p, "/")
		travel := parts[0]
		place := ""
		if len(parts) == 3 {
			place = parts[1]
		}
		compileTemplate(
			latestPictureTemplate,
			destFile,
			picturePageData{
				MetaTitle: getTitle(travel, place),
				MetaUrl:   url,
				MetaImage: "/pictures/512x269x1/" + p,
				Travel:    travel,
				Place:     place,
				Picture:   p,
				Previous:  previous,
				Next:      next,
			},
		)
	}
}

func prepareLatestImages(destDir string, albums Albums) {
	for _, resolution := range []string{"100x100x1", "512x269x1"} {
		// delete previous Latest
		err := os.RemoveAll(destDir + "/images/" + resolution + "/latest")
		if err != nil {
			log.Fatal(err)
		}
		// create 100x100x1/latest
		err = os.MkdirAll(destDir+"/images/"+resolution+"/latest", 0755)
		if err != nil {
			log.Fatal(err)
		}
		// copy latest images in directory
		for i, img := range albums.Latest {
			err = os.Symlink("../"+img, fmt.Sprintf("%s/images/"+resolution+"/latest/latest_%d.jpg", destDir, i))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func processTravels(destDir string, travels map[string]AlbumTravel) {
	for _, travel := range travels {
		err := os.MkdirAll(destDir+"/travels/"+travel.Title, 0755)
		if err != nil {
			log.Fatal(err)
		}

		compileRandomHomepage(destDir, travel)
		if len(travel.Places) == 0 {
			compilePicturesList(destDir, travel, "", 0)
		} else {
			compilePlacesList(destDir, travel)
			for i, place := range travel.Places {
				err := os.MkdirAll(destDir+"/travels/"+travel.Title+"/"+place, 0755)
				if err != nil {
					log.Fatal(err)
				}
				compileRandomPlace(destDir, travel, place)
				compilePicturesList(destDir, travel, place, i)
			}
		}
	}
}

func compileRandomLatest(destDir string) {
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/100x100x1/latest/random.php",
		randomPictureData{
			DestWidth:      400,
			DestHeight:     100,
			CountThumbs:    4,
			ThumbWidth:     100,
			ThumbHeight:    100,
			ThumbPositions: []coordinate{{0, 0}, {100, 0}, {200, 0}, {300, 0}},
		},
	)
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/512x269x1/latest/random.php",
		randomPictureData{
			DestWidth:      512,
			DestHeight:     269,
			CountThumbs:    1,
			ThumbWidth:     512,
			ThumbHeight:    269,
			ThumbPositions: []coordinate{{0, 0}},
		},
	)
}

func compileRandomHomepage(destDir string, travel AlbumTravel) {
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/100x100x1/"+travel.Title+"/random.php",
		randomPictureData{
			DestWidth:      200,
			DestHeight:     200,
			CountThumbs:    4,
			ThumbWidth:     100,
			ThumbHeight:    100,
			ThumbPositions: []coordinate{{0, 0}, {100, 0}, {0, 100}, {100, 100}},
		},
	)
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/100x100x1/"+travel.Title+"/random-mobile.php",
		randomPictureData{
			DestWidth:      400,
			DestHeight:     100,
			CountThumbs:    4,
			ThumbWidth:     100,
			ThumbHeight:    100,
			ThumbPositions: []coordinate{{0, 0}, {100, 0}, {200, 0}, {300, 0}},
		},
	)
}

func compilePlacesList(destDir string, travel AlbumTravel) {
	url := "/travels/" + travel.Title + "/index.html"
	compileTemplate(
		placesTemplate,
		destDir+url,
		placesListData{
			MetaTitle: getTitle(travel.Title, ""),
			MetaUrl:   url,
			MetaImage: "/pictures/512x269x1/" + travel.Title + "/random.php",
			Travel:    travel.Title,
			Places:    travel.Places,
		},
	)
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/512x269x1/"+travel.Title+"/random.php",
		randomPictureData{
			DestWidth:      512,
			DestHeight:     269,
			CountThumbs:    1,
			ThumbWidth:     512,
			ThumbHeight:    269,
			ThumbPositions: []coordinate{{0, 0}},
		},
	)
}

func compileRandomPlace(destDir string, travel AlbumTravel, place string) {
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/118x133x1/"+travel.Title+"/"+place+"/random.php",
		randomPictureData{
			DestWidth:      118,
			DestHeight:     133,
			CountThumbs:    1,
			ThumbWidth:     118,
			ThumbHeight:    133,
			ThumbPositions: []coordinate{{0, 0}},
		},
	)
}

func compilePicturesList(destDir string, travel AlbumTravel, place string, placeIndex int) {
	var destFile string
	var placePath string = place
	if place != "" {
		placePath = "/" + place
	}
	url := "/travels/" + travel.Title + placePath + "/index.html"
	destFile = destDir + url
	compileTemplate(
		picturesTemplate,
		destFile,
		picturesListData{
			MetaTitle: getTitle(travel.Title, place),
			MetaUrl:   url,
			MetaImage: "/pictures/512x269x1/" + travel.Title + placePath + "/random.php",
			Travel:    travel.Title,
			Place:     place,
			Pictures:  travel.Pictures[placeIndex],
			Places:    travel.Places,
		},
	)
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/512x269x1/"+travel.Title+placePath+"/random.php",
		randomPictureData{
			DestWidth:      512,
			DestHeight:     269,
			CountThumbs:    1,
			ThumbWidth:     512,
			ThumbHeight:    269,
			ThumbPositions: []coordinate{{0, 0}},
		},
	)
	for i, p := range travel.Pictures[placeIndex] {
		url = "/travels/" + travel.Title + placePath + "/" + strconv.Itoa(i) + ".html"
		destFile = destDir + url
		var previous, next *int
		if i < len(travel.Pictures[placeIndex])-1 {
			n := i + 1
			next = &n
		}
		if i > 0 {
			p := i - 1
			previous = &p
		}
		compileTemplate(
			pictureTemplate,
			destFile,
			picturePageData{
				MetaTitle: getTitle(travel.Title, place),
				MetaUrl:   url,
				MetaImage: "/pictures/512x269x1/" + p,
				Travel:    travel.Title,
				Place:     place,
				Picture:   p,
				Places:    travel.Places,
				Previous:  previous,
				Next:      next,
			},
		)
	}
}

func getAlbums(jsonFile string) Albums {
	fileContent, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}

	defer fileContent.Close()

	byteResult, err := ioutil.ReadAll(fileContent)
	if err != nil {
		panic(err)
	}

	var albums Albums
	err = json.Unmarshal(byteResult, &albums)
	if err != nil {
		panic(err)
	}

	return albums
}

func compileTemplate(t *template.Template, destFile string, data interface{}) {
	f, err := os.Open(destFile)
	if err == nil && !force {
		f.Close()
		return
	}

	f, err = os.Create(destFile)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	err = t.Funcs(funcMap).Execute(f, data)
	if err != nil {
		panic(err)
	}
}
