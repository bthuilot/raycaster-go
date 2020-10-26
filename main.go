package tinyraycaster_go

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"tinyraycaster-go/internal"
)

func LaunchWindow() {
	pixelgl.Run(run)
}


func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Simple Raycaster",
		Bounds: pixel.R(0, 0, internal.ScreenWidth, internal.ScreenHeight),
	}


	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	world := internal.CreateWorldMap()
	for !win.Closed() {
		buffer := pixel.PictureData{
			Pix: make([]color.RGBA, internal.ScreenHeight * internal.ScreenWidth),
			Stride: internal.ScreenWidth,
			Rect: pixel.R(0,0,internal.ScreenWidth, internal.ScreenHeight),
		}

		// Calculate ray cast
		internal.RaycasterLoop(&world, &buffer)
		// Calculate FPS and movement
		internal.CalculateMovement(&world, win)

		// Clear window and draw scene
		win.Clear(colornames.Black)
		//imd.Draw()
		picture := pixel.NewSprite(&buffer, buffer.Bounds())
		picture.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		// Draw FPS to the screen
		fpsText := text.New(pixel.V(10, internal.ScreenHeight * .98), text.NewAtlas(basicfont.Face7x13, text.ASCII))
		fmt.Fprint(fpsText, "FPS: ", world.Fps)
		fpsText.Draw(win, pixel.IM)
		win.Update()
	}
}