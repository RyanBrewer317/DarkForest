package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	// "os"
	// "time"

	// "github.com/faiface/beep"
	// "github.com/faiface/beep/effects"
	// "github.com/faiface/beep/mp3"
	// "github.com/faiface/beep/speaker"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// var (
// 	audio_paused bool
// 	audio_queue  []beep.StreamSeeker
// 	music_silent bool
// 	music_volume float64
// 	sfx_mixer    beep.Mixer
// 	audio_mixer  beep.Mixer
// 	music_mixer  beep.Mixer
// )

// type music struct {
// 	streamer    *beep.StreamSeekCloser
// 	resampler   *beep.Resampler
// 	sample_rate beep.SampleRate
// }

// type sfx struct {
// 	buffer      *beep.Buffer
// 	sample_rate beep.SampleRate
// }

// type music_streamer struct{}

// func (s *music_streamer) Stream(samples [][2]float64) (int, bool) {
// 	if audio_paused {
// 		for i := range samples {
// 			samples[i] = [2]float64{}
// 		}
// 		return len(samples), true
// 	}
// 	filled := 0
// 	for filled < len(samples) {
// 		if len(audio_queue) == 0 {
// 			for i := range samples[filled:] {
// 				samples[i][0] = 0
// 				samples[i][1] = 0
// 			}
// 			break
// 		}

// 		n, ok := audio_queue[0].Stream(samples[filled:])
// 		gain := 0.0
// 		if !music_silent {
// 			gain = math.Pow(2, music_volume)
// 		}
// 		for i := range samples[:n] {
// 			samples[i][0] *= gain
// 			samples[i][1] *= gain
// 		}
// 		if !ok {
// 			audio_queue = audio_queue[1:]
// 			if len(audio_queue) > 0 {
// 				audio_queue[0].Seek(0)
// 			}
// 		}
// 		filled += n
// 	}
// 	return len(samples), true
// }

// func (q *music_streamer) Err() error {
// 	return nil
// }

const AXETIME = 40

type Timer int

func (t *Timer) update(callback func()) {
	*t -= Timer(1)
	if *t == 0 {
		callback()
	}
}
func (t *Timer) resetTo(tm int) { *t = Timer(tm) }

type Direction int

const (
	DIR_RIGHT Direction = iota
	DIR_LEFT
)

type RigidBody struct {
	x          float64
	y          float64
	xVel       float64
	yVel       float64
	xAcc       float64
	yAcc       float64
	generation int
}

type TreeType int

const (
	TREE_FULL TreeType = iota
	TREE_TOP
	TREE_EMPTY
)

type TreeCell struct {
	x          int
	y          int
	t          TreeType
	lightLevel int
}

func getTree(x float64, y float64) *TreeCell {
	i := int(y / 32)
	j := int(x / 32)
	return &treecells[i*8+j]
}

var treecells [64]TreeCell

type Monster struct {
	physics *RigidBody
	dir     float64 // radians
}

var monsters [16]Monster

type GameState int

const (
	PLAYING GameState = iota
	PAUSED
	START
	VICTORY
	LOSS
)

type Game struct {
	state             GameState
	playerDir         Direction
	axeAnimationFrame int
	axeAnimationTimer Timer
	playerWalking     bool
	playerRigidBody   *RigidBody
	progressBarFrame  int
	progressBarTimer  Timer
	playerAxing       bool
	monsterFrame      int
	monsterTimer      Timer
	togglingPause     bool
	restarting        bool
	health            int
}

func gameYtoPx(y float64, h int, playerY float64) float64 {
	return (140.0 - float64(h)) - ((y - playerY) * 10)
}

func gameXtoPx(x float64, playerX float64) float64 {
	return (x-playerX)*10 + 144
}

var gamefont font.Face

var playerImg *ebiten.Image
var axeImg *ebiten.Image
var treeMiddleImg *ebiten.Image
var fullLeavesImg *ebiten.Image
var treeTopImg *ebiten.Image
var leavesTopImg *ebiten.Image
var skyTopImg *ebiten.Image
var skyImg *ebiten.Image
var progressBarImg *ebiten.Image
var monsterImg *ebiten.Image
// var helloSound sfx
// var helloByeBye1Sound sfx
// var helloByeBye2Sound sfx
// var wee1Sound sfx
// var wee2Sound sfx
// var choppingSound sfx
// var treeFallSound sfx
// var song1 music

func init() {
	var err error
	parsedFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	gamefont = truetype.NewFace(parsedFont, &truetype.Options{})
	playerImg, _, err = ebitenutil.NewImageFromFile("DarkForestPlayer1.png")
	if err != nil {
		log.Fatal(err)
	}
	axeImg, _, err = ebitenutil.NewImageFromFile("AxeCarry.png")
	if err != nil {
		log.Fatal(err)
	}
	treeMiddleImg, _, err = ebitenutil.NewImageFromFile("treemiddle.png")
	if err != nil {
		log.Fatal(err)
	}
	fullLeavesImg, _, err = ebitenutil.NewImageFromFile("fullleaves.png")
	if err != nil {
		log.Fatal(err)
	}
	treeTopImg, _, err = ebitenutil.NewImageFromFile("treetop.png")
	if err != nil {
		log.Fatal(err)
	}
	leavesTopImg, _, err = ebitenutil.NewImageFromFile("leavestop.png")
	if err != nil {
		log.Fatal(err)
	}
	skyTopImg, _, err = ebitenutil.NewImageFromFile("skytop.png")
	if err != nil {
		log.Fatal(err)
	}
	skyImg, _, err = ebitenutil.NewImageFromFile("sky.png")
	if err != nil {
		log.Fatal(err)
	}
	progressBarImg, _, err = ebitenutil.NewImageFromFile("progressBar.png")
	if err != nil {
		log.Fatal(err)
	}
	monsterImg, _, err = ebitenutil.NewImageFromFile("monster.png")
	if err != nil {
		log.Fatal(err)
	}
	// audio_paused = false
	// music_silent = false
	// music_volume = 1
	// fileNames := []string{"hello.mp3", "hellobyebye1.mp3", "hellobyebye2.mp3", "wee1.mp3", "wee2.mp3", "chopping.mp3", "treefall.mp3"}
	// files := []*sfx{&helloSound, &helloByeBye1Sound, &helloByeBye2Sound, &wee1Sound, &wee2Sound, &choppingSound, &treeFallSound}
	// for i := 0; i < len(fileNames); i++ {
	// 	f, err := os.Open(fileNames[i])
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	streamer, format, err := mp3.Decode(f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	buffer := beep.NewBuffer(format)
	// 	buffer.Append(streamer)
	// 	streamer.Close()
	// 	*files[i] = sfx{
	// 		buffer:      buffer,
	// 		sample_rate: format.SampleRate,
	// 	}
	// }
	// f, err := os.Open("music.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// streamer, format, err := mp3.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// loop := beep.Loop(-1, streamer)
	// resampled := beep.Resample(4, format.SampleRate, 44100, loop)
	// vol := &effects.Volume{
	// 	Streamer: resampled,
	// 	Base: 2,
	// 	Volume: music_volume - 5,
	// 	Silent: music_silent,
	// }
	// speaker.Init(beep.SampleRate(44100), beep.SampleRate(44100).N(time.Second/10))
	// go speaker.Play(&audio_mixer) // this is the only streamer that will ever actually be playing. It's silent at first
	// audio_mixer.Add(&music_mixer)
	// audio_mixer.Add(&sfx_mixer)
	// music_mixer.Add(vol)
}

// var playedSquirrel bool
// var playingSquirrel bool

// func play(sound *sfx) {
// 	buffer := sound.buffer
// 	s := buffer.Streamer(0, buffer.Len())
// 	resampled := beep.Resample(4, sound.sample_rate, 44100, s)
// 	volume := &effects.Volume{
// 		Base:     2,
// 		Volume:   music_volume,
// 		Streamer: resampled,
// 		Silent:   music_silent,
// 	}
// 	if sound == &helloByeBye1Sound || sound == &helloByeBye2Sound || sound == &helloSound {
// 		playedSquirrel = true
// 		playedSquirrel = true
// 		volume.Volume = -3
// 		seq := beep.Seq(volume, beep.Callback(func() {
// 			playingSquirrel = false
// 		}))
// 		sfx_mixer.Add(seq)
// 		return
// 	} else if sound == &wee1Sound || sound == &wee2Sound {
// 		volume.Volume = -3
// 	}
// 	sfx_mixer.Add(volume)
// }

func onGround(rb RigidBody) bool {
	tree := getTree(rb.x, rb.y)
	switch tree.t {
	case TREE_FULL:
		return rb.y == float64(int(rb.y)) && int(rb.y)%8 == 0
	case TREE_TOP:
		if rb.y < 32 {
			return rb.y == 0
		}
		if getTree(rb.x, rb.y-32).t == TREE_FULL {
			return rb.y == float64(tree.y)
		}
	case TREE_EMPTY:
		if rb.y < 32 {
			return rb.y == 0
		}
		if getTree(rb.x, rb.y-32).t == TREE_FULL {
			return rb.y == float64(tree.y)
		}
	}
	return false
}

func chopSquare(x float64, y float64) *TreeCell {
	tree := getTree(x, y)
	if math.Abs(float64(tree.x+16)-x) < 4 {
		if math.Abs(float64(tree.y)-y) < 1 {
			return tree
		}
	}
	return nil
}

func chop(tree *TreeCell) {
	// if rand.Intn(2) == 0 {
	// 	sound := ([]*sfx{&wee1Sound, &wee2Sound})[rand.Intn(2)]
	// 	play(sound)
	// }
	tree.t = TREE_EMPTY
	if tree.y < 7*32 {
		chop(getTree(float64(tree.x), float64(tree.y+32)))
	}
}

// no idea how to explain the general utility of this function...
// it maps x < 0 to 1.0 and x >= 0 to -1.0
func negCoef(x float64) float64 {
	if x < 0 {
		return 1
	}
	return -1
}

func (g *Game) Update() error {
	switch g.state {
	case VICTORY, LOSS:
		if ebiten.IsKeyPressed(ebiten.KeyF) {
			ebiten.SetFullscreen(true)
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			ebiten.SetFullscreen(false)
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.restarting = true
		}
		if g.restarting && !ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.restarting = false
			for i := 0; i < 8; i++ {
				for j := 0; j < 8; j++ {
					t := TREE_FULL
					if i == 7 {
						t = TREE_TOP
					}
					treecells[i*8+j] = TreeCell{
						x:          j * 32,
						y:          i * 32,
						t:          t,
						lightLevel: i,
					}
				}
			}
			g.playerRigidBody = &RigidBody{x: 128, y: 255}
			for i := 0; i < 16; i++ {
				m := Monster{physics: &RigidBody{x: float64(i * 16), y: 64}}
				monsters[i] = m
			}
			g.axeAnimationTimer = Timer(20)
			g.monsterTimer = Timer(10)
			g.health = 128
			g.state = START
		}
	case START:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.state = PLAYING
		}
		if ebiten.IsKeyPressed(ebiten.KeyF) {
			ebiten.SetFullscreen(true)
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			ebiten.SetFullscreen(false)
		}
	case PAUSED:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.togglingPause = true
			return nil
		}
		if g.togglingPause && !ebiten.IsKeyPressed(ebiten.KeyP) {
			g.togglingPause = false
			g.state = PLAYING
			return nil
		}
	case PLAYING:
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			g.togglingPause = true
			return nil
		}
		if g.togglingPause && !ebiten.IsKeyPressed(ebiten.KeyP) {
			g.state = PAUSED
			g.togglingPause = false
			return nil
		}
		rb := g.playerRigidBody
		tree := getTree(rb.x, rb.y)
		if rb.x < 0 || rb.x > 255 || rb.y < 0 || rb.y > 255 {
			if rb.x < 0 {
				rb.x = 0
			} else if rb.x > 255 {
				rb.x = 255
			}
			if rb.y < 0 {
				rb.y = 0
			} else if rb.y > 255 {
				rb.y = 255
			}
			rb.xVel += rb.xAcc
			rb.yVel += rb.yAcc
			rb.x += rb.xVel
			if rb.x < 0 {
				rb.x = 0
			} else if rb.x > 255 {
				rb.x = 255
			}
			if rb.y < 0 {
				rb.y = 0
			} else if rb.y > 255 {
				rb.y = 255
			}
		} else {
			rb.xVel += rb.xAcc
			rb.yVel += rb.yAcc
			if rb.x < 255 && rb.xVel > 0 || rb.x > 0 && rb.xVel < 0 {
				rb.x += rb.xVel
			}
			if rb.y < 255 && rb.yVel > 0 || !onGround(*rb) && rb.yVel < 0 {
				rb.y += rb.yVel
			}
		}
		if tree.t == TREE_FULL && (math.Mod(rb.y, 8) != 0 && math.Abs(math.Mod(rb.y, 8)) < 0.6) {
			rb.y = math.Round(rb.y/8) * 8
		}
		if rb.y < 0 {
			rb.y = 0
		}
		if !onGround(*g.playerRigidBody) {
			g.playerRigidBody.yVel -= 0.2
		} else {
			g.playerRigidBody.yVel = 0
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			g.playerDir = DIR_LEFT
			if g.playerRigidBody.x > 0 {
				g.playerRigidBody.x -= 0.5
			}
			g.playerWalking = true
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.playerDir = DIR_RIGHT
			if g.playerRigidBody.x < 255 {
				g.playerRigidBody.x += 0.5
			}
			g.playerWalking = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			if g.playerRigidBody.y > 0 {
				g.playerRigidBody.y -= 0.1
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) && onGround(*g.playerRigidBody) {
			g.playerRigidBody.yVel += 1.7
		}
		if !ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyD) {
			g.playerWalking = false
		}
		if ebiten.IsKeyPressed(ebiten.KeyJ) && !g.playerAxing {
			t := chopSquare(g.playerRigidBody.x, g.playerRigidBody.y)
			if t != nil && t.t != TREE_EMPTY {
				g.progressBarTimer = Timer(AXETIME)
				g.progressBarFrame = 0
				g.playerAxing = true
				// play(&choppingSound)
			}
		}
		if !ebiten.IsKeyPressed(ebiten.KeyJ) {
			g.playerAxing = false
		}
		if ebiten.IsKeyPressed(ebiten.KeyF) {
			ebiten.SetFullscreen(true)
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			ebiten.SetFullscreen(false)
		}
		if g.playerWalking {
			g.axeAnimationTimer.update(func() {
				g.axeAnimationFrame = (g.axeAnimationFrame + 1) % 4
				g.axeAnimationTimer.resetTo(10)
			})
		}
		if g.playerAxing {
			g.progressBarTimer.update(func() {
				g.progressBarFrame = g.progressBarFrame + 1
				if g.progressBarFrame == 4 {
					g.playerAxing = false
					t := chopSquare(g.playerRigidBody.x, g.playerRigidBody.y)
					if t != nil {
						// play(&treeFallSound)
						chop(t)
					}
				} else {
					g.progressBarTimer.resetTo(AXETIME)
				}
			})
		}
		g.monsterTimer.update(func() {
			g.monsterFrame = (g.monsterFrame + 1) % 2
			g.monsterTimer.resetTo(10)
		})
		for i := 6; i > -1; i-- {
			for j := 0; j < 8; j++ {
				treecells[i*8+j].lightLevel = treecells[(i+1)*8+j].lightLevel - 1
				if treecells[(i+1)*8+j].t == TREE_EMPTY {
					treecells[i*8+j].lightLevel++
				}
			}
		}
		es := 0
		for i := 0; i < 64; i++ {
			if treecells[i].t == TREE_EMPTY {
				es++
			}
		}
		if es == 64 {
			g.state = VICTORY
			return nil
		}
		ceil := float64((4 - gameStage(es)) * 32)
		for i := 0; i < 16; i++ {
			m := &monsters[i]
			p := m.physics
			a := p.x - rb.x
			b := p.y - rb.y
			c := math.Sqrt(a*a + b*b)
			var d float64
			if c < 10 && c > 0.2 {
				d = math.Acos(a/c) + float64(rand.Intn(20)-10)*math.Pi/180
			} else {
				d = m.dir + float64(rand.Intn(60)-30)*math.Pi/180
			}
			if c < 0.2 {
				g.health--
				if g.health == 0 {
					g.state = LOSS
					return nil
				}
			}
			m.dir = math.Mod(d, 2*math.Pi)
			p.x -= math.Cos(d) / 3
			p.y += negCoef(b) * math.Sin(d) / 3
			if p.x < 0 {
				p.x = 0
				m.dir = math.Pi
			} else if p.x > 255 {
				p.x = 255
				m.dir = 0
			}
			if p.y < 0 {
				p.y = 0
				m.dir = 3 * math.Pi
			} else if p.y > ceil {
				p.y = ceil
				m.dir = math.Pi / 2
			}
		}
		// dx := math.Abs(math.Mod(g.playerRigidBody.x, 32) - 16)
		// dy := math.Abs(math.Mod(g.playerRigidBody.y, 32) - 20)
		// if !playingSquirrel && !playedSquirrel && dx < 5 && dy < 5 {
		// 	sound := ([]*sfx{&helloSound, &helloByeBye1Sound, &helloByeBye2Sound})[rand.Intn(3)]
		// 	play(sound)
		// }
		// if dx >= 5 || dy >= 5 {
		// 	playedSquirrel = false
		// }
	}
	return nil
}

func gameStage(empties int) int {
	if empties == 64 {
		return 4
	}
	return empties >> 4 // integer division empties/16 to get [0,3]
}

func healthStage(health int) int {
	if health == 128 {
		return 3
	}
	return health >> 5 // integer division health/32 to get [0,3]
}

func write(i *ebiten.Image, s string, line int) {
	var c color.Color
	if line == 0 {
		c = color.White
	} else {
		c = color.Gray{128}
	}
	text.Draw(i, s, gamefont, 20, 15*line+20, c)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case VICTORY:
		write(screen, "Victory!", 0)
		write(screen, "Credits", 2)
		write(screen, "Visual Artist: Wawa", 3)
		write(screen, "Voice Actor: Wawa", 4)
		write(screen, "Developer: Ryan Brewer", 5)
		write(screen, "A Raining Studios Production", 7)
		write(screen, "Press [space] to play again!", 13)
	case LOSS:
		write(screen, "Defeat!", 0)
		write(screen, "Credits", 2)
		write(screen, "Visual Artist: Wawa", 3)
		write(screen, "Voice Actor: Wawa", 4)
		write(screen, "Developer: Ryan Brewer", 5)
		write(screen, "A Raining Studios Production", 7)
		write(screen, "Press [space] to play again!", 13)
	case START:
		write(screen, "The Dark Forest, by Raining Studios", 0)
		write(screen, "You are a hopeful lumberjack, hoping to take on the", 1)
		write(screen, "Dark Forest. They called you a fool. But you know", 2)
		write(screen, "there's money to be made! Chop trees at the bottom", 3)
		write(screen, "of each 'light section' of the tree. Clear the forest", 4)
		write(screen, "to win! Though, they say 'beware the monsters'...", 5)
		write(screen, "[space] to begin play", 6)
		write(screen, "[w] to jump", 7)
		write(screen, "[a] to go left", 8)
		write(screen, "[s] to drop down", 9)
		write(screen, "[d] to go right", 10)
		write(screen, "[p] to pause/unpause", 11)
		write(screen, "[j] to chop (only at certain spots :)", 12)
		write(screen, "[f] to enter fullscreen mode", 13)
		write(screen, "[esc] to exit fullscreen mode", 14)
	case PAUSED:
		write(screen, "Game Paused", 0)
		write(screen, "Press [p] to play!", 1)
	case PLAYING:
		playerImgOpt := ebiten.DrawImageOptions{}
		progressBarImgOpt := ebiten.DrawImageOptions{}
		if g.playerDir == DIR_LEFT {
			playerImgOpt.GeoM.Scale(-1, 1)
			playerImgOpt.GeoM.Translate(32, 0)
		}
		xpx := gameXtoPx(g.playerRigidBody.x, g.playerRigidBody.x)
		ypx := gameYtoPx(g.playerRigidBody.y, 32, g.playerRigidBody.y)
		playerImgOpt.GeoM.Translate(xpx, ypx)
		progressBarImgOpt.GeoM.Translate(xpx, ypx)
		for i := 0; i < 8; i++ {
			skyTopOpt := ebiten.DrawImageOptions{}
			skyTopOpt.GeoM.Translate(gameXtoPx(float64(32*i), g.playerRigidBody.x), gameYtoPx(32*7, 320, g.playerRigidBody.y))
			screen.DrawImage(skyTopImg, &skyTopOpt)
		}
		for i := 0; i < 7; i++ {
			for j := 0; j < 8; j++ {
				skyOpt := colorm.DrawImageOptions{}
				skyOpt.GeoM.Translate(gameXtoPx(float64(32*j), g.playerRigidBody.x), gameYtoPx(float64(32*i), 320, g.playerRigidBody.y))
				skyColorOpt := colorm.ColorM{}
				skyColorOpt.ChangeHSV(0, 1, float64(treecells[i*8+j].lightLevel)/8)
				colorm.DrawImage(screen, skyImg, skyColorOpt, &skyOpt)
			}
		}
		for i := 0; i < 64; i++ {
			if treecells[i].t == TREE_EMPTY {
				continue
			}
			treeDrawOpt := colorm.DrawImageOptions{}
			treeDrawOpt.GeoM.Translate(gameXtoPx(float64(treecells[i].x), g.playerRigidBody.x), gameYtoPx(float64(treecells[i].y), 320, g.playerRigidBody.y))
			treeColorOpt := colorm.ColorM{}
			treeColorOpt.ChangeHSV(0, 1, float64(treecells[i].lightLevel)/8)
			switch treecells[i].t {
			case TREE_FULL:
				colorm.DrawImage(screen, fullLeavesImg, treeColorOpt, &treeDrawOpt)
				colorm.DrawImage(screen, treeMiddleImg, treeColorOpt, &treeDrawOpt)
			case TREE_TOP:
				colorm.DrawImage(screen, leavesTopImg, treeColorOpt, &treeDrawOpt)
				colorm.DrawImage(screen, treeTopImg, treeColorOpt, &treeDrawOpt)
			}
		}
		screen.DrawImage(playerImg, &playerImgOpt)
		framex := g.axeAnimationFrame * 32
		frameRect := image.Rect(framex, 0, framex+32, 32)
		frame := ebiten.NewImageFromImage(axeImg.SubImage(frameRect))
		screen.DrawImage(frame, &playerImgOpt)
		if g.playerAxing {
			screen.DrawImage(ebiten.NewImageFromImage(progressBarImg.SubImage(image.Rectangle{image.Point{g.progressBarFrame * 32, 0}, image.Point{g.progressBarFrame*32 + 32, 32}})), &progressBarImgOpt)
		}
		for i := 0; i < 16; i++ {
			m := monsters[i]
			monsterColorOpt := colorm.ColorM{}
			tree := getTree(m.physics.x, m.physics.y)
			monsterColorOpt.ChangeHSV(0, 1, float64(tree.lightLevel)/8)
			monsterOpt := colorm.DrawImageOptions{}
			monsterOpt.GeoM.Translate(gameXtoPx(m.physics.x, g.playerRigidBody.x), gameYtoPx(m.physics.y, 32, g.playerRigidBody.y))
			colorm.DrawImage(screen, ebiten.NewImageFromImage(monsterImg.SubImage(image.Rect(g.monsterFrame*32, 0, g.monsterFrame*32+32, 32))), monsterColorOpt, &monsterOpt)
		}
		healthOpt := &ebiten.DrawImageOptions{}
		healthOpt.GeoM.Scale(2, 1.5)
		healthOpt.GeoM.Translate(20, 20)
		healthFrame := 3 - healthStage(g.health)
		screen.DrawImage(ebiten.NewImageFromImage(progressBarImg.SubImage(image.Rect(healthFrame*32, 0, healthFrame*32+32, 32))), healthOpt)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			t := TREE_FULL
			if i == 7 {
				t = TREE_TOP
			}
			treecells[i*8+j] = TreeCell{
				x:          j * 32,
				y:          i * 32,
				t:          t,
				lightLevel: i,
			}
		}
	}
	playerRB := RigidBody{x: 128, y: 255}
	for i := 0; i < 16; i++ {
		m := Monster{physics: &RigidBody{x: float64(i * 16), y: 64}}
		monsters[i] = m
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("The Dark Forest")
	if err := ebiten.RunGame(&Game{axeAnimationTimer: Timer(10), state: START, playerRigidBody: &playerRB, monsterTimer: Timer(10), health: 128}); err != nil {
		log.Fatal(err)
	}
}
