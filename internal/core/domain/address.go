package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	TypeResidential Type = "residential"
	TypeCommercial  Type = "commercial"
	TypeShipping    Type = "shipping"
	TypeBilling     Type = "billing"
)

type Type string

type Address struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Street     string
	Number     string
	Complement string
	City       string
	StateID    string
	ZipCode    string
	Type       string
	ID         uuid.UUID
	UserID     uuid.UUID
	IsDefault  bool
}

type State struct {
	ID     string
	Name   string
	Region string
}

func (a *Address) Update(street, number, complement, city, stateID, zipCode, addressType string) error {
	if street != "" {
		a.Street = street
	}
	if number != "" {
		a.Number = number
	}
	a.Complement = complement
	if city != "" {
		a.City = city
	}
	if stateID != "" {
		a.StateID = stateID
	}
	if zipCode != "" {
		a.ZipCode = zipCode
	}
	if addressType != "" {
		a.Type = addressType
	}
	a.UpdatedAt = time.Now().UTC()
	return nil
}

func (a *Address) SetDefault(isDefault bool) {
	a.IsDefault = isDefault
	a.UpdatedAt = time.Now().UTC()
}

func (a *Address) FormattedAddress() string {
	complement := ""
	if a.Complement != "" {
		complement = ", " + a.Complement
	}

	return a.Street + ", " + a.Number + complement + " - " + ", " + a.City + " - " + a.StateID +
		", " + a.ZipCode
}

func NewAddress(address Address) (*Address, error) {
	if address.Street == "" {
		return nil, errors.New("street cannot be empty")
	}
	if address.Number == "" {
		return nil, errors.New("number cannot be empty")
	}
	if address.City == "" {
		return nil, errors.New("city cannot be empty")
	}
	if address.StateID == "" {
		return nil, errors.New("state cannot be empty")
	}
	if address.ZipCode == "" {
		return nil, errors.New("zip code cannot be empty")
	}

	if address.Type == "" {
		address.Type = string(TypeResidential)
	}

	now := time.Now().UTC()
	return &Address{
		ID:         uuid.New(),
		UserID:     address.UserID,
		Street:     address.Street,
		Number:     address.Number,
		Complement: address.Complement,
		City:       address.City,
		StateID:    address.StateID,
		ZipCode:    address.ZipCode,
		IsDefault:  address.IsDefault,
		Type:       address.Type,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
