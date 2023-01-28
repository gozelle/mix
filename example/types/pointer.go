package example_types

type IntPointer struct {
	Int   *int
	Int8  *int8
	Int16 *int16
	Int32 *int32
	Int64 *int64
}

type IntArrayPointer struct {
	Int   []*int
	Int8  []*int8
	Int16 []*int16
	Int32 []*int32
	Int64 []*int64
}

type IntArrayArrayPointer struct {
	Int   [][]*int
	Int8  [][][]*int8
	Int16 [][][][]*int16
	Int32 [][][][][]*int32
	Int64 [][][][][][]*int64
}

type UintPointer struct {
	Uint   *uint
	Uint8  *uint8
	Uint16 *uint16
	Uint32 *uint32
	Uint64 *uint64
}

type UintArrayPointer struct {
	Uint   []*uint
	Uint8  []*uint8
	Uint16 []*uint16
	Uint32 []*uint32
	Uint64 []*uint64
}

type UintArrayArrayPointer struct {
	Uint   [][]*uint
	Uint8  [][][]*uint8
	Uint16 [][][][]*uint16
	Uint32 [][][][][]*uint32
	Uint64 [][][][][][]*uint64
}

type StringPointer struct {
	String *string
}

type StringArrayPointer struct {
	String []*string
}

type StringArrayArrayPointer struct {
	String2 [][]*string
	String3 [][][]*string
	String4 [][][][]*string
}

type FloatPointer struct {
	Float32 *float32
	Float64 *float64
}

type FloatArrayPointer struct {
	Float32 []*float32
	Float64 []*float64
}

type FloatArrayArrayPointer struct {
	Float322 [][]*float32
	Float323 [][][]*float32
	Float642 [][]*float64
	Float643 [][][]*float64
}

type BoolPointer struct {
	Bool *bool
}

type BoolArrayPointer struct {
	Bool []*bool
}

type BoolArrayArrayPointer struct {
	Bool2 [][]*bool
	Bool3 [][][]*bool
	Bool4 [][][]*bool
}

type FullPointer struct {
	*Int
	*IntArray
	*IntArrayArray
	*Uint
	*UintArray
	*UintArrayArray
	*String
	*StringArray
	*StringArrayArray
	*Float
	*FloatArray
	*FloatArrayArray
	*Bool
	*BoolArray
	*BoolArrayArray
}
