package domain

import (
	"fmt"
	"time"

	"github.com/MatheusHenrique129/bemax-api/pkg/hash"
	"github.com/google/uuid"
	"github.com/klassmann/cpfcnpj"
)

const (
	MinNameLen     = 3
	MinPasswordLen = 6

	UserStatusActive   Status = "active"
	UserStatusInactive Status = "inactive"
	UserStatusBlocked  Status = "blocked"
)

type Status string

type User struct {
	BirthDate time.Time  `json:"birth_date"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
	Password  string     `json:"password"`
	Phone     string     `json:"phone"`
	Status    Status     `json:"status"`
	CPF       string     `json:"cpf"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Addresses []Address  `json:"addresses"`
	Roles     []Role     `json:"roles"`
	ID        uuid.UUID  ``
}

func (u *User) Activate() {
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Deactivate() {
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) Block() {
	u.Status = UserStatusBlocked
	u.UpdatedAt = time.Now().UTC()
}

func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

func (u *User) CheckPassword(password string) error {
	return hash.CheckPassword(password, u.Password)
}

func (u *User) UpdateLastLogin() {
	now := time.Now().UTC()
	u.LastLogin = &now
}

func (u *User) MaskCPF() {
	plainCPF := cpfcnpj.Clean(u.CPF)

	maskedCPF := fmt.Sprintf("***.***.%s-%s",
		plainCPF[6:9],
		plainCPF[9:11],
	)

	u.CPF = maskedCPF
}

func (u *User) Update(email, fullName, phone string) {
	if email != "" {
		u.Email = email
	}
	if fullName != "" {
		u.FullName = fullName
	}
	if phone != "" {
		u.Phone = phone
	}
	u.UpdatedAt = time.Now().UTC()
}
