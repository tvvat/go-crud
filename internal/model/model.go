package model

type User struct {
	Id       int     `json:"id" gorm:"primary_key"`
	Login    string  `json:"login" `
	Password string  `json:"password"`
	Groups   []Group `json:"groupid" gorm:"many2many:user_group;"`
}

type Group struct {
	Id        int     `json:"id" `
	Subgroups []Group `json:"groupid" gorm:"many2many:group_subgroups;"`
}
