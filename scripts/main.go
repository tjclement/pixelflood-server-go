package main

import (
	"github.com/tjclement/pixelflut-go"
	"fmt"
	"flag"
	"runtime/pprof"
	"os"
	"log"
	"github.com/tjclement/framebuffer-go"
	"os/signal"
)

func main() {
	screen_width := flag.Int("screen_width", 320, "Width of the screen to draw on")
	screen_height := flag.Int("screen_height", 400, "Height of the screen to draw on")
	display := flag.String("display", "/dev/fb0", "Name of the framebuffer device to write to")
	profile := flag.Bool("profile", false, "Set to true to enable CPU profiling > cpu.profile")
	flag.Parse()

	if *profile {
		fmt.Println("Running with CPU profiling enabled")
		prof_file, err := os.Create("cpu.profile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(prof_file)
		defer pprof.StopCPUProfile()
	}

	fb, err := framebuffer.Open(*display)

	if err != nil {
		log.Panic(err)
	}
	fb.Flush()

	fmt.Println("Starting server")
	server := pixelflood_server.NewServer(fb, uint16(*screen_width), uint16(*screen_height))
	go server.Run()
	defer server.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<- c
}
