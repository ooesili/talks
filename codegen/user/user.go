package user

// take this
// +getters // HL
type User struct {
	name     string
	password string
}

// generate these
func (u User) Name() string { // HL
	return u.name // HL
} // HL

func (u User) Password() string { // HL
	return u.password // HL
} // HL
