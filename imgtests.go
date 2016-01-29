package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/gographics/imagick/imagick"
	"io/ioutil"
	"os"
	"strconv"
)

type Page struct {
	Title string
	Body  []byte
}

func main() {
	fmt.Println("Imagick start")

	images, _ := ioutil.ReadDir("./images")
	proceedImages(images)

	images2, _ := ioutil.ReadDir("./imagesDone")
	proceedHtml(images2)

}

func proceedImages(i []os.FileInfo) {
	imagick.Initialize()
	defer imagick.Terminate()

	green := imagick.NewPixelWand()
	green.SetColor("#1cd000")
	none := imagick.NewPixelWand()
	none.SetColor("none")
	channel := imagick.CHANNEL_OPACITY
	for _, f := range i {
		importImage := imagick.NewMagickWand()
		importImage.ReadImage("./images/" + f.Name())
		importImage.FloodfillPaintImage(channel, none, 20000, green, 0, 0, false)
		importImage.WriteImage("imagesDone/" + f.Name())
		importImage.Destroy()
	}
	green.Destroy()
	none.Destroy()
}

func proceedHtml(i []os.FileInfo) {
	fps := 17
	count := len(i)

	_ = count

	page := &Page{}
	page.Title = "animation.html"
	outPut := `
		<svg version="1.1" baseProfile="tiny" id="svg-root"''
	 		width="100%" height="100%" viewBox="0 0 716 578"
	  		xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">'
	`

	for key, f := range i {
		keyWrite := strconv.FormatInt(int64(key), 10)

		imgFile, err := os.Open("imagesDone/" + f.Name())

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fInfo, _ := imgFile.Stat()
		var size int64 = fInfo.Size()
		buf := make([]byte, size)

		// read file content into buffer
		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		imgBase64Str := base64.StdEncoding.EncodeToString(buf)

		data := imgBase64Str
		imgFile.Close()

		outPut += `<image id="frame` + keyWrite + `" width="716" height="578" xlink:href="data:image/png;base64,` + data + `" display="inline">8
			<set id="show` + keyWrite + `" attributeName="display" to="inline" begin="<?=($i==0?"0s;":"")?>show<?=($i+$numFrames-1)%$numFrames?>.end" dur="` + string(1/fps) + `s" fill="freeze"/>
			<set id="hide` + keyWrite + `" attributeName="display" to="none"  begin="show` + keyWrite + `.end" dur="0.01s" fill="freeze"/>      
			</image>
		`
	}

	outPut += `</svg>`

	page.Body = []byte(outPut)
	page.save()
}

func (p *Page) save() error {
	filename := p.Title
	return ioutil.WriteFile(filename, p.Body, 0600)
}
