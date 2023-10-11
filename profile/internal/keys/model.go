package keys

import "time"

type Type string

const (
	Cpf    Type = "cpf"
	Phone  Type = "phone"
	Email  Type = "email"
	Random Type = "random"
)

type Key struct {
	PK        int
	AccountID int `dynamodbav:"SK"`
	Name      string
	Type      Type
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type KeyRequest struct {
	keyID int
}

type ListKeyRequest struct {
	keyIDs []int
}
