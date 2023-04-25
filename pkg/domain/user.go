package domain

type User struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"uniqueIndex;not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null,email"`
	Password  string `json:"password" gorm:"not null"`
	Age       int    `json:"age" gorm:"not null"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//==================================================================================================
// Repository
//==================================================================================================

type UserRepository interface {
	InsertUser(user User) (err error)
	FindUserByEmail(email string) (User, error)
	FindUserByID(id int64) (User, error)
}
