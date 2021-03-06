// automatically generated by the FlatBuffers compiler, do not modify

package tracing

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NumericValue struct {
	_tab flatbuffers.Table
}

func GetRootAsNumericValue(buf []byte, offset flatbuffers.UOffsetT) *NumericValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &NumericValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *NumericValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NumericValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *NumericValue) Value() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *NumericValue) MutateValue(n float64) bool {
	return rcv._tab.MutateFloat64Slot(4, n)
}

func NumericValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func NumericValueAddValue(builder *flatbuffers.Builder, value float64) {
	builder.PrependFloat64Slot(0, value, 0.0)
}
func NumericValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
