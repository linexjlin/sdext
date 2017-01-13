package sdext

import "testing"

func TestMain(t *testing.T) {
	//loadUrlDone("/tmp/")
	articles := Extracter("http://192.168.168.74:18080/Slashdot.html", "/tmp/")
	for _, v := range articles {
		t.Log(v)
	}

}

func TestLoadUrlDone(t *testing.T) {
	loadUrlDone("/tmp/")
}

func TestExtractContent(t *testing.T) {
	url := "https://science.slashdot.org/story/17/01/10/1255251/scientists-predict-star-collision-visible-to-the-naked-eye-in-2022"
	article := extractContent(url, "88469813")
	if len(article) > 1 {
		t.Log("ok!")
	}
}
