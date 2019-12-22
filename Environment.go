package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
)

const (
	GOLD_SCORE = 1000
	PENALTY_SCORE = -1000
	ACTION_SCORE = -1
	WIY_ID = 1000
	PANNA_ID = 100
	GOLD_ID = 10
	HOMA_ID = 1
	ARROW_LAUNCH_SCORE = -9
	SPAWN_CHANCE = 0.2
	FIELD_SIZE = 4
)

type Arrow struct {
	pos Pos
	dir LookDirection
	isLaunched bool
}

type Environment struct {
	arrow Arrow
	homa *Homa
	wiyPos Pos
	pannaCount uint
	pannaPoses []Pos
	goldPoses []Pos
}

func NewEnvironment() *Environment {
	var env Environment
	env.homa = NewHoma()

	return &env
}

func GetMovement(lookDir LookDirection) Pos {
	switch lookDir {
	case UP:
		return Pos{
			x: 0,
			y: -1,
		}
	case DOWN:
		return Pos{
			x: 0,
			y: 1,
		}
	case LEFT:
		return Pos{
			x: -1,
			y: 0,
		}
	case RIGHT:
		return Pos{
			x: 1,
			y: 0,
		}
	}
	return Pos{}
}

func GetRotationClockwise(lookDir LookDirection) LookDirection {
	switch lookDir {
	case UP:
		return RIGHT
	case DOWN:
		return  LEFT
	case LEFT:
		return UP
	case RIGHT:
		return DOWN
	}
	return UP
}

func IsPosFree(poses []Pos, posToCheck Pos) bool {
	for _, p := range poses {
		if p == posToCheck {
			return false
		}
	}
	return true
}

func InBounds(homaPos Pos) bool {
	return homaPos.x >= 0 && homaPos.x < FIELD_SIZE && homaPos.y >= 0 && homaPos.y < FIELD_SIZE
}

func (e *Environment) InitField(goldCount uint, pannaCount uint) {
	e.pannaCount = pannaCount
	e.wiyPos = e.GetRandomFreePosExclude(e.goldPoses)
	for i := uint(0); i < goldCount; i++ {
		e.goldPoses = append(e.goldPoses, e.GetRandomFreePosExclude(e.goldPoses))
	}
}

func (e *Environment) GetScore() int {
	return e.homa.GetScore()
}

func (e *Environment) PrintField() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()

	var field [FIELD_SIZE][FIELD_SIZE]int
	if e.wiyPos.y != -1000 {
		field[e.wiyPos.y][e.wiyPos.x] = WIY_ID
	}
	field[e.homa.homaPos.y][e.homa.homaPos.x] += HOMA_ID

	for _, gold := range e.goldPoses {
		field[gold.y][gold.x] += GOLD_ID
	}

	for _, panna := range e.pannaPoses {
		field[panna.y][panna.x] += PANNA_ID
	}

	for i := 0; i < FIELD_SIZE; i++ {
		for j := 0; j < FIELD_SIZE; j++ {
			fmt.Printf("%d", field[i][j])
		}
		fmt.Printf("\n")
	}
}

func (e *Environment) StartLoop(usePrint bool) {
	for {
		e.SetPannas()
		e.UpdateArrow()
		if usePrint {
			e.PrintField()
		}
		if !e.TryMakeAction(e.homa) {
			e.homa.AddScore(PENALTY_SCORE)
			fmt.Printf("Game over. Score %d\n", e.homa.GetScore())
			break
		}
		// Loop escape
		if len(e.goldPoses) == 0 {
			break
		}
	}
	fmt.Printf("Your score: %d\n", e.homa.GetScore())
}

func (e *Environment) SetPannas() {
	e.pannaPoses = make([]Pos, 0, 5)

	for i := uint(0); i < e.pannaCount; i++ {
		// If random value is less than chance
		rVal := rand.Float64()
		if rVal < SPAWN_CHANCE {
			e.pannaPoses = append(e.pannaPoses)
		}
	}
}

func (e *Environment) UpdateArrow() bool {
	for {
		if e.arrow.pos == e.wiyPos {
			e.arrow.isLaunched = false
			e.wiyPos = Pos{
				x: -1000,
				y: -1000,
			}
			return true
		}

		pos := GetMovement(e.arrow.dir)
		e.arrow.pos.x += pos.x
		e.arrow.pos.y += pos.y

		// Loop escape
		if e.arrow.isLaunched && InBounds(e.arrow.pos) {
			continue
		} else {
			break
		}
	}
	e.arrow.isLaunched = false
	return false
}

func (e *Environment) GetRandomFreePosExclude(poses []Pos) Pos {
	var pos Pos
	for {
		pos = Pos{
			x: rand.Intn(FIELD_SIZE - 1),
			y: rand.Intn(FIELD_SIZE - 1),
		}

		if IsPosFree(poses, pos) && e.wiyPos != pos && e.homa.homaPos != pos {
			return pos
		}
	}
	return Pos{}
}

func (e *Environment) HasSmell(pos Pos) bool {
	return math.Abs(float64(pos.x) - float64(e.wiyPos.x)) + math.Abs(float64(pos.y) - float64(e.wiyPos.y)) <= 1
}

func (e *Environment) HasGold(pos Pos) bool {
	return !IsPosFree(e.goldPoses, pos)
}

func (e *Environment) ArrowKilled() bool {
	return e.wiyPos == Pos{
		x: -1000,
		y: -1000,
	}
}

func (e *Environment) FeelWind(pos Pos) bool {
	for _, panna := range e.pannaPoses {
		if math.Abs(float64(pos.x) - float64(panna.x)) + math.Abs(float64(pos.y) - float64(panna.y)) <= 1 {
			return true
		}
	}
	return false
}

func (e *Environment) TryMakeAction(homa *Homa) bool {
	if !IsPosFree(e.pannaPoses, homa.homaPos) || e.wiyPos == homa.homaPos {
		return false
	}

	homa.AddScore(ACTION_SCORE)
	receptors := e.GetReceptorsState(e.homa)
	nextAction := e.homa.GetDesiredAction(receptors)

	switch nextAction {
	case MoveUP:
		return e.TryMove(e.homa)
	case PickUp:
		return e.TryPickUp(e.homa)
	case Shoot:
		return e.TryShoot(e.homa)
	default:
		return e.Rotate(e.homa, nextAction)
	}

	return true
}

func (e *Environment) TryMove(homa *Homa) bool {
	homaPos := homa.homaPos
	change := GetMovement(homa.lookDir)
	homaPos.x += change.x
	homaPos.y += change.y

	if !InBounds(homaPos) {
		homa.hit = true
		return true
	}

	homa.homaPos = homaPos
	if !IsPosFree(e.pannaPoses, homaPos) {
		return false
	}

	if homaPos == e.wiyPos {
		return false
	}

	return true
}

func (e *Environment) TryPickUp(homa *Homa) bool {
	if !IsPosFree(e.goldPoses, homa.homaPos) {
		goldPos := FindPosInSlice(e.goldPoses, homa.homaPos)
		e.goldPoses = append(e.goldPoses[:goldPos], e.goldPoses[goldPos+1:]...)
		e.homa.AddScore(GOLD_SCORE)
	}
	return true
}

func (e *Environment) TryShoot(homa *Homa) bool {
	if !homa.isArrowAvailable {
		return true
	}

	homa.AddScore(ARROW_LAUNCH_SCORE)
	e.arrow.isLaunched = true
	e.arrow.dir = homa.lookDir
	e.arrow.pos = homa.homaPos
	homa.isArrowAvailable = false
	return true
}

func (e *Environment) Rotate(homa *Homa, nextAction Action) bool {
	if nextAction == RotateRight {
		homa.lookDir = GetRotationClockwise(homa.lookDir)
	} else if nextAction == RotateLeft {
		for i := 0; i < 3; i++ {
			homa.lookDir = GetRotationClockwise(homa.lookDir)
		}
	}
	return true
}

func (e *Environment) GetReceptorsState(homa *Homa) ReceptorsState {
	pos := homa.homaPos
	return ReceptorsState{
		Smell: e.HasSmell(pos),
		Wind: e.FeelWind(pos),
		Blink: e.HasGold(pos),
		Hit: homa.hit,
		Cry: e.ArrowKilled(),
	}
}

func FindPosInSlice(poses []Pos, pos Pos) int {
	for index, p := range poses {
		if p == pos {
			return index
		}
	}
	return -1
}