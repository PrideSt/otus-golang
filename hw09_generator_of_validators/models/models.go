package models

//go:generate go run ./../go-validate/main.go -in ./models.go

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}
)

type Token struct {
	Header    []byte
	Payload   []byte
	Signature []byte
}

type Response struct {
	Code int    `validate:"in:200,404,500"`
	Body string `json:"omitempty"`
}

type WithEmbedded struct {
	ID string `validate:"len:36"`
	Response
}

type ClassProperty struct {
	ID string `validate:"len:36"`
	// OneResp is instance of Validate interface
	OneResp          Response
	NoValidatorClass Token
}
