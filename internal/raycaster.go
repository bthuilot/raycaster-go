package internal

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	color2 "image/color"
	"math"
	"time"
)

func RaycasterLoop(world *World, buffer *pixel.PictureData) {
	for x := 0; x < ScreenWidth; x += 1 {
		cameraX := float64(2 * x)/ float64(ScreenWidth) - 1
		rayDirX := world.playerDir.x + world.cameraPlane.x * cameraX
		rayDirY := world.playerDir.y + world.cameraPlane.y * cameraX

		// Find coordinate of map currently in
		mapX := int(world.playerPos.x)
		mapY := int(world.playerPos.y)


		deltaDistX := math.Abs(1 / rayDirX)
		deltaDistY := math.Abs(1 / rayDirY)

		hit := 0

		stepX, sideDistX := calculateStepAndSideDist(rayDirX, world.playerPos.x, float64(mapX), deltaDistX)
		stepY, sideDistY := calculateStepAndSideDist(rayDirY, world.playerPos.y, float64(mapY), deltaDistY)

		var side int
		// DDA
		for hit == 0  {
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = 0
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = 1
			}
			if world.worldMap[mapX][mapY] > 0 {
				hit = 1
			}
		}

		var prepWallDist float64
		if side == 0 {
			prepWallDist = calculatePrepWallDist(float64(mapX), world.playerPos.x, float64(stepX), rayDirX)
		} else {
			prepWallDist = calculatePrepWallDist(float64(mapY), world.playerPos.y, float64(stepY), rayDirY)
		}

		lineHeight := int(ScreenHeight / prepWallDist)

		drawStart := -lineHeight / 2 + ScreenHeight/ 2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight / 2 + ScreenHeight/ 2
		if drawEnd >= ScreenHeight {
			drawEnd = ScreenHeight - 1
		}

		texNum := world.worldMap[mapX][mapY] - 1
		texture := world.textures[texNum]
		bounds := texture.Bounds()
		texH := bounds.H()
		texW := bounds.W()


		var wallX float64
		if side == 0 {
			wallX = world.playerPos.y + prepWallDist * rayDirY
		} else {
			wallX = world.playerPos.x + prepWallDist * rayDirX
		}

		wallX -= math.Floor(wallX)

		texX := int(wallX * texW)
		if (side == 0 && rayDirX > 0) || (side == 1 && rayDirY < 0){
			texX = TextureWidth - texX - 1
		}

		step := texH / float64(lineHeight)
		texPos := float64(drawStart - ScreenHeight/2 + lineHeight/ 2) * step

		for y := drawStart; y < drawEnd; y++ {
			texY := int(texPos) & (int(texH) - 1)
			texPos += step
			texIndex := world.textures[texNum].Index(pixel.V(float64(texX), float64(texY)))
			color := world.textures[texNum].Pix[texIndex]
			if side == 1 {
				color = color2.RGBA{
					R: color.R/2,
					B: color.B/2,
					G: color.G/2,
					A: color.A/2,
				}
			}
			index := buffer.Index(pixel.V(float64(x), float64(y)))
			buffer.Pix[index] = color
		}
	}
}

func calculateStepAndSideDist(rayDir float64, playerPos float64, mapPos float64, delta float64) (step int, sideDist float64) {
	if rayDir < 0 {
		step = -1
		sideDist = (playerPos - mapPos) * delta
	} else {
		step = 1
		sideDist = (mapPos + 1.0 - playerPos) * delta
	}
	return
}

func calculatePrepWallDist(mapPos float64, pos float64, step float64, rayDir float64) float64 {
	return (mapPos - pos + (1 - step) / 2) / rayDir
}

func CalculateMovement(world *World, win *pixelgl.Window) {
	frameTime := float64(time.Since(world.prevTime).Nanoseconds()) / 1000000000.0
	world.prevTime = time.Now()
	world.Fps = 1.0 / frameTime
	// Add FPS to top of screen


	moveSpeed := frameTime * 5
	rotSpeed := frameTime * 3

	if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW){
		if  world.worldMap[int(world.playerPos.x + world.playerDir.x * moveSpeed)][int(world.playerPos.y)] == 0 {
			world.playerPos.x += world.playerDir.x * moveSpeed
		}
		if world.worldMap[int(world.playerPos.x)][int(world.playerPos.y + world.playerDir.y * moveSpeed)] == 0 {
			world.playerPos.y += world.playerDir.y * moveSpeed
		}
	}

	if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS){
		if  world.worldMap[int(world.playerPos.x - world.playerDir.x * moveSpeed)][int(world.playerPos.y)] == 0 {
			world.playerPos.x -= world.playerDir.x * moveSpeed
		}
		if world.worldMap[int(world.playerPos.x)][int(world.playerPos.y - world.playerDir.y * moveSpeed)] == 0 {
			world.playerPos.y -= world.playerDir.y * moveSpeed
		}
	}

	if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD){
		prevDirX := world.playerDir.x
		world.playerDir.x = world.playerDir.x * math.Cos(-rotSpeed) - world.playerDir.y * math.Sin(-rotSpeed)
		world.playerDir.y = prevDirX* math.Sin(-rotSpeed) + world.playerDir.y * math.Cos(-rotSpeed)
		prevCameraPlaneX := world.cameraPlane.x
		world.cameraPlane.x = world.cameraPlane.x * math.Cos(-rotSpeed) - world.cameraPlane.y * math.Sin(-rotSpeed)
		world.cameraPlane.y = prevCameraPlaneX * math.Sin(-rotSpeed) + world.cameraPlane.y * math.Cos(-rotSpeed)
	}
	if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA){
		prevDirX := world.playerDir.x
		world.playerDir.x = world.playerDir.x * math.Cos(rotSpeed) - world.playerDir.y * math.Sin(rotSpeed)
		world.playerDir.y = prevDirX* math.Sin(rotSpeed) + world.playerDir.y * math.Cos(rotSpeed)
		prevCameraPlaneX := world.cameraPlane.x
		world.cameraPlane.x = world.cameraPlane.x * math.Cos(rotSpeed) - world.cameraPlane.y * math.Sin(rotSpeed)
		world.cameraPlane.y = prevCameraPlaneX * math.Sin(rotSpeed) + world.cameraPlane.y * math.Cos(rotSpeed)
	}
}