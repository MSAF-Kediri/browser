package browser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/tebeka/selenium"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Browser struct {
	Wd selenium.WebDriver
	LogFile *os.File
}

func (b *Browser) Init(port int) {
	// Connect to the WebDriver instance running locally.
	//coba := map[string]interface{}{
	//	"goog:chromeOptions": {
	//		"excludeSwitches":        [ "enable-automation" ],
	//		"useAutomationExtension": false
	//	}

	coba := map[string]interface{}{
		"useAutomationExtension": false,
		"excludeSwitches": []string{
			"enable-automation",
		},
	}
	caps := selenium.Capabilities{
		"goog:chromeOptions": coba,
		"browserName":        "chrome",
	}

	b.Wd, _ = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	sessionID := b.Wd.SessionID()
	log.SetOutput(b.LogFile)
	b.writePrintln("Chrome driver sessionID: " + sessionID)

}

func (b *Browser) writePrintln(s interface{}) {
	fmt.Println(time.Now().Format("02-01-2006 15:04:05"), s)
	log.Println(time.Now().Format("02-01-2006 15:04:05"), s)
}

// ElementIsVisible returns a condition that checks if the element is visible.
func (b *Browser) ElementIsVisible(elt selenium.WebElement) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		visible, err := elt.IsDisplayed()
		return visible, err
	}
}

// ElementIsEnabled returns a condition that checks if element's enabled.
func (b *Browser) ElementIsEnabled(elt selenium.WebElement) selenium.Condition {
	return func (wd selenium.WebDriver) (bool, error) {
		enabled, err := elt.IsEnabled()
		return enabled, err
	}
}

// ElementIsLocatedAndVisible returns a condition that checks if the element is found on page and is visible.
func (b *Browser) ElementIsLocatedAndVisible(by, selector string) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		element, err := wd.FindElement(by, selector)
		if err != nil {
			return false, nil
		}
		visible, err := element.IsDisplayed()
		return visible, err
	}
}

func (b *Browser) InputText() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	b.writePrintln("inputtext: " + text)
	return text
}

func (b *Browser) SendKeys(by string, value string, keys string) bool {
	b.writePrintln("sendkeys dengan" + " " + by + " " + value)
	inputan, err := b.Wd.FindElement(by, value)
	if err == nil {
		err = b.Wd.WaitWithTimeout(b.ElementIsVisible(inputan), 10 * time.Second)
		if err == nil {
			_ = inputan.Clear()
			err = inputan.SendKeys(keys)
			if err == nil {
				valueInputan, _ := inputan.GetAttribute("value")
				if valueInputan == keys {
					b.writePrintln("sendkeys dengan" + " " + by + " " + value + " berhasil")
					return true
				}
			}
		}
	}
	b.writePrintln("sendkeys dengan" + " " + by + " " + value + " gagal")
	return false
}

func (b *Browser) Click(by string, value string) bool {
	b.writePrintln("klik tombol dengan" + " " + by + " " + value)
	btn, err := b.Wd.FindElement(by, value)
	if err == nil {
		err = b.Wd.WaitWithTimeout(b.ElementIsEnabled(btn), 10 * time.Second)
		if err == nil {
			err = btn.Click()
			if err == nil {
				b.writePrintln("klik tombol dengan" + " " + by + " " + value + " berhasil")
				return true
			}
		}
	}
	b.writePrintln("klik tombol dengan" + " " + by + " " + value + " gagal")
	return false
}

func (b * Browser) SwitchDefaultFrame() bool {
	b.writePrintln("kembali ke frame awal")
	err := b.Wd.SwitchFrame(nil)
	if err == nil {
		b.writePrintln("kembali ke frame awal berhasil")
		return true
	}
	b.writePrintln("kembali ke frame awal gagal")
	return false
}

func (b *Browser) SwitchFrame(by string, value string) bool {
	b.writePrintln("pindah ke frame dengan " + by + " " + value)
	frame, err := b.Wd.FindElement(by, value)
	if err == nil {
		err = b.Wd.SwitchFrame(frame)
		if err == nil {
			b.writePrintln("pindah ke frame dengan " + by + " " + value + " berhasil")
			return true
		}
	}
	b.writePrintln("pindah ke frame dengan " + by + " " + value + " gagal")
	return false
}

func (b *Browser) serveFrames(imgByte []byte) {
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	if err != nil {
		log.Fatalln(err)
	}

	out, _ := os.Create("img.jpeg")
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)
	//jpeg.Encode(out, img, nil)
	if err != nil {
		log.Println(err)
	}

}

func (b *Browser) ScreenShot(x0, y0, x1, y1, namefile interface{}) {
	imgByte, err := b.Wd.Screenshot()
	if err != nil {
		b.writePrintln(err)
	}
	err = ioutil.WriteFile(namefile.(string), imgByte, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if x0 == nil && y0 == nil && x1 == nil && y1 == nil {
		return
	}

	fImg1, _ := os.Open(namefile.(string))
	defer fImg1.Close()
	img1, _, _ := image.Decode(fImg1)

	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	sImg, ok := img1.(subImager)
	if !ok {
		b.writePrintln("image does not support cropping")
	}

	crop := image.Rect(x0.(int), y0.(int), x1.(int), y1.(int))
	cropImg := sImg.SubImage(crop)

	namecropImg := "crop_" + namefile.(string)
	fd, err := os.Create(namecropImg)
	if err != nil {
		b.writePrintln(err)
	}
	defer fd.Close()

	err = png.Encode(fd, cropImg)
	if err != nil {
		b.writePrintln(err)
	}
}

func (b *Browser) SelectOptions(by string, value string, attr string, target string) bool{
	b.writePrintln("select options dengan " + by + " " + value + " " + attr + " " + target)
	opts, err := b.Wd.FindElements(by, value)
	if err == nil {
		for _, opt := range opts {
			text := ""
			if attr == "value" {
				text, _ = opt.GetAttribute(attr)
			} else {
				text, _ = opt.Text()

			}
			if text == target {
				err = opt.Click()
				if err == nil {
					b.writePrintln("select options dengan " + by + " " + value + " " + attr + " " + target + " berhasil")
					return true
				}
			}
		}
	}
	b.writePrintln("select options dengan " + by + " " + value + " " + attr + " " + target + " gagal")
	return false
}

func (b *Browser) RemoveArrayElem(arr []selenium.WebElement,index uint64) []selenium.WebElement{
	copy(arr[index:], arr[index+1:])
	arr[len(arr)-1] = nil
	return arr
}

func (b *Browser) QuitBrowser() {
	b.writePrintln("hapus session")
	err := b.Wd.Quit()
	if err != nil {
		panic(err)
	}
	b.writePrintln("hapus session berhasil\n")
}

