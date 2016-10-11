package main

import (
	"log"
	"image/color"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

const (
	teksturaKotaScratcha = "scratch_cat.png"
	teksturaInnegoKota = "cartoon-cat-free.png"
	przyspieszenieZiemskie = 10000
)

type kociSystem struct {
	świat *ecs.World
	kot
}

func (ks *kociSystem) New(świat *ecs.World) {
	ks.świat = świat

	ks.stwórzKota()
	
	ks.zarejestrujKlawisze()
}

func (ks *kociSystem) stwórzKota() {
	ks.kot.BasicEntity = ecs.NewBasic()
	ks.kot.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 670},
		Width:    128,
		Height:   128,
	}

	tekstura, błąd := common.LoadedSprite(teksturaKotaScratcha)
	if błąd != nil {
		log.Println("Nie udało się załadować tekstury: " + błąd.Error())
	}

	ks.kot.RenderComponent = common.RenderComponent{
		Drawable: tekstura,
		Scale:    engo.Point{1, 1},
	}

	for _, system := range ks.świat.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ks.kot.BasicEntity, &ks.kot.RenderComponent, &ks.kot.SpaceComponent)
		}
	}
}

func (ks *kociSystem) zarejestrujKlawisze() {
	engo.Input.RegisterButton("lewo", engo.ArrowLeft)
	engo.Input.RegisterButton("prawo", engo.ArrowRight)
	engo.Input.RegisterButton("góra", engo.ArrowUp)
	engo.Input.RegisterButton("dół", engo.ArrowDown)
	engo.Input.RegisterButton("spacja", engo.Space)
}

func (ks *kociSystem) Update(dt float32) {
	lewo := engo.Input.Button("lewo")
	prawo := engo.Input.Button("prawo")
	spacja := engo.Input.Button("spacja")

	staryX := ks.kot.SpaceComponent.Position.X
	staryY := ks.kot.SpaceComponent.Position.Y
	nowyX := staryX
	nowyY := staryY

	ks.kot.stoi = false
	nowyY += dt * (ks.kot.szybkośćY + (dt * przyspieszenieZiemskie / 2))
	ks.kot.szybkośćY += dt * przyspieszenieZiemskie

	if lewo.Down() && !prawo.Down() {
		nowyX -= dt * 500
	}
	if prawo.Down() && !lewo.Down() {
		nowyX += dt * 500
	}

	if nowyY >= 670 && staryY <= 670 {
		nowyY = 670
		ks.kot.stoi = true
	}

	if spacja.Down() && ks.kot.stoi {
		ks.kot.szybkośćY = -2500
	}

	ks.kot.SpaceComponent.Position.X = nowyX
	ks.kot.SpaceComponent.Position.Y = nowyY
}

func (ks *kociSystem) Remove(usunięty ecs.BasicEntity) {
	if usunięty == ks.kot.BasicEntity {
		panic("Usunięto kota!")
	}
}


type kot struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	szybkośćY float32
	stoi bool
}


type kociaScena struct{}

func (*kociaScena) Type() string { return "kociaScena" }

func (*kociaScena) Preload() {
	engo.Files.Load(teksturaKotaScratcha, teksturaInnegoKota)
}

func (*kociaScena) Setup(świat *ecs.World) {
	common.SetBackground(color.White)

	świat.AddSystem(&common.RenderSystem{})

	świat.AddSystem(&kociSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:  "Kot",
		Width:  1200,
		Height: 800,
		VSync: true,
	}

	engo.Run(opts, &kociaScena{})
}
