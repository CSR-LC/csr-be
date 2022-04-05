// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Equipment equipment
//
// swagger:model Equipment
type Equipment struct {

	// description
	// Example: This is a dog harness.\nWARNING: do not put on cats!
	// Required: true
	Description *string `json:"description"`

	// kind
	// Example: 5
	// Required: true
	Kind *int64 `json:"kind"`

	// location
	// Example: 71
	// Required: true
	Location *int64 `json:"location"`

	// name
	// Example: Dog harness 3000
	// Required: true
	Name *string `json:"name"`

	// photo
	// Example: https://...
	// Required: true
	Photo *string `json:"photo"`

	// rate day
	// Required: true
	RateDay *int64 `json:"rate_day"`

	// rate hour
	// Required: true
	RateHour *int64 `json:"rate_hour"`

	// sku
	// Example: ABC012345678
	// Required: true
	Sku *string `json:"sku"`

	// status
	// Example: 1
	// Required: true
	Status *int64 `json:"status"`
}

// Validate validates this equipment
func (m *Equipment) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKind(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePhoto(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRateDay(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRateHour(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSku(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Equipment) validateDescription(formats strfmt.Registry) error {

	if err := validate.Required("description", "body", m.Description); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateKind(formats strfmt.Registry) error {

	if err := validate.Required("kind", "body", m.Kind); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateLocation(formats strfmt.Registry) error {

	if err := validate.Required("location", "body", m.Location); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validatePhoto(formats strfmt.Registry) error {

	if err := validate.Required("photo", "body", m.Photo); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateRateDay(formats strfmt.Registry) error {

	if err := validate.Required("rate_day", "body", m.RateDay); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateRateHour(formats strfmt.Registry) error {

	if err := validate.Required("rate_hour", "body", m.RateHour); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateSku(formats strfmt.Registry) error {

	if err := validate.Required("sku", "body", m.Sku); err != nil {
		return err
	}

	return nil
}

func (m *Equipment) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this equipment based on context it is used
func (m *Equipment) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Equipment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Equipment) UnmarshalBinary(b []byte) error {
	var res Equipment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
