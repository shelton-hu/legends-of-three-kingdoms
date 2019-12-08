package manager

import (
	"strings"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/ptypes/wrappers"
)

const (
	TagName = "json"
)

type Request interface {
	Reset()
	String() string
	ProtoMessage()
	Descriptor() ([]byte, []int)
}

type RequestWithSort interface {
	Request
	GetSortKey() *wrappers.StringValue
	GetReverse() *wrappers.BoolValue
}

func getFieldName(field *structs.Field) string {
	tag := field.Tag(TagName)
	t := strings.Split(tag, ",")
	if len(t) == 0 {
		return "-"
	}
	return t[0]
}
