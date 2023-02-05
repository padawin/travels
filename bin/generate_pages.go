package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"text/template"
)

// Declare type pointer to a template
var (
	funcMap = template.FuncMap{
		"url": createURL,
	}
	randomImageTemplate   = template.Must(template.New("random.php").Funcs(funcMap).ParseFiles("templates/random.php"))
	homepageTemplate      = template.Must(template.New("index.html").Funcs(funcMap).ParseFiles("templates/header.html", "templates/footer.html", "templates/index.html"))
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

type AlbumTravel struct {
	Pictures [][]string `json:"pics"`
	Places   []string   `json:"places"`
	Title    string     `json:"title"`
}

type Albums struct {
	Latest  []string               `json:"latest"`
	Travels map[string]AlbumTravel `json:"travels"`
}

func (a *Albums) getTravelsOrdered() []AlbumTravel {
	var res []AlbumTravel
	for _, travel := range a.Travels {
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
	Travel string
	Places []string
}

type picturesListData struct {
	Travel   string
	Place    string
	Pictures []string
	Places   []string
}

type picturePageData struct {
	Travel   string
	Place    string
	Picture  string
	Previous *int
	Next     *int
	Places   []string
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
	latestPicturesData := picturesListData{
		Travel:   "Latest",
		Pictures: albums.Latest,
	}
	compileTemplate(picturesTemplate, destDir+"/latest.html", latestPicturesData)
	prepareLatestImages(destDir, albums)
	compileRandomLatest(destDir)

	// Travels
	processTravels(destDir, albums.Travels)
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

func prepareLatestImages(destDir string, albums Albums) {
	// delete previous __latest__
	err := os.RemoveAll(destDir + "/images/100x100x1/__latest__")
	if err != nil {
		log.Fatal(err)
	}
	// create 100x100x1/latest
	err = os.MkdirAll(destDir+"/images/100x100x1/__latest__", 0755)
	if err != nil {
		log.Fatal(err)
	}
	// copy latest images in directory
	for i, img := range albums.Latest {
		err = os.Symlink("../"+img, fmt.Sprintf("%s/images/100x100x1/__latest__/latest_%d.jpg", destDir, i))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func compileRandomLatest(destDir string) {
	compileTemplate(
		randomImageTemplate,
		destDir+"/images/100x100x1/__latest__/random.php",
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
}

func compilePlacesList(destDir string, travel AlbumTravel) {
	compileTemplate(
		placesTemplate,
		destDir+"/travels/"+travel.Title+"/index.html",
		placesListData{
			Travel: travel.Title,
			Places: travel.Places,
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
	destFile = destDir + "/travels/" + travel.Title + "/" + placePath + "/index.html"
	compileTemplate(
		picturesTemplate,
		destFile,
		picturesListData{
			Travel:   travel.Title,
			Place:    place,
			Pictures: travel.Pictures[placeIndex],
			Places:   travel.Places,
		},
	)
	for i, p := range travel.Pictures[placeIndex] {
		destFile = destDir + "/travels/" + travel.Title + "/" + placePath + "/" + strconv.Itoa(i) + ".html"
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
				Travel:   travel.Title,
				Place:    place,
				Picture:  p,
				Places:   travel.Places,
				Previous: previous,
				Next:     next,
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
