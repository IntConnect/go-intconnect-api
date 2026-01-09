package mapper

import (
	"go-intconnect-api/internal/entity"
	"go-intconnect-api/internal/model"
	"go-intconnect-api/pkg/helper"
)

func MapMachineDocument(machineEntity *entity.Machine, machineResponse *model.MachineResponse) {
	if machineEntity == nil || machineResponse == nil {
		return
	}

	if len(machineEntity.MachineDocuments) == 0 {
		machineResponse.MachineDocuments = []*model.MachineDocumentResponse{}
		return
	}

	result := make([]*model.MachineDocumentResponse, 0, len(machineEntity.MachineDocuments))

	for i := range machineEntity.MachineDocuments {
		machineDocument := machineEntity.MachineDocuments[i]

		mapped := helper.MapEntityIntoResponse[*entity.MachineDocument, *model.MachineDocumentResponse](
			machineDocument)

		result = append(result, mapped)
	}

	machineResponse.MachineDocuments = result
}
