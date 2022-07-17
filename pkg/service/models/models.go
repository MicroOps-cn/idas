package models

type AuthAlgorithm string

func (x AppMeta_GrantType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}

func (x AppMeta_GrantMode) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}

func (x AppMeta_Status) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
