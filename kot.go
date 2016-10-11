package main

import (
	"log"
	"image/color"
	"math/rand"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

const (
	teksturaKotaScratcha = "scratch_cat.png"
	teksturaKotaScratchaMrug = "scratch_cat_blink.png"
	teksturaInnegoKota = "cartoon-cat-free.png"
	przyspieszenieZiemskie = 10000
	czasMrugnięcia = 0.05
	bazaOkresuMrugnięcia = 3.0
	losowyCzasDoMrugnięcia = 5.0
	szansaPodwójnegoMrugnięcia = 0.2
)

type systemObsługiKota struct {
	świat *ecs.World
	kot
}

func (sok *systemObsługiKota) New(świat *ecs.World) {
	sok.świat = świat

	sok.stwórzKota()
	
	sok.zarejestrujKlawisze()
}

func (sok *systemObsługiKota) stwórzKota() {
	sok.kot.BasicEntity = ecs.NewBasic()
	sok.kot.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 670},
		Width:    128,
		Height:   128,
	}

	tekstura, błąd := common.LoadedSprite(teksturaKotaScratcha)
	if błąd != nil {
		log.Println("Nie udało się załadować tekstury: " + błąd.Error())
	}
	sok.kot.tekstura = tekstura

	tekstura, błąd = common.LoadedSprite(teksturaKotaScratchaMrug)
	if błąd != nil {
		log.Println("Nie udało się załadować tekstury: " + błąd.Error())
	}
	sok.kot.teksturaMrug = tekstura

	sok.kot.RenderComponent = common.RenderComponent{
		Drawable: sok.kot.tekstura,
		Scale:    engo.Point{1, 1},
	}

	for _, system := range sok.świat.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&sok.kot.BasicEntity, &sok.kot.RenderComponent, &sok.kot.SpaceComponent)
		}
	}
}

func (sok *systemObsługiKota) zarejestrujKlawisze() {
	engo.Input.RegisterButton("lewo", engo.ArrowLeft)
	engo.Input.RegisterButton("prawo", engo.ArrowRight)
	engo.Input.RegisterButton("góra", engo.ArrowUp)
	engo.Input.RegisterButton("dół", engo.ArrowDown)
	engo.Input.RegisterButton("spacja", engo.Space)
}

func (sok *systemObsługiKota) Update(dt float32) {
	lewo := engo.Input.Button("lewo")
	prawo := engo.Input.Button("prawo")
	spacja := engo.Input.Button("spacja")

	staryX := sok.kot.SpaceComponent.Position.X
	staryY := sok.kot.SpaceComponent.Position.Y
	nowyX := staryX
	nowyY := staryY

	sok.kot.stoi = false
	nowyY += dt * (sok.kot.szybkośćY + (dt * przyspieszenieZiemskie / 2))
	sok.kot.szybkośćY += dt * przyspieszenieZiemskie

	if lewo.Down() && !prawo.Down() {
		nowyX -= dt * 500
	}
	if prawo.Down() && !lewo.Down() {
		nowyX += dt * 500
	}

	// Czy kot stoi na podłodze?
	if nowyY >= 670 && staryY <= 670 {
		nowyY = 670
		sok.kot.stoi = true
	}

	if spacja.Down() && sok.kot.stoi {
		sok.kot.szybkośćY = -2500
	}

	sok.kot.SpaceComponent.Position.X = nowyX
	sok.kot.SpaceComponent.Position.Y = nowyY

	sok.kot.doKońcaMrugnięcia -= dt
	sok.kot.doNastępnegoMrugnięcia -= dt
	switch {
		case sok.kot.doNastępnegoMrugnięcia <= 0 : {
			sok.kot.RenderComponent.Drawable = sok.kot.teksturaMrug
			sok.kot.doKońcaMrugnięcia = czasMrugnięcia
			sok.zaplanujMrugnięcie()
		}
		case sok.kot.doKońcaMrugnięcia <= 0 : {
			sok.kot.RenderComponent.Drawable = sok.kot.tekstura
		}
	}
}

func (sok *systemObsługiKota) zaplanujMrugnięcie() {
	if !sok.kot.podwójneMrugnięcie && rand.Float32() < szansaPodwójnegoMrugnięcia {
		sok.kot.doNastępnegoMrugnięcia = 2 * czasMrugnięcia
		sok.kot.podwójneMrugnięcie = true
	} else {
		sok.kot.doNastępnegoMrugnięcia = czasMrugnięcia + rand.Float32() * bazaOkresuMrugnięcia
		sok.kot.doNastępnegoMrugnięcia += rand.Float32() * losowyCzasDoMrugnięcia
		sok.kot.podwójneMrugnięcie = false
	}
}

func (sok *systemObsługiKota) Remove(usunięty ecs.BasicEntity) {
	if usunięty == sok.kot.BasicEntity {
		panic("Usunięto kota!")
	}
}

type kot struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	szybkośćY float32
	stoi bool
	tekstura common.Drawable
	teksturaMrug common.Drawable
	doNastępnegoMrugnięcia float32
	doKońcaMrugnięcia float32
	podwójneMrugnięcie bool
}


type kociaScena struct{}

func (*kociaScena) Type() string { return "kociaScena" }

func (*kociaScena) Preload() {
	engo.Files.Load(teksturaKotaScratcha, teksturaKotaScratchaMrug, teksturaInnegoKota)
}

func (*kociaScena) Setup(świat *ecs.World) {
	common.SetBackground(color.White)

	świat.AddSystem(&common.RenderSystem{})

	świat.AddSystem(&systemObsługiKota{})
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
