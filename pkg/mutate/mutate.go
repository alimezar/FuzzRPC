// pkg/mutate/mutate.go
package mutate

import (
	"fmt"
	"github.com/jhump/protoreflect/dynamic"
	"math"
	mrand "math/rand"

	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

// CloneMessage clones a dynamic.Message via binary marshal/unmarshal.
func CloneMessage(msg *dynamic.Message) (*dynamic.Message, error) {
	buf, err := msg.Marshal()
	if err != nil {
		return nil, fmt.Errorf("clone marshal: %w", err)
	}
	clone := dynamic.NewMessage(msg.GetMessageDescriptor())
	if err := clone.Unmarshal(buf); err != nil {
		return nil, fmt.Errorf("clone unmarshal: %w", err)
	}
	return clone, nil
}

// MutateSeed returns a slice of mutated messages (one mutation per field).
func MutateSeed(seed *dynamic.Message) ([]*dynamic.Message, error) {
	md := seed.GetMessageDescriptor()
	var muts []*dynamic.Message

	for _, f := range md.GetFields() {
		// skip maps and repeated fields
		if f.IsMap() || f.IsRepeated() {
			continue
		}
		base := seed.GetField(f)

		// STRING mutations
		if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_STRING {
			variants := []string{"A", RandString(64), "' OR '1'='1"}
			for _, v := range variants {
				m, err := CloneMessage(seed)
				if err != nil {
					return nil, err
				}
				m.SetField(f, v)
				muts = append(muts, m)
			}
		}

		// BOOL flip
		if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_BOOL {
			m, err := CloneMessage(seed)
			if err != nil {
				return nil, err
			}
			m.SetField(f, !base.(bool))
			muts = append(muts, m)
		}

		// ENUM: all values
		if enum := f.GetEnumType(); enum != nil {
			for _, ev := range enum.GetValues() {
				m, err := CloneMessage(seed)
				if err != nil {
					return nil, err
				}
				m.SetField(f, ev.GetNumber())
				muts = append(muts, m)
			}
		}

		// INTEGER mutations
		switch f.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_INT32,
			descriptorpb.FieldDescriptorProto_TYPE_SINT32,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			ints := []int32{-1, 1, math.MaxInt32, math.MinInt32}
			for _, v := range ints {
				m, err := CloneMessage(seed)
				if err != nil {
					return nil, err
				}
				m.SetField(f, v)
				muts = append(muts, m)
			}
		case descriptorpb.FieldDescriptorProto_TYPE_INT64,
			descriptorpb.FieldDescriptorProto_TYPE_SINT64,
			descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			ints64 := []int64{-1, 1, math.MaxInt64, math.MinInt64}
			for _, v := range ints64 {
				m, err := CloneMessage(seed)
				if err != nil {
					return nil, err
				}
				m.SetField(f, v)
				muts = append(muts, m)
			}
		}
	}

	return muts, nil
}

// RandString returns a random alphanumeric string of length n.
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mrand.Intn(len(letters))]
	}
	return string(b)
}
