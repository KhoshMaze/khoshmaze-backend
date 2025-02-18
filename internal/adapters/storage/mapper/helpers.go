package mapper

import "github.com/KhoshMaze/khoshmaze-backend/internal/domain/user/model"

func IntegerToSubscribtionType(role uint) model.Subscribtion {
	switch role {
	case 1:
		return model.Premium1
	case 2:
		return model.Premium2

	default:
		return model.Normal
	}

}
