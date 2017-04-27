package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	dir := os.Args[1]
	files, err := filepath.Glob(dir)
	if err != nil {
		log.Fatal(err)
	}
	var inputs string
	var filters string
	concats := "[0:v]"
	for i, e := range files[:len(files)-1] {
		inputs = fmt.Sprintf("%s -loop 1 -t 1 -i %s", inputs, e)
		filters = fmt.Sprintf(
			"%s[%d:v][%d:v]blend=all_expr='A*(if(gte(T,0.5),1,T/0.5))+B*(1-(if(gte(T,0.5),1,T/0.5)))'[b%dv]; ",
			filters, i+1, i, i+1)
		concats = fmt.Sprintf("%s[b%dv][%d:v]", concats, i+1, i+1)
	}
	// TODO templating ?
	imageNumber := strconv.Itoa((len(files)-1)*2 + 1)
	inputs = fmt.Sprintf("%s -loop 1 -t 1 -i %s", inputs, files[len(files)-1 : len(files)][0])
	cmdString := "ffmpeg" + inputs + " -filter_complex \"" + filters + concats + "concat=n=" + string(imageNumber) + ":v=1:a=0,format=yuv420p[v]\"" + " -map \"[v]\" out.mp4"
	fmt.Println(cmdString)
}
