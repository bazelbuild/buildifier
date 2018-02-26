// Code generated by protoc-gen-go. DO NOT EDIT.
// source: extra_actions_base_proto/extra_actions_base.proto

/*
Package blaze is a generated protocol buffer package.

It is generated from these files:
	extra_actions_base_proto/extra_actions_base.proto

It has these top-level messages:
	ExtraActionSummary
	DetailedExtraActionInfo
	ExtraActionInfo
	EnvironmentVariable
	SpawnInfo
	CppCompileInfo
	CppLinkInfo
	JavaCompileInfo
	PythonInfo
*/
package blaze

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// A list of extra actions and metadata for the print_action command.
type ExtraActionSummary struct {
	Action           []*DetailedExtraActionInfo `protobuf:"bytes,1,rep,name=action" json:"action,omitempty"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (m *ExtraActionSummary) Reset()                    { *m = ExtraActionSummary{} }
func (m *ExtraActionSummary) String() string            { return proto.CompactTextString(m) }
func (*ExtraActionSummary) ProtoMessage()               {}
func (*ExtraActionSummary) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ExtraActionSummary) GetAction() []*DetailedExtraActionInfo {
	if m != nil {
		return m.Action
	}
	return nil
}

// An individual action printed by the print_action command.
type DetailedExtraActionInfo struct {
	// If the given action was included in the output due to a request for a
	// specific file, then this field contains the name of that file so that the
	// caller can correctly associate the extra action with that file.
	//
	// The data in this message is currently not sufficient to run the action on a
	// production machine, because not all necessary input files are identified,
	// especially for C++.
	//
	// There is no easy way to fix this; we could require that all header files
	// are declared and then add all of them here (which would be a huge superset
	// of the files that are actually required), or we could run the include
	// scanner and add those files here.
	TriggeringFile *string `protobuf:"bytes,1,opt,name=triggering_file,json=triggeringFile" json:"triggering_file,omitempty"`
	// The actual action.
	Action           *ExtraActionInfo `protobuf:"bytes,2,req,name=action" json:"action,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *DetailedExtraActionInfo) Reset()                    { *m = DetailedExtraActionInfo{} }
func (m *DetailedExtraActionInfo) String() string            { return proto.CompactTextString(m) }
func (*DetailedExtraActionInfo) ProtoMessage()               {}
func (*DetailedExtraActionInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DetailedExtraActionInfo) GetTriggeringFile() string {
	if m != nil && m.TriggeringFile != nil {
		return *m.TriggeringFile
	}
	return ""
}

func (m *DetailedExtraActionInfo) GetAction() *ExtraActionInfo {
	if m != nil {
		return m.Action
	}
	return nil
}

// Provides information to an extra_action on the original action it is
// shadowing.
type ExtraActionInfo struct {
	// The label of the ActionOwner of the shadowed action.
	Owner *string `protobuf:"bytes,1,opt,name=owner" json:"owner,omitempty"`
	// Only set if the owner is an Aspect.
	// Corresponds to AspectValue.AspectKey.getAspectClass.getName()
	// This field is deprecated as there might now be
	// multiple aspects applied to the same target.
	// This is the aspect name of the last aspect
	// in 'aspects' (8) field.
	AspectName *string `protobuf:"bytes,6,opt,name=aspect_name,json=aspectName" json:"aspect_name,omitempty"`
	// Only set if the owner is an Aspect.
	// Corresponds to AspectValue.AspectKey.getParameters()
	// This field is deprecated as there might now be
	// multiple aspects applied to the same target.
	// These are the aspect parameters of the last aspect
	// in 'aspects' (8) field.
	AspectParameters map[string]*ExtraActionInfo_StringList `protobuf:"bytes,7,rep,name=aspect_parameters,json=aspectParameters" json:"aspect_parameters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// If the owner is an aspect, all aspects applied to the target
	Aspects []*ExtraActionInfo_AspectDescriptor `protobuf:"bytes,8,rep,name=aspects" json:"aspects,omitempty"`
	// An id uniquely describing the shadowed action at the ActionOwner level.
	Id *string `protobuf:"bytes,2,opt,name=id" json:"id,omitempty"`
	// The mnemonic of the shadowed action. Used to distinguish actions with the
	// same ActionType.
	Mnemonic                     *string `protobuf:"bytes,5,opt,name=mnemonic" json:"mnemonic,omitempty"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *ExtraActionInfo) Reset()                    { *m = ExtraActionInfo{} }
func (m *ExtraActionInfo) String() string            { return proto.CompactTextString(m) }
func (*ExtraActionInfo) ProtoMessage()               {}
func (*ExtraActionInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

var extRange_ExtraActionInfo = []proto.ExtensionRange{
	{1000, 536870911},
}

func (*ExtraActionInfo) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ExtraActionInfo
}

func (m *ExtraActionInfo) GetOwner() string {
	if m != nil && m.Owner != nil {
		return *m.Owner
	}
	return ""
}

func (m *ExtraActionInfo) GetAspectName() string {
	if m != nil && m.AspectName != nil {
		return *m.AspectName
	}
	return ""
}

func (m *ExtraActionInfo) GetAspectParameters() map[string]*ExtraActionInfo_StringList {
	if m != nil {
		return m.AspectParameters
	}
	return nil
}

func (m *ExtraActionInfo) GetAspects() []*ExtraActionInfo_AspectDescriptor {
	if m != nil {
		return m.Aspects
	}
	return nil
}

func (m *ExtraActionInfo) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *ExtraActionInfo) GetMnemonic() string {
	if m != nil && m.Mnemonic != nil {
		return *m.Mnemonic
	}
	return ""
}

type ExtraActionInfo_StringList struct {
	Value            []string `protobuf:"bytes,1,rep,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ExtraActionInfo_StringList) Reset()                    { *m = ExtraActionInfo_StringList{} }
func (m *ExtraActionInfo_StringList) String() string            { return proto.CompactTextString(m) }
func (*ExtraActionInfo_StringList) ProtoMessage()               {}
func (*ExtraActionInfo_StringList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 1} }

func (m *ExtraActionInfo_StringList) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

type ExtraActionInfo_AspectDescriptor struct {
	// Corresponds to AspectDescriptor.getName()
	AspectName *string `protobuf:"bytes,1,opt,name=aspect_name,json=aspectName" json:"aspect_name,omitempty"`
	// Corresponds to AspectDescriptor.getParameters()
	AspectParameters map[string]*ExtraActionInfo_AspectDescriptor_StringList `protobuf:"bytes,2,rep,name=aspect_parameters,json=aspectParameters" json:"aspect_parameters,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_unrecognized []byte                                                  `json:"-"`
}

func (m *ExtraActionInfo_AspectDescriptor) Reset()         { *m = ExtraActionInfo_AspectDescriptor{} }
func (m *ExtraActionInfo_AspectDescriptor) String() string { return proto.CompactTextString(m) }
func (*ExtraActionInfo_AspectDescriptor) ProtoMessage()    {}
func (*ExtraActionInfo_AspectDescriptor) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 2}
}

func (m *ExtraActionInfo_AspectDescriptor) GetAspectName() string {
	if m != nil && m.AspectName != nil {
		return *m.AspectName
	}
	return ""
}

func (m *ExtraActionInfo_AspectDescriptor) GetAspectParameters() map[string]*ExtraActionInfo_AspectDescriptor_StringList {
	if m != nil {
		return m.AspectParameters
	}
	return nil
}

type ExtraActionInfo_AspectDescriptor_StringList struct {
	Value            []string `protobuf:"bytes,1,rep,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *ExtraActionInfo_AspectDescriptor_StringList) Reset() {
	*m = ExtraActionInfo_AspectDescriptor_StringList{}
}
func (m *ExtraActionInfo_AspectDescriptor_StringList) String() string {
	return proto.CompactTextString(m)
}
func (*ExtraActionInfo_AspectDescriptor_StringList) ProtoMessage() {}
func (*ExtraActionInfo_AspectDescriptor_StringList) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 2, 1}
}

func (m *ExtraActionInfo_AspectDescriptor_StringList) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

type EnvironmentVariable struct {
	// It is possible that this name is not a valid variable identifier.
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// The value is unescaped and unquoted.
	Value            *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *EnvironmentVariable) Reset()                    { *m = EnvironmentVariable{} }
func (m *EnvironmentVariable) String() string            { return proto.CompactTextString(m) }
func (*EnvironmentVariable) ProtoMessage()               {}
func (*EnvironmentVariable) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *EnvironmentVariable) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *EnvironmentVariable) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

// Provides access to data that is specific to spawn actions.
// Usually provided by actions using the "Spawn" & "Genrule" Mnemonics.
type SpawnInfo struct {
	Argument []string `protobuf:"bytes,1,rep,name=argument" json:"argument,omitempty"`
	// A list of environment variables and their values. No order is enforced.
	Variable         []*EnvironmentVariable `protobuf:"bytes,2,rep,name=variable" json:"variable,omitempty"`
	InputFile        []string               `protobuf:"bytes,4,rep,name=input_file,json=inputFile" json:"input_file,omitempty"`
	OutputFile       []string               `protobuf:"bytes,5,rep,name=output_file,json=outputFile" json:"output_file,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *SpawnInfo) Reset()                    { *m = SpawnInfo{} }
func (m *SpawnInfo) String() string            { return proto.CompactTextString(m) }
func (*SpawnInfo) ProtoMessage()               {}
func (*SpawnInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SpawnInfo) GetArgument() []string {
	if m != nil {
		return m.Argument
	}
	return nil
}

func (m *SpawnInfo) GetVariable() []*EnvironmentVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

func (m *SpawnInfo) GetInputFile() []string {
	if m != nil {
		return m.InputFile
	}
	return nil
}

func (m *SpawnInfo) GetOutputFile() []string {
	if m != nil {
		return m.OutputFile
	}
	return nil
}

var E_SpawnInfo_SpawnInfo = &proto.ExtensionDesc{
	ExtendedType:  (*ExtraActionInfo)(nil),
	ExtensionType: (*SpawnInfo)(nil),
	Field:         1003,
	Name:          "blaze.SpawnInfo.spawn_info",
	Tag:           "bytes,1003,opt,name=spawn_info,json=spawnInfo",
	Filename:      "extra_actions_base_proto/extra_actions_base.proto",
}

// Provides access to data that is specific to C++ compile actions.
// Usually provided by actions using the "CppCompile" Mnemonic.
type CppCompileInfo struct {
	Tool           *string  `protobuf:"bytes,1,opt,name=tool" json:"tool,omitempty"`
	CompilerOption []string `protobuf:"bytes,2,rep,name=compiler_option,json=compilerOption" json:"compiler_option,omitempty"`
	SourceFile     *string  `protobuf:"bytes,3,opt,name=source_file,json=sourceFile" json:"source_file,omitempty"`
	OutputFile     *string  `protobuf:"bytes,4,opt,name=output_file,json=outputFile" json:"output_file,omitempty"`
	// Due to header discovery, this won't include headers unless the build is
	// actually performed. If set, this field will include the value of
	// "source_file" in addition to the headers.
	SourcesAndHeaders []string `protobuf:"bytes,5,rep,name=sources_and_headers,json=sourcesAndHeaders" json:"sources_and_headers,omitempty"`
	// A list of environment variables and their values. No order is enforced.
	Variable         []*EnvironmentVariable `protobuf:"bytes,6,rep,name=variable" json:"variable,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *CppCompileInfo) Reset()                    { *m = CppCompileInfo{} }
func (m *CppCompileInfo) String() string            { return proto.CompactTextString(m) }
func (*CppCompileInfo) ProtoMessage()               {}
func (*CppCompileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *CppCompileInfo) GetTool() string {
	if m != nil && m.Tool != nil {
		return *m.Tool
	}
	return ""
}

func (m *CppCompileInfo) GetCompilerOption() []string {
	if m != nil {
		return m.CompilerOption
	}
	return nil
}

func (m *CppCompileInfo) GetSourceFile() string {
	if m != nil && m.SourceFile != nil {
		return *m.SourceFile
	}
	return ""
}

func (m *CppCompileInfo) GetOutputFile() string {
	if m != nil && m.OutputFile != nil {
		return *m.OutputFile
	}
	return ""
}

func (m *CppCompileInfo) GetSourcesAndHeaders() []string {
	if m != nil {
		return m.SourcesAndHeaders
	}
	return nil
}

func (m *CppCompileInfo) GetVariable() []*EnvironmentVariable {
	if m != nil {
		return m.Variable
	}
	return nil
}

var E_CppCompileInfo_CppCompileInfo = &proto.ExtensionDesc{
	ExtendedType:  (*ExtraActionInfo)(nil),
	ExtensionType: (*CppCompileInfo)(nil),
	Field:         1001,
	Name:          "blaze.CppCompileInfo.cpp_compile_info",
	Tag:           "bytes,1001,opt,name=cpp_compile_info,json=cppCompileInfo",
	Filename:      "extra_actions_base_proto/extra_actions_base.proto",
}

// Provides access to data that is specific to C++ link  actions.
// Usually provided by actions using the "CppLink" Mnemonic.
type CppLinkInfo struct {
	InputFile               []string `protobuf:"bytes,1,rep,name=input_file,json=inputFile" json:"input_file,omitempty"`
	OutputFile              *string  `protobuf:"bytes,2,opt,name=output_file,json=outputFile" json:"output_file,omitempty"`
	InterfaceOutputFile     *string  `protobuf:"bytes,3,opt,name=interface_output_file,json=interfaceOutputFile" json:"interface_output_file,omitempty"`
	LinkTargetType          *string  `protobuf:"bytes,4,opt,name=link_target_type,json=linkTargetType" json:"link_target_type,omitempty"`
	LinkStaticness          *string  `protobuf:"bytes,5,opt,name=link_staticness,json=linkStaticness" json:"link_staticness,omitempty"`
	LinkStamp               []string `protobuf:"bytes,6,rep,name=link_stamp,json=linkStamp" json:"link_stamp,omitempty"`
	BuildInfoHeaderArtifact []string `protobuf:"bytes,7,rep,name=build_info_header_artifact,json=buildInfoHeaderArtifact" json:"build_info_header_artifact,omitempty"`
	// The list of command line options used for running the linking tool.
	LinkOpt          []string `protobuf:"bytes,8,rep,name=link_opt,json=linkOpt" json:"link_opt,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *CppLinkInfo) Reset()                    { *m = CppLinkInfo{} }
func (m *CppLinkInfo) String() string            { return proto.CompactTextString(m) }
func (*CppLinkInfo) ProtoMessage()               {}
func (*CppLinkInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CppLinkInfo) GetInputFile() []string {
	if m != nil {
		return m.InputFile
	}
	return nil
}

func (m *CppLinkInfo) GetOutputFile() string {
	if m != nil && m.OutputFile != nil {
		return *m.OutputFile
	}
	return ""
}

func (m *CppLinkInfo) GetInterfaceOutputFile() string {
	if m != nil && m.InterfaceOutputFile != nil {
		return *m.InterfaceOutputFile
	}
	return ""
}

func (m *CppLinkInfo) GetLinkTargetType() string {
	if m != nil && m.LinkTargetType != nil {
		return *m.LinkTargetType
	}
	return ""
}

func (m *CppLinkInfo) GetLinkStaticness() string {
	if m != nil && m.LinkStaticness != nil {
		return *m.LinkStaticness
	}
	return ""
}

func (m *CppLinkInfo) GetLinkStamp() []string {
	if m != nil {
		return m.LinkStamp
	}
	return nil
}

func (m *CppLinkInfo) GetBuildInfoHeaderArtifact() []string {
	if m != nil {
		return m.BuildInfoHeaderArtifact
	}
	return nil
}

func (m *CppLinkInfo) GetLinkOpt() []string {
	if m != nil {
		return m.LinkOpt
	}
	return nil
}

var E_CppLinkInfo_CppLinkInfo = &proto.ExtensionDesc{
	ExtendedType:  (*ExtraActionInfo)(nil),
	ExtensionType: (*CppLinkInfo)(nil),
	Field:         1002,
	Name:          "blaze.CppLinkInfo.cpp_link_info",
	Tag:           "bytes,1002,opt,name=cpp_link_info,json=cppLinkInfo",
	Filename:      "extra_actions_base_proto/extra_actions_base.proto",
}

// Provides access to data that is specific to java compile actions.
// Usually provided by actions using the "Javac" Mnemonic.
type JavaCompileInfo struct {
	Outputjar        *string  `protobuf:"bytes,1,opt,name=outputjar" json:"outputjar,omitempty"`
	Classpath        []string `protobuf:"bytes,2,rep,name=classpath" json:"classpath,omitempty"`
	Sourcepath       []string `protobuf:"bytes,3,rep,name=sourcepath" json:"sourcepath,omitempty"`
	SourceFile       []string `protobuf:"bytes,4,rep,name=source_file,json=sourceFile" json:"source_file,omitempty"`
	JavacOpt         []string `protobuf:"bytes,5,rep,name=javac_opt,json=javacOpt" json:"javac_opt,omitempty"`
	Processor        []string `protobuf:"bytes,6,rep,name=processor" json:"processor,omitempty"`
	Processorpath    []string `protobuf:"bytes,7,rep,name=processorpath" json:"processorpath,omitempty"`
	Bootclasspath    []string `protobuf:"bytes,8,rep,name=bootclasspath" json:"bootclasspath,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *JavaCompileInfo) Reset()                    { *m = JavaCompileInfo{} }
func (m *JavaCompileInfo) String() string            { return proto.CompactTextString(m) }
func (*JavaCompileInfo) ProtoMessage()               {}
func (*JavaCompileInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *JavaCompileInfo) GetOutputjar() string {
	if m != nil && m.Outputjar != nil {
		return *m.Outputjar
	}
	return ""
}

func (m *JavaCompileInfo) GetClasspath() []string {
	if m != nil {
		return m.Classpath
	}
	return nil
}

func (m *JavaCompileInfo) GetSourcepath() []string {
	if m != nil {
		return m.Sourcepath
	}
	return nil
}

func (m *JavaCompileInfo) GetSourceFile() []string {
	if m != nil {
		return m.SourceFile
	}
	return nil
}

func (m *JavaCompileInfo) GetJavacOpt() []string {
	if m != nil {
		return m.JavacOpt
	}
	return nil
}

func (m *JavaCompileInfo) GetProcessor() []string {
	if m != nil {
		return m.Processor
	}
	return nil
}

func (m *JavaCompileInfo) GetProcessorpath() []string {
	if m != nil {
		return m.Processorpath
	}
	return nil
}

func (m *JavaCompileInfo) GetBootclasspath() []string {
	if m != nil {
		return m.Bootclasspath
	}
	return nil
}

var E_JavaCompileInfo_JavaCompileInfo = &proto.ExtensionDesc{
	ExtendedType:  (*ExtraActionInfo)(nil),
	ExtensionType: (*JavaCompileInfo)(nil),
	Field:         1000,
	Name:          "blaze.JavaCompileInfo.java_compile_info",
	Tag:           "bytes,1000,opt,name=java_compile_info,json=javaCompileInfo",
	Filename:      "extra_actions_base_proto/extra_actions_base.proto",
}

// Provides access to data that is specific to python rules.
// Usually provided by actions using the "Python" Mnemonic.
type PythonInfo struct {
	SourceFile       []string `protobuf:"bytes,1,rep,name=source_file,json=sourceFile" json:"source_file,omitempty"`
	DepFile          []string `protobuf:"bytes,2,rep,name=dep_file,json=depFile" json:"dep_file,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *PythonInfo) Reset()                    { *m = PythonInfo{} }
func (m *PythonInfo) String() string            { return proto.CompactTextString(m) }
func (*PythonInfo) ProtoMessage()               {}
func (*PythonInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *PythonInfo) GetSourceFile() []string {
	if m != nil {
		return m.SourceFile
	}
	return nil
}

func (m *PythonInfo) GetDepFile() []string {
	if m != nil {
		return m.DepFile
	}
	return nil
}

var E_PythonInfo_PythonInfo = &proto.ExtensionDesc{
	ExtendedType:  (*ExtraActionInfo)(nil),
	ExtensionType: (*PythonInfo)(nil),
	Field:         1005,
	Name:          "blaze.PythonInfo.python_info",
	Tag:           "bytes,1005,opt,name=python_info,json=pythonInfo",
	Filename:      "extra_actions_base_proto/extra_actions_base.proto",
}

func init() {
	proto.RegisterType((*ExtraActionSummary)(nil), "blaze.ExtraActionSummary")
	proto.RegisterType((*DetailedExtraActionInfo)(nil), "blaze.DetailedExtraActionInfo")
	proto.RegisterType((*ExtraActionInfo)(nil), "blaze.ExtraActionInfo")
	proto.RegisterType((*ExtraActionInfo_StringList)(nil), "blaze.ExtraActionInfo.StringList")
	proto.RegisterType((*ExtraActionInfo_AspectDescriptor)(nil), "blaze.ExtraActionInfo.AspectDescriptor")
	proto.RegisterType((*ExtraActionInfo_AspectDescriptor_StringList)(nil), "blaze.ExtraActionInfo.AspectDescriptor.StringList")
	proto.RegisterType((*EnvironmentVariable)(nil), "blaze.EnvironmentVariable")
	proto.RegisterType((*SpawnInfo)(nil), "blaze.SpawnInfo")
	proto.RegisterType((*CppCompileInfo)(nil), "blaze.CppCompileInfo")
	proto.RegisterType((*CppLinkInfo)(nil), "blaze.CppLinkInfo")
	proto.RegisterType((*JavaCompileInfo)(nil), "blaze.JavaCompileInfo")
	proto.RegisterType((*PythonInfo)(nil), "blaze.PythonInfo")
	proto.RegisterExtension(E_SpawnInfo_SpawnInfo)
	proto.RegisterExtension(E_CppCompileInfo_CppCompileInfo)
	proto.RegisterExtension(E_CppLinkInfo_CppLinkInfo)
	proto.RegisterExtension(E_JavaCompileInfo_JavaCompileInfo)
	proto.RegisterExtension(E_PythonInfo_PythonInfo)
}

func init() { proto.RegisterFile("extra_actions_base_proto/extra_actions_base.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1038 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x56, 0xdd, 0x6e, 0x23, 0x35,
	0x14, 0x56, 0xa6, 0x4d, 0x9b, 0x39, 0xd1, 0xa6, 0xa9, 0x4b, 0xd9, 0x6c, 0x80, 0xdd, 0x12, 0x10,
	0x5b, 0x01, 0x9a, 0x15, 0xb9, 0x58, 0x10, 0x08, 0xa1, 0x6e, 0xb7, 0xa8, 0x40, 0x45, 0xab, 0xe9,
	0x0a, 0x21, 0x21, 0x34, 0x72, 0x26, 0x6e, 0xea, 0x76, 0x66, 0x6c, 0x79, 0x9c, 0x94, 0x70, 0xd5,
	0x97, 0xe0, 0x8e, 0x47, 0xe1, 0x25, 0x78, 0x8b, 0x5d, 0x10, 0x37, 0x88, 0x07, 0x40, 0x3e, 0xf6,
	0xcc, 0x34, 0x69, 0xda, 0xed, 0x9d, 0xfd, 0x9d, 0x6f, 0xce, 0xcf, 0xf7, 0xd9, 0x4e, 0xe0, 0x13,
	0xf6, 0x8b, 0x56, 0x34, 0xa2, 0xb1, 0xe6, 0x22, 0xcb, 0xa3, 0x01, 0xcd, 0x59, 0x24, 0x95, 0xd0,
	0xe2, 0xc9, 0xf5, 0x40, 0x80, 0x01, 0x52, 0x1f, 0x24, 0xf4, 0x57, 0xd6, 0x3b, 0x00, 0xb2, 0x67,
	0x28, 0x3b, 0xc8, 0x38, 0x1e, 0xa7, 0x29, 0x55, 0x53, 0xf2, 0x14, 0x56, 0xec, 0x27, 0x9d, 0xda,
	0xd6, 0xd2, 0x76, 0xb3, 0xff, 0x30, 0x40, 0x76, 0xf0, 0x9c, 0x69, 0xca, 0x13, 0x36, 0xbc, 0xf2,
	0xc9, 0x37, 0xd9, 0x89, 0x08, 0x1d, 0xbb, 0xa7, 0xe0, 0xfe, 0x0d, 0x14, 0xf2, 0x18, 0xd6, 0xb4,
	0xe2, 0xa3, 0x11, 0x53, 0x3c, 0x1b, 0x45, 0x27, 0x3c, 0x61, 0x9d, 0xda, 0x56, 0x6d, 0xdb, 0x0f,
	0x5b, 0x15, 0xfc, 0x35, 0x4f, 0x18, 0x09, 0xca, 0xda, 0xde, 0x96, 0xb7, 0xdd, 0xec, 0xbf, 0xe9,
	0x6a, 0xdf, 0x54, 0xf3, 0xbf, 0x3a, 0xac, 0xcd, 0x17, 0x7b, 0x03, 0xea, 0xe2, 0x22, 0x63, 0xca,
	0x95, 0xb0, 0x1b, 0xf2, 0x1e, 0x34, 0x69, 0x2e, 0x59, 0xac, 0xa3, 0x8c, 0xa6, 0xac, 0xb3, 0x62,
	0x62, 0xcf, 0xbc, 0x4e, 0x2d, 0x04, 0x0b, 0x7f, 0x4f, 0x53, 0x46, 0x7e, 0x86, 0x75, 0x47, 0x92,
	0x54, 0xd1, 0x94, 0x69, 0xa6, 0xf2, 0xce, 0x2a, 0xaa, 0xf0, 0xf1, 0xe2, 0x4e, 0x82, 0x1d, 0xe4,
	0x1f, 0x95, 0xf4, 0xbd, 0x4c, 0xab, 0x29, 0x26, 0x6e, 0xd3, 0xb9, 0x10, 0xd9, 0x81, 0x55, 0x8b,
	0xe5, 0x9d, 0x06, 0x26, 0x7d, 0x7c, 0x6b, 0xd2, 0xe7, 0x2c, 0x8f, 0x15, 0x97, 0x5a, 0xa8, 0xb0,
	0xf8, 0x8e, 0xb4, 0xc0, 0xe3, 0xc3, 0x8e, 0x87, 0x93, 0x79, 0x7c, 0x48, 0xba, 0xd0, 0x48, 0x33,
	0x96, 0x8a, 0x8c, 0xc7, 0x9d, 0x3a, 0xa2, 0xe5, 0xbe, 0x7b, 0x02, 0x9b, 0x0b, 0xbb, 0x23, 0x6d,
	0x58, 0x3a, 0x67, 0x53, 0xa7, 0x8f, 0x59, 0x92, 0x4f, 0xa1, 0x3e, 0xa1, 0xc9, 0x98, 0x61, 0xe6,
	0x66, 0xff, 0xdd, 0x1b, 0xfa, 0x3a, 0xd6, 0xc6, 0xa9, 0x03, 0x9e, 0xeb, 0xd0, 0xf2, 0x3f, 0xf7,
	0x3e, 0xab, 0x75, 0x3f, 0x00, 0xa8, 0x02, 0x46, 0x7e, 0x9b, 0xca, 0x9c, 0x1e, 0xbf, 0xe4, 0x75,
	0x6a, 0xdd, 0x3f, 0x3c, 0x68, 0xcf, 0x4f, 0x46, 0x1e, 0xcd, 0xfa, 0x62, 0x7b, 0xba, 0xea, 0xc9,
	0xd9, 0x22, 0x4f, 0x3c, 0x94, 0xef, 0xcb, 0x3b, 0xca, 0xb7, 0xd8, 0xa4, 0xeb, 0x06, 0x75, 0x2f,
	0xee, 0xae, 0xd8, 0xfe, 0xac, 0x62, 0xfd, 0xbb, 0xb6, 0xb2, 0x58, 0xc2, 0xde, 0xeb, 0x25, 0xfc,
	0xd0, 0x6f, 0xbc, 0x5c, 0x6d, 0x5f, 0x5e, 0x5e, 0x5e, 0x7a, 0xbd, 0xaf, 0x60, 0x63, 0x2f, 0x9b,
	0x70, 0x25, 0xb2, 0x94, 0x65, 0xfa, 0x07, 0xaa, 0x38, 0x1d, 0x24, 0x8c, 0x10, 0x58, 0x76, 0x22,
	0x7a, 0xdb, 0x7e, 0x88, 0xeb, 0x2a, 0x97, 0x87, 0xa0, 0xdd, 0xf4, 0x5e, 0xd5, 0xc0, 0x3f, 0x96,
	0xf4, 0xc2, 0xde, 0x98, 0x2e, 0x34, 0xa8, 0x1a, 0x8d, 0x4d, 0x2e, 0x57, 0xb2, 0xdc, 0x93, 0xa7,
	0xd0, 0x98, 0xb8, 0xfc, 0x4e, 0xf5, 0x6e, 0x31, 0xea, 0xf5, 0x0e, 0xc2, 0x92, 0x4b, 0xde, 0x01,
	0xe0, 0x99, 0x1c, 0x6b, 0x7b, 0xdb, 0x97, 0x31, 0xab, 0x8f, 0x08, 0x5e, 0xf4, 0x47, 0xd0, 0x14,
	0x63, 0x5d, 0xc6, 0xeb, 0x18, 0x07, 0x0b, 0x19, 0x42, 0x7f, 0x1f, 0x20, 0x37, 0x0d, 0x46, 0xdc,
	0x74, 0x78, 0xc3, 0x3b, 0xd0, 0xf9, 0x7b, 0x15, 0xd5, 0x6f, 0xbb, 0x70, 0x39, 0x52, 0xe8, 0xe7,
	0xc5, 0xb2, 0xf7, 0xa7, 0x07, 0xad, 0x5d, 0x29, 0x77, 0x45, 0x2a, 0x79, 0xc2, 0x70, 0x60, 0x02,
	0xcb, 0x5a, 0x88, 0xc4, 0xf9, 0x89, 0x6b, 0xf3, 0x46, 0xc5, 0x96, 0xa2, 0x22, 0x21, 0xdd, 0x1b,
	0x64, 0xba, 0x6a, 0x15, 0xf0, 0x21, 0xa2, 0xa6, 0xf5, 0x5c, 0x8c, 0x55, 0xcc, 0x6c, 0xeb, 0x4b,
	0xf6, 0xc4, 0x5a, 0x68, 0xd1, 0x6c, 0xcb, 0x96, 0x50, 0xcd, 0x46, 0x02, 0xd8, 0xb0, 0xf4, 0x3c,
	0xa2, 0xd9, 0x30, 0x3a, 0x65, 0x74, 0x68, 0x0e, 0xb5, 0x15, 0x61, 0xdd, 0x85, 0x76, 0xb2, 0xe1,
	0xbe, 0x0d, 0xcc, 0x78, 0xb0, 0x72, 0x77, 0x0f, 0xfa, 0x3f, 0x42, 0x3b, 0x96, 0x32, 0x72, 0xfd,
	0xdf, 0xae, 0xe4, 0x2b, 0xab, 0xe4, 0xa6, 0x0b, 0xcf, 0x0a, 0x16, 0xb6, 0xe2, 0x99, 0x7d, 0xef,
	0xf7, 0x25, 0x68, 0xee, 0x4a, 0x79, 0xc0, 0xb3, 0x73, 0x14, 0x74, 0xd6, 0xed, 0xda, 0x6b, 0xdc,
	0xf6, 0xae, 0x29, 0xd2, 0x87, 0x4d, 0x9e, 0x69, 0xa6, 0x4e, 0x68, 0xcc, 0xa2, 0xab, 0x54, 0xab,
	0xee, 0x46, 0x19, 0x3c, 0xac, 0xbe, 0xd9, 0x86, 0x76, 0xc2, 0xb3, 0xf3, 0x48, 0x53, 0x35, 0x62,
	0x3a, 0xd2, 0x53, 0x59, 0x68, 0xdd, 0x32, 0xf8, 0x0b, 0x84, 0x5f, 0x4c, 0x25, 0x33, 0xd6, 0x22,
	0x33, 0xd7, 0x54, 0xf3, 0x38, 0x63, 0x79, 0xee, 0xde, 0x4a, 0x24, 0x1e, 0x97, 0xa8, 0x19, 0xa3,
	0x20, 0xa6, 0x12, 0xa5, 0xf6, 0x43, 0xdf, 0x71, 0x52, 0x49, 0xbe, 0x80, 0xee, 0x60, 0xcc, 0x93,
	0x21, 0x2a, 0xe9, 0x6c, 0x8b, 0xa8, 0xd2, 0xfc, 0x84, 0xc6, 0x1a, 0x7f, 0x27, 0xfc, 0xf0, 0x3e,
	0x32, 0x8c, 0x28, 0xd6, 0xbd, 0x1d, 0x17, 0x26, 0x0f, 0xa0, 0x81, 0xb9, 0x85, 0xd4, 0xf8, 0xfa,
	0xfb, 0xe1, 0xaa, 0xd9, 0x1f, 0x4a, 0xdd, 0x3f, 0x84, 0x7b, 0xc6, 0x27, 0x0c, 0xdf, 0x6a, 0xd2,
	0x5f, 0xd6, 0x24, 0x52, 0x99, 0x54, 0x38, 0x10, 0x36, 0xe3, 0x6a, 0xd3, 0xfb, 0xd7, 0x83, 0xb5,
	0x6f, 0xe9, 0x84, 0x5e, 0x3d, 0xf3, 0x6f, 0x83, 0x6f, 0x85, 0x3d, 0xa3, 0xc5, 0x4f, 0x63, 0x05,
	0x98, 0x68, 0x9c, 0xd0, 0x3c, 0x97, 0x54, 0x9f, 0xba, 0x73, 0x5f, 0x01, 0xe4, 0x21, 0xb8, 0xf3,
	0x8d, 0xe1, 0x25, 0x7b, 0x59, 0x2b, 0x64, 0xfe, 0x4a, 0x2c, 0x5f, 0x25, 0xa0, 0x57, 0x6f, 0x81,
	0x7f, 0x46, 0x27, 0x34, 0xc6, 0xe9, 0xed, 0x39, 0x6f, 0x20, 0x70, 0x28, 0xb5, 0xa9, 0x2d, 0x95,
	0x88, 0x59, 0x9e, 0x0b, 0x55, 0x88, 0x5e, 0x02, 0xe4, 0x7d, 0xb8, 0x57, 0x6e, 0xb0, 0xbc, 0xd5,
	0x79, 0x16, 0x34, 0xac, 0x81, 0x10, 0xba, 0x9a, 0xc1, 0x4a, 0x3c, 0x0b, 0xf6, 0x7f, 0x82, 0x75,
	0x53, 0xf5, 0x6e, 0x37, 0xe2, 0xa5, 0x15, 0xbb, 0x08, 0xcf, 0xe9, 0x19, 0xae, 0x9d, 0xcd, 0x02,
	0xbd, 0xdf, 0x6a, 0x00, 0x47, 0x53, 0x7d, 0xea, 0xfe, 0x86, 0xcc, 0x69, 0x52, 0xbb, 0xa6, 0xc9,
	0x03, 0x68, 0x0c, 0x99, 0x2c, 0x6e, 0x04, 0x1e, 0x88, 0x21, 0x93, 0xf8, 0xf8, 0x7d, 0x07, 0x4d,
	0x89, 0x99, 0x6e, 0xef, 0xf0, 0x1f, 0xdb, 0xe1, 0xba, 0x0b, 0x57, 0xc5, 0x43, 0x90, 0xe5, 0xfa,
	0xd9, 0x13, 0xf8, 0x28, 0x16, 0x69, 0x30, 0x12, 0x62, 0x94, 0xb0, 0x60, 0xc8, 0x26, 0xe6, 0xb9,
	0xcb, 0x03, 0x3c, 0xa7, 0x41, 0xc2, 0x07, 0x81, 0xfb, 0x83, 0x18, 0xe0, 0xdf, 0xc5, 0xa3, 0xda,
	0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x81, 0x96, 0x76, 0xca, 0x51, 0x0a, 0x00, 0x00,
}
