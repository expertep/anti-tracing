package otel

import (
	"context"
	"reflect"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Attributes represent additional key-value descriptors that can be bound
func CreateSpan(level string, ctx context.Context, name string, data interface{}, fields ...attribute.KeyValue) (context.Context, trace.Span) { // Attributes represent additional key-value descriptors that can be bound

	tracer := otel.Tracer(name)
	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(fields...))
	span.SetName(name)
	//span.AddEvent(level, trace.WithAttributes(fields...))
	if data != nil {
		attributes := addStructAsAttributes(span, "", data)
		span.AddEvent(level, trace.WithAttributes(attributes...))
	}

	switch level {
	case "fatal":
		span.SetStatus(codes.Error, name)
	case "info":
		span.SetStatus(codes.Ok, name)
	}
	return ctx, span
}

func addStructAsAttributes(span trace.Span, prefix string, data interface{}) []attribute.KeyValue {
	v := reflect.ValueOf(data)
	var attrAttributes []attribute.KeyValue
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		fieldValue := v.Field(i)

		switch fieldValue.Kind() {
		case reflect.String:
			// span.SetAttributes(attribute.String(fieldName, fieldValue.String()))
			attrAttributes = append(attrAttributes, attribute.String(prefix+fieldName, fieldValue.String()))
		case reflect.Bool:
			//span.SetAttributes(attribute.Bool(fieldName, fieldValue.Bool()))
			attrAttributes = append(attrAttributes, attribute.Bool(prefix+fieldName, fieldValue.Bool()))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			//	span.SetAttributes(attribute.Int64(fieldName, fieldValue.Int()))
			attrAttributes = append(attrAttributes, attribute.Int64(prefix+fieldName, fieldValue.Int()))

		case reflect.Float32, reflect.Float64:
			//span.SetAttributes(attribute.Float64(fieldName, fieldValue.Float()))
			attrAttributes = append(attrAttributes, attribute.Float64(prefix+fieldName, fieldValue.Float()))

		case reflect.Slice:
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				switch elem.Kind() {
				case reflect.String:
					//	span.SetAttributes(attribute.String(fieldName+"."+strconv.Itoa(j), elem.String()))
					attrAttributes = append(attrAttributes, attribute.String(prefix+fieldName+"."+strconv.Itoa(j), elem.String()))

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					//span.SetAttributes(attribute.Int64(fieldName+"."+strconv.Itoa(j), elem.Int()))
					attrAttributes = append(attrAttributes, attribute.Int64(prefix+fieldName+"."+strconv.Itoa(j), elem.Int()))

				case reflect.Float32, reflect.Float64:
					//span.SetAttributes(attribute.Float64(fieldName+"."+strconv.Itoa(j), elem.Float()))
					attrAttributes = append(attrAttributes, attribute.Float64(prefix+fieldName+"."+strconv.Itoa(j), elem.Float()))

				case reflect.Bool:
					//span.SetAttributes(attribute.Bool(fieldName+"."+strconv.Itoa(j), elem.Bool()))
					attrAttributes = append(attrAttributes, attribute.Bool(prefix+fieldName+"."+strconv.Itoa(j), elem.Bool()))

				case reflect.Struct:
					attributes := addStructAsAttributes(span, prefix+fieldName+"."+strconv.Itoa(j)+".", elem.Interface())
					// Add other types if needed...
					attrAttributes = append(attrAttributes, attributes...)

				}
			}
		case reflect.Struct:
			attributes := addStructAsAttributes(span, prefix+fieldName+".", fieldValue.Interface())
			// Add cases for other data types as needed...
			attrAttributes = append(attrAttributes, attributes...)

		default:
			// Handle unsupported types if needed
		}
	}
	return attrAttributes
}