package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/models"
	eqStatus "git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations/equipment_status"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/pkg/domain"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/ent/enttest"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/mocks"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi"
	"git.epam.com/epm-lstr/epm-lstr-lc/be/internal/generated/swagger/restapi/operations"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

func TestSetEquipmentStatusHandler(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:eqstatushandler?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	logger := zap.NewNop()

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		t.Fatal(err)
	}
	api := operations.NewBeAPI(swaggerSpec)
	SetEquipmentStatusHandler(logger, api)
	assert.NotEmpty(t, api.EquipmentStatusCheckEquipmentStatusHandler)
	assert.NotEmpty(t, api.EquipmentStatusUpdateEquipmentStatusOnAvailableHandler)
	assert.NotEmpty(t, api.EquipmentStatusUpdateEquipmentStatusOnUnavailableHandler)
	assert.NotEmpty(t, api.EquipmentStatusUpdateRepairedEquipmentStatusDatesHandler)
}

type EquipmentStatusTestSuite struct {
	suite.Suite
	logger                    *zap.Logger
	equipmentStatusRepository *mocks.EquipmentStatusRepository
	orderStatusRepository     *mocks.OrderStatusRepository
	handler                   *EquipmentStatus
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(EquipmentStatusTestSuite))
}

func (s *EquipmentStatusTestSuite) SetupTest() {
	s.logger = zap.NewNop()
	s.equipmentStatusRepository = &mocks.EquipmentStatusRepository{}
	s.orderStatusRepository = &mocks.OrderStatusRepository{}
	s.handler = NewEquipmentStatus(s.logger)
}

func (s *EquipmentStatusTestSuite) Test_Put_EquipmentStatusInRepairFunc_OK() {
	t := s.T()
	request := http.Request{}
	ctx := request.Context()

	statusName := "statusName"
	timeNow := time.Now()
	startDate := time.Date(2023, time.February, 14, 12, 34, 56, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 10)

	data := eqStatus.UpdateEquipmentStatusOnUnavailableParams{
		HTTPRequest:       &request,
		EquipmentstatusID: 1,
		Name: &models.EquipmentStatusInRepairRequest{
			EndDate:    (*strfmt.DateTime)(&endDate),
			StartDate:  (*strfmt.DateTime)(&startDate),
			StatusName: &statusName,
		},
	}

	// err := errors.New("test")

	eqStatusModel := models.EquipmentStatus{
		StartDate:  (*strfmt.DateTime)(&startDate),
		EndDate:    (*strfmt.DateTime)(&endDate),
		StatusName: &statusName,
		ID:         &data.EquipmentstatusID,
	}

	eqStatusResponseModel := ent.EquipmentStatus{
		ID:        int(data.EquipmentstatusID),
		StartDate: startDate,
		EndDate:   endDate,
	}

	s.equipmentStatusRepository.On("Update", ctx, &eqStatusModel).Return(&eqStatusResponseModel, nil)

	orderResult := ent.Order{ID: 1}
	userResult := ent.User{ID: 1}

	s.equipmentStatusRepository.On(
		"GetOrderAndUserByDates",
		ctx,
		int(*eqStatusModel.ID),
		time.Time(*eqStatusModel.StartDate),
		time.Time(*eqStatusModel.EndDate),
	).Return(&orderResult, &userResult, nil)

	comment := EQUIPMENT_UNDER_REPAIR_COMMENT_FOR_ORDER
	orderID := int64(orderResult.ID)

	orderModel := models.NewOrderStatus{
		Comment:   &comment,
		CreatedAt: (*strfmt.DateTime)(&timeNow),
		OrderID:   &orderID,
		Status:    &domain.OrderStatusRejected,
	}

	s.orderStatusRepository.On("UpdateStatus", ctx, userResult.ID, orderModel).Return(nil)

	s.equipmentStatusRepository.On(
		"GetEquipmentStatusByID",
		ctx,
		int(*eqStatusModel.ID),
	).Return(&eqStatusResponseModel, nil)

	handlerFunc := s.handler.PutEquipmentStatusInRepairFunc(
		s.equipmentStatusRepository, s.orderStatusRepository,
	)

	access := "dummy access"
	resp := handlerFunc(data, access)
	responseRecorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(responseRecorder, producer)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	s.equipmentStatusRepository.AssertExpectations(t)

	actualEquipment := &models.EquipmentStatusRepairResponse{}
	err := json.Unmarshal(responseRecorder.Body.Bytes(), actualEquipment)
	if err != nil {
		t.Errorf("unable to unmarshal response body: %v", err)
	}
	// assert.Equal(t, equipmentToReturn.Name, *actualEquipment.Name)
}
