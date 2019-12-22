package main

type Action int
const (
	MoveUP Action = iota
	RotateLeft
	RotateRight
	PickUp
	Shoot
)

type Pos struct {
	x int
	y int
}

type LookDirection int
const (
	UP LookDirection = iota
	RIGHT
	DOWN
	LEFT
)
