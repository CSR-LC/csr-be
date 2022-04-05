// Code generated by go-swagger; DO NOT EDIT.

package equipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// CreateNewEquipmentCreatedCode is the HTTP code returned for type CreateNewEquipmentCreated
const CreateNewEquipmentCreatedCode int = 201

/*CreateNewEquipmentCreated Equipment created

swagger:response createNewEquipmentCreated
*/
type CreateNewEquipmentCreated struct {

	/*
	  In: Body
	*/
	Payload *models.EquipmentResponse `json:"body,omitempty"`
}

// NewCreateNewEquipmentCreated creates CreateNewEquipmentCreated with default headers values
func NewCreateNewEquipmentCreated() *CreateNewEquipmentCreated {

	return &CreateNewEquipmentCreated{}
}

// WithPayload adds the payload to the create new equipment created response
func (o *CreateNewEquipmentCreated) WithPayload(payload *models.EquipmentResponse) *CreateNewEquipmentCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create new equipment created response
func (o *CreateNewEquipmentCreated) SetPayload(payload *models.EquipmentResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNewEquipmentCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*CreateNewEquipmentDefault Unexpected error.

swagger:response createNewEquipmentDefault
*/
type CreateNewEquipmentDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewCreateNewEquipmentDefault creates CreateNewEquipmentDefault with default headers values
func NewCreateNewEquipmentDefault(code int) *CreateNewEquipmentDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateNewEquipmentDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create new equipment default response
func (o *CreateNewEquipmentDefault) WithStatusCode(code int) *CreateNewEquipmentDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create new equipment default response
func (o *CreateNewEquipmentDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create new equipment default response
func (o *CreateNewEquipmentDefault) WithPayload(payload *models.Error) *CreateNewEquipmentDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create new equipment default response
func (o *CreateNewEquipmentDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateNewEquipmentDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
