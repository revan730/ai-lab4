package main

type ReceptorsState struct {
	Smell bool
	Wind bool
	Blink bool
	Hit bool
	Cry bool
}

type Homa struct {
	hit bool
	isArrowAvailable bool
	score int
	homaPos Pos
	lookDir LookDirection
	movementDir LookDirection
	nextMovementDir LookDirection
	facedWiy bool
	actionSequence []LookDirection
}

func NewHoma() *Homa {
	var h Homa
	h.hit = false
	h.isArrowAvailable = true
	h.score = 0
	h.homaPos = Pos{
		x: 0,
		y: 0,
	}
	h.lookDir = RIGHT
	h.movementDir = RIGHT
	h.facedWiy = false

	return &h
}

func (h *Homa) GetDesiredAction(state ReceptorsState) Action {
	h.hit = false
	if state.Blink {
		return PickUp
	}

	if state.Wind {
		if h.movementDir == h.lookDir {
			return RotateRight
		} else {
			h.restoreDirection()
		}
	}

	if state.Smell && !state.Cry && h.isArrowAvailable {
		if h.homaPos.y != 0 || h.homaPos.y == 0 && h.homaPos.x == 0 || h.homaPos.y == 3 {
			if len(h.actionSequence) >= 1 {
				if h.lookDir == RIGHT {
					return Shoot
				} else {
					return RotateLeft
				}
			}

			// Line 35
			if h.lookDir == DOWN {
				return Shoot
			} else if h.lookDir == RIGHT {
				return RotateRight
			} else {
				return  RotateLeft
			}
		}
		if !h.facedWiy {
			h.facedWiy = true
			h.movementDir = LEFT
			h.actionSequence = append(h.actionSequence, RIGHT)
			h.actionSequence = append(h.actionSequence, UP)
			h.actionSequence = append(h.actionSequence, DOWN)
			return h.restoreDirection()
		} else if len(h.actionSequence) == 0 {
			if h.movementDir == h.lookDir {
				return Shoot
			}
			return h.restoreDirection()
		}
	}

	if state.Hit {
		var rotateAction Action
		if h.lookDir == RIGHT {
			rotateAction = RotateRight
		} else {
			rotateAction = RotateLeft
		}
		var appendAction LookDirection
		if h.movementDir == RIGHT {
			appendAction = LEFT
		} else {
			appendAction = RIGHT
		}
		h.actionSequence = append(h.actionSequence, appendAction)
		h.movementDir = DOWN
		return rotateAction
	}

	if h.movementDir == h.lookDir {
		if len(h.actionSequence) != 0 {
			h.movementDir = h.actionSequence[len(h.actionSequence) - 1]
			h.actionSequence = h.actionSequence[:len(h.actionSequence) - 1]
		}
		return MoveUP
	}
	return h.restoreDirection()
}

func (h *Homa) AddScore(score int) {
	h.score += score
}

func (h *Homa) GetScore() int {
	return h.score
}

func (h *Homa) restoreDirection() Action {
	curLookDir := h.lookDir
	desiredLookDir := h.movementDir
	if curLookDir - desiredLookDir == 1 || desiredLookDir - curLookDir == 3 {
		return RotateLeft
	} else if desiredLookDir - curLookDir == 1 || curLookDir - desiredLookDir == 3 {
		return RotateRight
	}
	return RotateRight
}