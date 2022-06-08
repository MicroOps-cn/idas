package api

type CustomType interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) (err error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) (err error)
}
