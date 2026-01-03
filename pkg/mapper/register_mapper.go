package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/helper"
)

func FuncMapRegister(registerEntity *entity.Register, registerResponse *model.RegisterResponse) {
	registerResponse.MachineResponse = helper.MapEntityIntoResponse[
		*entity.Machine,
		*model.MachineResponse,
	](registerEntity.Machine)
}
