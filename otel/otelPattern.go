package otel

import (
	"context"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type IOtelPattern interface {
	ParentSpan(name string)
	ChildSpan(name string)
	ChildSpanDown(name string)
	HandlerSuccessParent(name string)
	HandlerSuccess(name string)
	HandlerFail(status string, err string)
	SetAttributes(att ...attribute.KeyValue)            // add tag
	SetAttributesByStruct(att interface{})              // add tag
	AddEvent(name string, options ...trace.EventOption) // add log
	AddEventByStruct(name string, options interface{})  // add log
	SpanEnd()
	SetAttributesParent(att ...attribute.KeyValue)            // add tag
	SetAttributesByStructParent(att interface{})              // add tag
	AddEventParent(name string, options ...trace.EventOption) // add log
	AddEventByStructParent(name string, options interface{})  // add log
	SpanEndParent()
	Att(key string, val string) trace.SpanStartEventOption
	St(key string, val string) attribute.KeyValue
	OtelEnd()
}

type OtelPattern struct {
	Ctx      context.Context
	ChildCtx context.Context
	Tracer   trace.Tracer
	Parent   trace.Span
	Span     trace.Span
	Otel     trace.Span
}

func CreateOtel(ctx context.Context, name string, rid string) IOtelPattern {
	tracer := otel.Tracer(name)
	context, span := tracer.Start(ctx, name)
	span.SetStatus(codes.Ok, name)
	span.SetAttributes(attribute.String("route", name), attribute.String("rid", rid))
	return &OtelPattern{
		Ctx:    context,
		Tracer: tracer,
		Otel:   span,
	}
}

func (ot *OtelPattern) ParentSpan(name string) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	context, span := ot.Tracer.Start(ot.Ctx, name)
	ot.ChildCtx = context
	ot.Parent = span
}

func (ot *OtelPattern) ChildSpan(name string) {
	_, span := ot.Tracer.Start(ot.ChildCtx, name)
	// ot.ChildCtx = context
	ot.Span = span
}

func (ot *OtelPattern) ChildSpanDown(name string) {
	context, span := ot.Tracer.Start(ot.ChildCtx, name)
	ot.ChildCtx = context
	ot.Span = span
}

func (ot *OtelPattern) HandlerSuccess(name string) {
	ot.Span.SetStatus(codes.Ok, name)
	ot.Span.End()
}

func (ot *OtelPattern) HandlerSuccessParent(name string) {
	ot.Parent.SetStatus(codes.Ok, name)
	ot.Parent.End()
}

func (ot *OtelPattern) HandlerFail(status string, err string) {
	if ot.Span != nil {
		ot.Span.RecordError(fmt.Errorf(err))
		ot.Span.SetStatus(codes.Error, status)
		ot.Span.End()
	}
	if ot.Parent != nil {
		ot.Parent.RecordError(fmt.Errorf(err))
		ot.Parent.SetStatus(codes.Error, status)
		ot.Parent.End()
	}
	if ot.Otel != nil {
		ot.Otel.End()
	}
}

func (ot *OtelPattern) SetAttributes(att ...attribute.KeyValue) {
	ot.Span.SetAttributes(att...)
}

func (ot *OtelPattern) SetAttributesByStruct(att interface{}) {
	fields := convertStructToAttributes(att)
	ot.Span.SetAttributes(fields...)
}

func (ot *OtelPattern) AddEvent(name string, options ...trace.EventOption) {
	ot.Span.AddEvent(name, options...)
}

func (ot *OtelPattern) AddEventByStruct(name string, options interface{}) {
	fields := convertStructToAttributes(options)
	ot.Span.AddEvent(name, trace.WithAttributes(fields...))
}

func (ot *OtelPattern) SpanEnd() {
	ot.Span.End()
}

func (ot *OtelPattern) SetAttributesParent(att ...attribute.KeyValue) {
	ot.Parent.SetAttributes(att...)
}

func (ot *OtelPattern) SetAttributesByStructParent(att interface{}) {
	fields := convertStructToAttributes(att)
	ot.Parent.SetAttributes(fields...)
}

func (ot *OtelPattern) AddEventParent(name string, options ...trace.EventOption) {
	ot.Parent.AddEvent(name, options...)
}

func (ot *OtelPattern) AddEventByStructParent(name string, options interface{}) {
	fields := convertStructToAttributes(options)
	ot.Parent.AddEvent(name, trace.WithAttributes(fields...))
}

func (ot *OtelPattern) SpanEndParent() {
	ot.Parent.End()
}

func (ot *OtelPattern) Att(key string, val string) trace.SpanStartEventOption {
	return trace.WithAttributes(attribute.String(key, val))
}

func (ot *OtelPattern) St(key string, val string) attribute.KeyValue {
	return attribute.String(key, val)
}

func (ot *OtelPattern) OtelEnd() {
	ot.Otel.End()
}

func convertStructToAttributes(s interface{}) []attribute.KeyValue {
	var attributes []attribute.KeyValue
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i).Interface()
			key := attribute.Key(field.Name)
			attributes = appendAttribute(attributes, key, value)
		}
	}

	return attributes
}

func appendAttribute(attributes []attribute.KeyValue, key attribute.Key, value interface{}) []attribute.KeyValue {
	switch v := value.(type) {
	case string:
		attributes = append(attributes, key.String(v))
	case int:
		attributes = append(attributes, key.Int(v))
	case bool:
		attributes = append(attributes, key.Bool(v))
	// และเพิ่ม case สำหรับชนิดข้อมูลอื่น ๆ ตามที่คุณต้องการ
	default:
		// ชนิดข้อมูลที่ไม่รองรับ
		attributes = append(attributes, key.String(fmt.Sprintf("%v", v)))
	}
	return attributes
}
