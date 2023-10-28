// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufimage

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func stripSourceOnlyOptions[M proto.Message](options M) M {
	optionsRef := options.ProtoReflect()
	// See if there are any options to strip.
	var found bool
	optionsRef.Range(func(field protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		if field.Options().(*descriptorpb.FieldOptions).GetRetention() == descriptorpb.FieldOptions_RETENTION_SOURCE {
			found = true
			return false
		}
		return true
	})
	if !found {
		return options
	}

	// There is at least one. So we need to make a copy that does not have those options.
	newOptions := optionsRef.New()
	optionsRef.Range(func(field protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		if field.Options().(*descriptorpb.FieldOptions).GetRetention() != descriptorpb.FieldOptions_RETENTION_SOURCE {
			newOptions.Set(field, val)
		}
		return true
	})
	return newOptions.Interface().(M)
}

func stripSourceOnlyOptionsFromFile(file *descriptorpb.FileDescriptorProto) *descriptorpb.FileDescriptorProto {
	var dirty bool
	newOpts := stripSourceOnlyOptions(file.Options)
	if newOpts != file.Options {
		dirty = true
	}
	newMsgs, changed := updateAll(file.MessageType, stripSourceOnlyOptionsFromMessage)
	if changed {
		dirty = true
	}
	newEnums, changed := updateAll(file.EnumType, stripSourceOnlyOptionsFromEnum)
	if changed {
		dirty = true
	}
	newExts, changed := updateAll(file.Extension, stripSourceOnlyOptionsFromField)
	if changed {
		dirty = true
	}
	newSvcs, changed := updateAll(file.Service, stripSourceOnlyOptionsFromService)
	if changed {
		dirty = true
	}

	if !dirty {
		return file
	}

	newFile := shallowCopy(file)
	newFile.Options = newOpts
	newFile.MessageType = newMsgs
	newFile.EnumType = newEnums
	newFile.Extension = newExts
	newFile.Service = newSvcs
	return newFile
}

func stripSourceOnlyOptionsFromMessage(msg *descriptorpb.DescriptorProto) *descriptorpb.DescriptorProto {
	var dirty bool
	newOpts := stripSourceOnlyOptions(msg.Options)
	if newOpts != msg.Options {
		dirty = true
	}
	newFields, changed := updateAll(msg.Field, stripSourceOnlyOptionsFromField)
	if changed {
		dirty = true
	}
	newOneofs, changed := updateAll(msg.OneofDecl, stripSourceOnlyOptionsFromOneof)
	if changed {
		dirty = true
	}
	newExtRanges, changed := updateAll(msg.ExtensionRange, stripSourceOnlyOptionsFromExtensionRange)
	if changed {
		dirty = true
	}
	newMsgs, changed := updateAll(msg.NestedType, stripSourceOnlyOptionsFromMessage)
	if changed {
		dirty = true
	}
	newEnums, changed := updateAll(msg.EnumType, stripSourceOnlyOptionsFromEnum)
	if changed {
		dirty = true
	}
	newExts, changed := updateAll(msg.Extension, stripSourceOnlyOptionsFromField)
	if changed {
		dirty = true
	}

	if !dirty {
		return msg
	}

	newMsg := shallowCopy(msg)
	newMsg.Options = newOpts
	newMsg.Field = newFields
	newMsg.OneofDecl = newOneofs
	newMsg.ExtensionRange = newExtRanges
	newMsg.NestedType = newMsgs
	newMsg.EnumType = newEnums
	newMsg.Extension = newExts
	return newMsg
}

func stripSourceOnlyOptionsFromField(field *descriptorpb.FieldDescriptorProto) *descriptorpb.FieldDescriptorProto {
	newOpts := stripSourceOnlyOptions(field.Options)
	if newOpts == field.Options {
		return field
	}
	newField := shallowCopy(field)
	newField.Options = newOpts
	return newField
}

func stripSourceOnlyOptionsFromOneof(oneof *descriptorpb.OneofDescriptorProto) *descriptorpb.OneofDescriptorProto {
	newOpts := stripSourceOnlyOptions(oneof.Options)
	if newOpts == oneof.Options {
		return oneof
	}
	newOneof := shallowCopy(oneof)
	newOneof.Options = newOpts
	return newOneof
}

func stripSourceOnlyOptionsFromExtensionRange(extRange *descriptorpb.DescriptorProto_ExtensionRange) *descriptorpb.DescriptorProto_ExtensionRange {
	newOpts := stripSourceOnlyOptions(extRange.Options)
	if newOpts == extRange.Options {
		return extRange
	}
	newExtRange := shallowCopy(extRange)
	newExtRange.Options = newOpts
	return newExtRange
}

func stripSourceOnlyOptionsFromEnum(enum *descriptorpb.EnumDescriptorProto) *descriptorpb.EnumDescriptorProto {
	var dirty bool
	newOpts := stripSourceOnlyOptions(enum.Options)
	if newOpts != enum.Options {
		dirty = true
	}
	newVals, changed := updateAll(enum.Value, stripSourceOnlyOptionsFromEnumValue)
	if changed {
		dirty = true
	}

	if !dirty {
		return enum
	}

	newEnum := shallowCopy(enum)
	newEnum.Options = newOpts
	newEnum.Value = newVals
	return newEnum
}

func stripSourceOnlyOptionsFromEnumValue(enumVal *descriptorpb.EnumValueDescriptorProto) *descriptorpb.EnumValueDescriptorProto {
	newOpts := stripSourceOnlyOptions(enumVal.Options)
	if newOpts == enumVal.Options {
		return enumVal
	}
	newEnumVal := shallowCopy(enumVal)
	newEnumVal.Options = newOpts
	return newEnumVal
}

func stripSourceOnlyOptionsFromService(svc *descriptorpb.ServiceDescriptorProto) *descriptorpb.ServiceDescriptorProto {
	var dirty bool
	newOpts := stripSourceOnlyOptions(svc.Options)
	if newOpts != svc.Options {
		dirty = true
	}
	newMethods, changed := updateAll(svc.Method, stripSourceOnlyOptionsFromMethod)
	if changed {
		dirty = true
	}

	if !dirty {
		return svc
	}

	newSvc := shallowCopy(svc)
	newSvc.Options = newOpts
	newSvc.Method = newMethods
	return newSvc
}

func stripSourceOnlyOptionsFromMethod(method *descriptorpb.MethodDescriptorProto) *descriptorpb.MethodDescriptorProto {
	newOpts := stripSourceOnlyOptions(method.Options)
	if newOpts == method.Options {
		return method
	}
	newMethod := shallowCopy(method)
	newMethod.Options = newOpts
	return newMethod
}

func shallowCopy[M proto.Message](msg M) M {
	msgRef := msg.ProtoReflect()
	other := msgRef.New()
	msgRef.Range(func(field protoreflect.FieldDescriptor, val protoreflect.Value) bool {
		other.Set(field, val)
		return true
	})
	return other.Interface().(M)
}

func updateAll[T comparable](slice []T, updateFunc func(T) T) ([]T, bool) {
	var updated []T // initialized lazily, only when/if a copy is needed
	for i, item := range slice {
		newItem := updateFunc(item)
		if updated != nil {
			updated[i] = newItem
		} else if newItem != item {
			updated = make([]T, len(slice))
			copy(updated[:i], slice)
			updated[i] = newItem
		}
	}
	if updated != nil {
		return updated, true
	}
	return slice, false
}
