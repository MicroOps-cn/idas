package fs

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
)

func (x *FsPath) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, b []byte) error {
	return json.Unmarshal(b, &x.FsPath)
}
