// automatically generated by the FlatBuffers compiler, do not modify

package tracing

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ChildOf struct {
	_tab flatbuffers.Table
}

func GetRootAsChildOf(buf []byte, offset flatbuffers.UOffsetT) *ChildOf {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ChildOf{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *ChildOf) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ChildOf) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *ChildOf) Parent(obj *SpanContext) *SpanContext {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(SpanContext)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func ChildOfStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ChildOfAddParent(builder *flatbuffers.Builder, parent flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(parent), 0)
}
func ChildOfEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
