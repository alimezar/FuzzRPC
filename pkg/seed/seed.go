// pkg/seed/seed.go
package seed

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

// BuildSeed constructs a dynamic.Message populated with minimal default values for every field.
func BuildSeed(md *desc.MessageDescriptor) *dynamic.Message {
	msg := dynamic.NewMessage(md)
	for _, f := range md.GetFields() {
		// Skip map and repeated fields
		if f.IsMap() || f.IsRepeated() {
			continue
		}
		// Nested message: recurse
		if nested := f.GetMessageType(); nested != nil {
			msg.SetField(f, BuildSeed(nested))
			continue
		}
		// Enum: pick first value
		if enum := f.GetEnumType(); enum != nil {
			vals := enum.GetValues()
			if len(vals) > 0 {
				msg.SetField(f, vals[0].GetNumber())
			}
			continue
		}
		// Scalar types: zero/default values
		switch f.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			msg.SetField(f, false)
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			msg.SetField(f, "")
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			msg.SetField(f, []byte{})
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
			msg.SetField(f, float32(0))
		case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			msg.SetField(f, float64(0))
		case descriptorpb.FieldDescriptorProto_TYPE_INT32,
			descriptorpb.FieldDescriptorProto_TYPE_SINT32,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			msg.SetField(f, int32(0))
		case descriptorpb.FieldDescriptorProto_TYPE_INT64,
			descriptorpb.FieldDescriptorProto_TYPE_SINT64,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			msg.SetField(f, int64(0))
		case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
			msg.SetField(f, uint32(0))
		case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
			descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			msg.SetField(f, uint64(0))
		}
	}
	return msg
}
