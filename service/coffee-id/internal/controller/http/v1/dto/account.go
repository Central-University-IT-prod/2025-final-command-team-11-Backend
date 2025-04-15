package dto

type Account struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateUser struct {
	Email    string `json:"email"    validate:"required,email"`
	Name     string `json:"name"     validate:"required,min=2,max=30"`
	Password string `json:"password" validate:"required,min=8,max=50,password"`
}

type UpdateUser struct {
	Email       string `json:"email"       validate:"omitempty,email"`
	Name        string `json:"name"        validate:"omitempty,min=2,max=30"`
	Password    string `json:"password"    validate:"omitempty,min=8,max=50,password"`
	OldPassword string `json:"oldPassword" validate:"omitempty,min=8,max=50,password"`
}

type SetRole struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

type AP struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Verified bool   `json:"verified"`
}

type AcountsPagintation struct {
	Accounts []*AP `json:"accounts"`
	Count    int   `json:"count"`
}

type DeleteUser struct {
	Password string `json:"password" validate:"required,min=8,max=50,password"`
}

type AccountAnswer struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
