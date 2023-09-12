package handlers

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/pet_size"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/repositories"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"
)

func SetPetSizeHandler(logger *zap.Logger, api *operations.BeAPI) {
	petSizeRepo := repositories.NewPetSizeRepository()
	petSizeHandler := NewPetSize(logger)

	api.PetSizeGetAllPetSizeHandler = petSizeHandler.GetAllPetSizeFunc(petSizeRepo)
	api.PetSizeEditPetSizeHandler = petSizeHandler.UpdatePetSizeByID(petSizeRepo)
	api.PetSizeDeletePetSizeHandler = petSizeHandler.DeletePetSizeByID(petSizeRepo)
	api.PetSizeCreateNewPetSizeHandler = petSizeHandler.CreatePetSizeFunc(petSizeRepo)
	api.PetSizeGetPetSizeHandler = petSizeHandler.GetPetSizeByID(petSizeRepo)
}

type PetSize struct {
	logger *zap.Logger
}

func NewPetSize(logger *zap.Logger) *PetSize {
	return &PetSize{
		logger: logger,
	}
}

func (ps PetSize) CreatePetSizeFunc(repository domain.PetSizeRepository) pet_size.CreateNewPetSizeHandlerFunc {
	return func(p pet_size.CreateNewPetSizeParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		allPetSizes, err := repository.GetAll(ctx)
		if err != nil {
			ps.logger.Error("Error while getting pet size", zap.Error(err))
			return pet_size.NewCreateNewPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while creating pet size"))
		}
		for _, petSize := range allPetSizes {
			if *p.NewPetSize.Name == petSize.Name {
				ps.logger.Error("Error while creating pet size", zap.Error(err))
				return pet_size.NewCreateNewPetSizeDefault(http.StatusInternalServerError).
					WithPayload(buildInternalErrorPayload("Error while creating pet size: the name already exist"))
			}
		}
		petSize, err := repository.Create(ctx, *p.NewPetSize)
		if err != nil {
			ps.logger.Error("Error while creating pet size", zap.Error(err))
			return pet_size.NewCreateNewPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildBadRequestErrorPayload("Error while creating pet size"))
		}
		if petSize == nil {
			ps.logger.Error("Pet size is nil")
			return pet_size.NewCreateNewPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while creating pet size"))
		}
		id := int64(petSize.ID)
		return pet_size.NewCreateNewPetSizeCreated().WithPayload(&models.PetSizeResponse{
			ID:          &id,
			Name:        &petSize.Name,
			Size:        &petSize.Size,
			IsUniversal: false,
		},
		)
	}
}

func (ps PetSize) GetAllPetSizeFunc(repository domain.PetSizeRepository) pet_size.GetAllPetSizeHandlerFunc {
	return func(p pet_size.GetAllPetSizeParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		petSizes, err := repository.GetAll(ctx)
		if err != nil {
			ps.logger.Error("Error while getting pet size", zap.Error(err))
			return pet_size.NewGetAllPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while getting pet size"))
		}
		if len(petSizes) == 0 {
			return pet_size.NewGetAllPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("No pet size found"))
		}
		listOfPetSize := models.ListOfPetSizes{}
		for _, v := range petSizes {
			id := int64(v.ID)
			listOfPetSize = append(listOfPetSize, &models.PetSizeResponse{
				ID:          &id,
				Name:        &v.Name,
				Size:        &v.Size,
				IsUniversal: v.IsUniversal,
			})
		}
		return pet_size.NewGetAllPetSizeOK().WithPayload(listOfPetSize)
	}
}

func (ps PetSize) GetPetSizeByID(repo domain.PetSizeRepository) pet_size.GetPetSizeHandlerFunc {
	return func(p pet_size.GetPetSizeParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		petSize, err := repo.GetByID(ctx, int(p.PetSizeID))
		if err != nil {
			ps.logger.Error("Error while getting pet size by id", zap.Error(err))
			return pet_size.NewGetPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while getting pet size"))
		}
		if petSize == nil {
			return pet_size.NewGetPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while getting pet size"))
		}
		id := int64(petSize.ID)
		return pet_size.NewGetPetSizeOK().WithPayload(&models.PetSizeResponse{ID: &id, Name: &petSize.Name, Size: &petSize.Size})
	}
}

func (ps PetSize) DeletePetSizeByID(repo domain.PetSizeRepository) pet_size.DeletePetSizeHandlerFunc {
	return func(p pet_size.DeletePetSizeParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		err := repo.Delete(ctx, int(p.PetSizeID))
		if err != nil {
			ps.logger.Error("Error while deleting pet size by id", zap.Error(err))
			return pet_size.NewDeletePetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while deleting pet size"))
		}
		return pet_size.NewDeletePetSizeOK().WithPayload("Pet size deleted")
	}
}

func (ps PetSize) UpdatePetSizeByID(repo domain.PetSizeRepository) pet_size.EditPetSizeHandlerFunc {
	return func(p pet_size.EditPetSizeParams, _ *models.Principal) middleware.Responder {
		ctx := p.HTTPRequest.Context()
		petSize, err := repo.Update(ctx, int(p.PetSizeID), p.EditPetSize)
		if err != nil {
			ps.logger.Error("Error while updating pet size by id", zap.Error(err))
			return pet_size.NewEditPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while updating pet size"))
		}
		if petSize == nil {
			ps.logger.Error("Error while updating pet size by id", zap.Error(err))
			return pet_size.NewEditPetSizeDefault(http.StatusInternalServerError).
				WithPayload(buildInternalErrorPayload("Error while updating pet size"))
		}

		id := int64(petSize.ID)
		return pet_size.NewEditPetSizeOK().WithPayload(&models.PetSizeResponse{ID: &id, Name: &petSize.Name, Size: &petSize.Size})
	}
}
