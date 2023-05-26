//go:generate go-enum

package code

import "github.com/totvs-cloud/pflagstruct/projscan"

// FieldKind
// ENUM(Native,StdLib,StringMap,TCloudTag,Struct)
type FieldKind string

func KindOf(field *projscan.Field) FieldKind {
	if field.IsTCloudTags() {
		return FieldKindTCloudTag
	}

	if field.Type.IsValid() {
		return FieldKindNative
	}

	if field.Type == "map[string]string" {
		return FieldKindStringMap
	}

	if field.FromStandardLibrary() {
		return FieldKindStdLib
	}

	if field.StructRef != nil && !field.Array {
		return FieldKindStruct
	}

	return ""
}
