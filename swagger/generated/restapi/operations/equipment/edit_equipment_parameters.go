// Code generated by go-swagger; DO NOT EDIT.

package equipment

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	"git.epam.com/epm-lstr/epm-lstr-lc/be/swagger/generated/models"
)

// NewEditEquipmentParams creates a new EditEquipmentParams object
//
// There are no default values defined in the spec.
func NewEditEquipmentParams() EditEquipmentParams {

	return EditEquipmentParams{}
}

// EditEquipmentParams contains all the bound params for the edit equipment operation
// typically these are obtained from a http.Request
//
// swagger:parameters EditEquipment
type EditEquipmentParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Edit an equipment
	  In: body
	*/
	EditEquipment *models.Equipment
	/*equipment id
	  Required: true
	  In: path
	*/
	EquipmentID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewEditEquipmentParams() beforehand.
func (o *EditEquipmentParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Equipment
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("editEquipment", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.EditEquipment = &body
			}
		}
	}

	rEquipmentID, rhkEquipmentID, _ := route.Params.GetOK("equipmentId")
	if err := o.bindEquipmentID(rEquipmentID, rhkEquipmentID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindEquipmentID binds and validates parameter EquipmentID from path.
func (o *EditEquipmentParams) bindEquipmentID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.EquipmentID = raw

	return nil
}
