package example_types

type Int struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64
}

type IntArray struct {
	Int   []int
	Int8  []int8
	Int16 []int16
	Int32 []int32
	Int64 []int64
}

type IntArrayArray struct {
	Int   [][]int
	Int8  [][][]int8
	Int16 [][][][]int16
	Int32 [][][][][]int32
	Int64 [][][][][][]int64
}

type Uint struct {
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
}

type UintArray struct {
	Uint   []uint
	Uint8  []uint8
	Uint16 []uint16
	Uint32 []uint32
	Uint64 []uint64
}

type UintArrayArray struct {
	Uint   [][]uint
	Uint8  [][][]uint8
	Uint16 [][][][]uint16
	Uint32 [][][][][]uint32
	Uint64 [][][][][][]uint64
}

type String struct {
	String string
}

type StringArray struct {
	String []string
}

type StringArrayArray struct {
	String2 [][]string
	String3 [][][]string
	String4 [][][][]string
}

type Float struct {
	Float32 float32
	Float64 float64
}

type FloatArray struct {
	Float32 []float32
	Float64 []float64
}

type FloatArrayArray struct {
	Float322 [][]float32
	Float323 [][][]float32
	Float642 [][]float64
	Float643 [][][]float64
}

type Bool struct {
	Bool bool
}

type BoolArray struct {
	Bool []bool
}

type BoolArrayArray struct {
	Bool2 [][]bool
	Bool3 [][][]bool
	Bool4 [][][]bool
}

type Full struct {
	Int
	IntArray
	IntArrayArray
	Uint
	UintArray
	UintArrayArray
	String
	StringArray
	StringArrayArray
	Float
	FloatArray
	FloatArrayArray
	Bool
	BoolArray
	BoolArrayArray
}
