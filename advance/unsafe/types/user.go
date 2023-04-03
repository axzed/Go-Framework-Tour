package types

type User struct {
	Name    string
	age     int32
	Alias   []byte
	Address string
}

type UserV1 struct {
	Name    string
	age     int32
	agev1   int32
	Alias   []byte
	Address string
}

type UserV2 struct {
	Name    string
	Alias   []byte
	Address string
	age     int32
}
