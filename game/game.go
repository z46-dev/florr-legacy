package game

import (
	"math"
	"time"

	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/util"
	"github.com/z46-dev/golog"
)

const (
	TeamID_ID_SENTINEL uint64 = iota
	TeamID_Players
	TeamID_Mobs
	TeamID_END_SENTINEL
)

func NewGame(opts GameOptions) (g *Game) {
	g = &Game{
		Time:               0,
		Options:            opts,
		clients:            shared.NewCollection[*Client](),
		flowers:            shared.NewCollection[*Flower](),
		petals:             shared.NewCollection[*Petal](),
		mobs:               shared.NewCollection[*Mob](),
		flowersSpatialHash: util.NewSpatialHash[*Flower](),
		petalsSpatialHash:  util.NewSpatialHash[*Petal](),
		mobsSpatialHash:    util.NewSpatialHash[*Mob](),
		flowersCache:       shared.NewCollection[*FlowerCache](),
		petalsCache:        shared.NewCollection[*PetalCache](),
		mobsCache:          shared.NewCollection[*MobCache](),
		Width:              50000,
		Height:             3000,
		nextEntityID:       TeamID_END_SENTINEL,
		nextClientID:       1,
		stopChan:           make(chan struct{}),
		doneChan:           make(chan struct{}),
		lastUpdate:         time.Now(),
		L:                  golog.New().Prefix("[GAME]", golog.BoldBlue).Timestamp(),
	}

	g.NumTeams = opts.NumTeams
	return
}

func (g *Game) Start() {
	go g.Run()
}

func (g *Game) Run() {
	g.runOnce.Do(func() {
		var ticker *time.Ticker = time.NewTicker(TickInterval)
		defer ticker.Stop()
		defer close(g.doneChan)

		for {
			select {
			case <-ticker.C:
				g.Update()
			case <-g.stopChan:
				return
			}
		}
	})
}

func (g *Game) Stop() {
	g.stopOnce.Do(func() {
		close(g.stopChan)
	})

	<-g.doneChan
}

func (g *Game) Update() {
	var (
		start             time.Time = time.Now()
		deltaMilliseconds float64   = float64(time.Since(g.lastUpdate).Milliseconds())
		damping           float64   = math.Pow(velocityRetention, deltaMilliseconds)
		positionFactor    float64   = velocityIntegral * (1 - damping)
	)

	g.lastUpdate = time.Now()
	g.Time++

	g.flowersSpatialHash.Clear()
	g.petalsSpatialHash.Clear()
	g.mobsSpatialHash.Clear()

	// ==================================================
	// ================= Update Flowers =================
	// ==================================================

	g.flowers.Lock.Lock()

	for _, toAdd := range g.flowers.ToAppend {
		g.flowers.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.flowers.ToRemove {
		delete(g.flowers.Items, toRemove)
	}

	g.flowers.ToAppend = g.flowers.ToAppend[:0]
	g.flowers.ToRemove = g.flowers.ToRemove[:0]

	for _, flower := range g.flowers.Items {
		flower.Update(deltaMilliseconds, damping, positionFactor)
	}

	g.flowers.Lock.Unlock()

	// ==================================================
	// ================= Update  Petals =================
	// ==================================================

	g.petals.Lock.Lock()

	for _, toAdd := range g.petals.ToAppend {
		g.petals.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.petals.ToRemove {
		delete(g.petals.Items, toRemove)
	}

	g.petals.ToAppend = g.petals.ToAppend[:0]
	g.petals.ToRemove = g.petals.ToRemove[:0]

	for _, petal := range g.petals.Items {
		petal.Update(deltaMilliseconds, damping, positionFactor)
	}

	g.petals.Lock.Unlock()

	// ==================================================
	// ================== Update  Mobs ==================
	// ==================================================

	g.mobs.Lock.Lock()

	for _, toAdd := range g.mobs.ToAppend {
		g.mobs.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.mobs.ToRemove {
		delete(g.mobs.Items, toRemove)
	}

	g.mobs.ToAppend = g.mobs.ToAppend[:0]
	g.mobs.ToRemove = g.mobs.ToRemove[:0]

	for _, mob := range g.mobs.Items {
		mob.Update(deltaMilliseconds, damping, positionFactor)
	}

	g.mobs.Lock.Unlock()

	g.Collide(deltaMilliseconds)

	// ==================================================
	// ================= Update  Caches =================
	// ==================================================

	g.flowersCache.Lock.Lock()

	for _, toAdd := range g.flowersCache.ToAppend {
		g.flowersCache.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.flowersCache.ToRemove {
		delete(g.flowersCache.Items, toRemove)
	}

	g.flowersCache.ToAppend = g.flowersCache.ToAppend[:0]
	g.flowersCache.ToRemove = g.flowersCache.ToRemove[:0]

	g.flowersCache.Lock.Unlock()

	g.petalsCache.Lock.Lock()

	for _, toAdd := range g.petalsCache.ToAppend {
		g.petalsCache.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.petalsCache.ToRemove {
		delete(g.petalsCache.Items, toRemove)
	}

	g.petalsCache.ToAppend = g.petalsCache.ToAppend[:0]
	g.petalsCache.ToRemove = g.petalsCache.ToRemove[:0]

	g.petalsCache.Lock.Unlock()

	g.mobsCache.Lock.Lock()

	for _, toAdd := range g.mobsCache.ToAppend {
		g.mobsCache.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.mobsCache.ToRemove {
		delete(g.mobsCache.Items, toRemove)
	}

	g.mobsCache.ToAppend = g.mobsCache.ToAppend[:0]
	g.mobsCache.ToRemove = g.mobsCache.ToRemove[:0]

	g.mobsCache.Lock.Unlock()

	// ==================================================
	// ================= Update Clients =================
	// ==================================================

	g.clients.Lock.Lock()

	for _, toAdd := range g.clients.ToAppend {
		g.clients.Items[toAdd.ID] = toAdd
	}

	for _, toRemove := range g.clients.ToRemove {
		delete(g.clients.Items, toRemove)
	}

	g.clients.ToAppend = g.clients.ToAppend[:0]
	g.clients.ToRemove = g.clients.ToRemove[:0]

	for _, client := range g.clients.Items {
		client.Update()
	}

	g.clients.Lock.Unlock()

	g.SpawnMobs()

	g.framesNow++

	g.mspts = append(g.mspts, float64(time.Since(start).Milliseconds()))
	if g.Time%TicksPerSecond == 0 {
		g.fps = g.framesNow
		g.framesNow = 0

		g.mspt = 0
		for _, mspt := range g.mspts {
			g.mspt += mspt
		}

		if len(g.mspts) > 0 {
			g.mspt /= float64(len(g.mspts))
		}

		g.mspts = g.mspts[:0]

		if g.fps < TicksPerSecond-1 || g.mspt > 1000.0/float64(TicksPerSecond)+1 {
			g.L.Warningf("Game performance degraded (expected, actual, threshold): FPS(%d, %d, %d), MSPT (%.2f, %.2f, %.2f)", TicksPerSecond, g.fps, TicksPerSecond-1, 1000.0/float64(TicksPerSecond), g.mspt, 1000.0/float64(TicksPerSecond)+1)
		}
	}
}
