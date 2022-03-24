// Code generated by go-swagger; DO NOT EDIT.

package kinds

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// DeleteKindCreatedCode is the HTTP code returned for type DeleteKindCreated
const DeleteKindCreatedCode int = 201

/*DeleteKindCreated kind of equipment successfully deleted from db

swagger:response deleteKindCreated
*/
type DeleteKindCreated struct {

	/*
	  In: Body
	*/
	Payload *models.DeleteKindResponse `json:"body,omitempty"`
}

// NewDeleteKindCreated creates DeleteKindCreated with default headers values
func NewDeleteKindCreated() *DeleteKindCreated {

	return &DeleteKindCreated{}
}

// WithPayload adds the payload to the delete kind created response
func (o *DeleteKindCreated) WithPayload(payload *models.DeleteKindResponse) *DeleteKindCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete kind created response
func (o *DeleteKindCreated) SetPayload(payload *models.DeleteKindResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKindCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteKindDefault Unexpected error.

swagger:response deleteKindDefault
*/
type DeleteKindDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteKindDefault creates DeleteKindDefault with default headers values
func NewDeleteKindDefault(code int) *DeleteKindDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteKindDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete kind default response
func (o *DeleteKindDefault) WithStatusCode(code int) *DeleteKindDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete kind default response
func (o *DeleteKindDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete kind default response
func (o *DeleteKindDefault) WithPayload(payload *models.Error) *DeleteKindDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete kind default response
func (o *DeleteKindDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteKindDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}