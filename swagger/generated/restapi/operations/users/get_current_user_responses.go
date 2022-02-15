// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// GetCurrentUserOKCode is the HTTP code returned for type GetCurrentUserOK
const GetCurrentUserOKCode int = 200

/*GetCurrentUserOK Success

swagger:response getCurrentUserOK
*/
type GetCurrentUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.GetUserResponse `json:"body,omitempty"`
}

// NewGetCurrentUserOK creates GetCurrentUserOK with default headers values
func NewGetCurrentUserOK() *GetCurrentUserOK {

	return &GetCurrentUserOK{}
}

// WithPayload adds the payload to the get current user o k response
func (o *GetCurrentUserOK) WithPayload(payload *models.GetUserResponse) *GetCurrentUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get current user o k response
func (o *GetCurrentUserOK) SetPayload(payload *models.GetUserResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCurrentUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetCurrentUserDefault Unexpected error.

swagger:response getCurrentUserDefault
*/
type GetCurrentUserDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetCurrentUserDefault creates GetCurrentUserDefault with default headers values
func NewGetCurrentUserDefault(code int) *GetCurrentUserDefault {
	if code <= 0 {
		code = 500
	}

	return &GetCurrentUserDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get current user default response
func (o *GetCurrentUserDefault) WithStatusCode(code int) *GetCurrentUserDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get current user default response
func (o *GetCurrentUserDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get current user default response
func (o *GetCurrentUserDefault) WithPayload(payload *models.Error) *GetCurrentUserDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get current user default response
func (o *GetCurrentUserDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCurrentUserDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
