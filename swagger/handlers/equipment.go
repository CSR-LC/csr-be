package handlers

import (
	"fmt"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/equipment"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/restapi/operations/status"
	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Equipment struct {
	client *ent.Client
	logger *zap.Logger
}

func NewEquipment(client *ent.Client, logger *zap.Logger) *Equipment {
	return &Equipment{
		client: client,
		logger: logger,
	}
}

func (c Equipment) PostEquipmentFunc() equipment.CreateNewEquipmentHandlerFunc {
	return func(s equipment.CreateNewEquipmentParams) middleware.Responder {
		e, err := c.client.Equipment.Create().SetName(*s.NewEquipment.Name).SetDescription(*s.NewEquipment.Description).SetSku(*s.NewEquipment.Sku).SetRateDay(*s.NewEquipment.RateDay).SetRateHour(*s.NewEquipment.RateHour).Save(s.HTTPRequest.Context())
		if err != nil {
			return equipment.NewCreateNewEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}

		id := fmt.Sprintf("%d", e.ID)
		return equipment.NewCreateNewEquipmentCreated().WithPayload(&models.EquipmentResponse{
			ID:          &id,
			Description: &e.Description,
			Name:        &e.Name,
			RateDay:     &e.RateDay,
			RateHour:    &e.RateHour,
			Sku:         &e.Sku,
		})
	}
}

func (c Equipment) GetEquipmentFunc() equipment.GetEquipmentHandlerFunc {
	return func(s equipment.GetEquipmentParams) middleware.Responder {
		equipmentId, err := strconv.Atoi(s.EquipmentID)
		if err != nil {
			return equipment.NewGetEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		e, err := c.client.Equipment.Get(s.HTTPRequest.Context(), equipmentId)
		if err != nil {
			return equipment.NewGetEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}

		return equipment.NewGetEquipmentCreated().WithPayload(&models.Equipment{
			Description: &e.Description,
			Name:        &e.Name,
			RateDay:     &e.RateDay,
			RateHour:    &e.RateHour,
			Sku:         &e.Sku,
		})
	}
}

func (c Equipment) DeleteEquipmentFunc() equipment.DeleteEquipmentHandlerFunc {
	return func(s equipment.DeleteEquipmentParams) middleware.Responder {
		equipmentId, err := strconv.Atoi(s.EquipmentID)
		if err != nil {
			return status.NewDeleteStatusDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		e, err := c.client.Equipment.Get(s.HTTPRequest.Context(), equipmentId)
		if err != nil {
			return equipment.NewDeleteEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		err = c.client.Equipment.DeleteOne(e).Exec(s.HTTPRequest.Context())
		if err != nil {
			return equipment.NewDeleteEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}

		return equipment.NewDeleteEquipmentCreated().WithPayload(&models.Equipment{
			Description: &e.Description,
			Name:        &e.Name,
			RateDay:     &e.RateDay,
			RateHour:    &e.RateHour,
			Sku:         &e.Sku,
		})
	}
}

func (c Equipment) ListEquipmentFunc() equipment.GetAllEquipmentHandlerFunc {
	return func(s equipment.GetAllEquipmentParams) middleware.Responder {
		e, err := c.client.Equipment.Query().All(s.HTTPRequest.Context())
		if err != nil {
			return equipment.NewGetAllEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		listEquipment := models.ListEquipment{}
		for _, element := range e {
			id := strconv.Itoa(element.ID)
			listEquipment = append(listEquipment, &models.EquipmentResponse{
				ID:          &id,
				Name:        &element.Name,
				Description: &element.Description,
				Sku:         &element.Sku,
				RateDay:     &element.RateDay,
				RateHour:    &element.RateHour,
			})
		}
		return equipment.NewGetAllEquipmentCreated().WithPayload(listEquipment)
	}
}

func (c Equipment) EditEquipmentFunc() equipment.EditEquipmentHandlerFunc {
	return func(s equipment.EditEquipmentParams) middleware.Responder {
		equipmentId, err := strconv.Atoi(s.EquipmentID)
		if err != nil {
			return equipment.NewEditEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		e, err := c.client.Equipment.Get(s.HTTPRequest.Context(), equipmentId)
		if err != nil {
			return equipment.NewEditEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}
		edit := e.Update()
		if *s.EditEquipment.Name != "" {
			edit.SetName(*s.EditEquipment.Name)
		}
		if *s.EditEquipment.Sku != "" {
			edit.SetSku(*s.EditEquipment.Sku)
		}
		if *s.EditEquipment.Description != "" {
			edit.SetDescription(*s.EditEquipment.Description)
		}
		if *s.EditEquipment.RateDay != 0 {
			edit.SetRateDay(*s.EditEquipment.RateDay)
		}
		if *s.EditEquipment.RateHour != 0 {
			edit.SetRateHour(*s.EditEquipment.RateHour)
		}
		res, err := edit.Save(s.HTTPRequest.Context())
		//res, err := c.client.Equipment.Get(s.HTTPRequest.Context(), equipmentId)
		if err != nil {
			return equipment.NewEditEquipmentDefault(http.StatusInternalServerError).WithPayload(&models.Error{
				Data: &models.ErrorData{
					Message: err.Error(),
				},
			})
		}

		return equipment.NewEditEquipmentCreated().WithPayload(&models.Equipment{
			Description: &res.Description,
			Name:        &res.Name,
			RateDay:     &res.RateDay,
			RateHour:    &res.RateHour,
			Sku:         &res.Sku,
		})
	}
}
