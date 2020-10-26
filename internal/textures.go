package internal

import (
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"os"
	"path"
)

var textures = []string {
	"eagle.png",
	"redbrick.png",
	"purplestone.png",
	"greystone.png",
	"bluestone.png",
	"mossy.png",
	"wood.png",
	"colorstone.png",
	"barrel.png",
	"greenlight.png",
	"pillar.png",
}



func LoadTextures() []*pixel.PictureData {
	loadedTextures := make([]*pixel.PictureData, 0)
	for _, file := range textures {
		texture, err := loadTexture(path.Join("assets", file))
		if err == nil {
			loadedTextures = append(loadedTextures, texture)
		}
	}
	return loadedTextures
}

func loadTexture(filename string) (*pixel.PictureData, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}