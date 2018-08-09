package collector

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/aybabtme/tracing/gen/span/me/aybabt/tracing"
	"github.com/pkg/errors"
)

type DB interface {
	ExecContext(context.Context, string, ...interface{}) (*sql.Result, error)
}

func InsertSpans(ctx context.Context, db DB, spans []*tracing.Span) error {
	for i, span := range spans {
		if err := InsertSpan(ctx, db, span); err != nil {
			return errors.Wrapf(err, "writing span %d", i)
		}
	}
	return nil
}

func InsertSpan(ctx context.Context, db DB, span *tracing.Span) error {
	buf := bytes.NewBuffer(nil)
	arger := writeInsertSpanStmt(buf, span)
	_, err := db.ExecContext(
		ctx,
		buf.String(),
		arger()...,
	)
	return err
}

const Schema = `
CREATE TABLE IF NOT EXISTS spans (
	span_id INT NOT NULL,
	trace_id INT NOT NULL
	operation VARCHAR(128),
	start_at TIMESTAMPTZ,
	finish_at TIMESTAMPTZ,
)

CREATE TABLE IF NOT EXISTS span_reference_child_ofs (
	span_id INT NOT NULL REFERENCES spans(span_id),
	parent_span_id INT NOT NULL,
)

CREATE TABLE IF NOT EXISTS span_reference_follows_froms (
	span_id INT NOT NULL REFERENCES spans(span_id),
	parent_span_id INT NOT NULL,
)

CREATE TABLE IF NOT EXISTS tags (
	span_id INT NOT NULL REFERENCES spans(span_id),
	key STRING NOT NULL,
	string_value STRING,
	numeric_value DOUBLE,
)

CREATE TABLE IF NOT EXISTS baggages (
	span_id INT NOT NULL REFERENCES spans(span_id),
	key STRING NOT NULL,
	string_value STRING,
	numeric_value DOUBLE,
)
`

func writeInsertSpanStmt(buf *bytes.Buffer, span *tracing.Span) func() (args []interface{}) {

	spanCtx := span.Context(nil)
	startAt := span.Start(nil)
	finishAt := span.Finish(nil)
	spanRef := span.Reference(nil)

	buf.WriteString(`
INSERT INTO spans(
	span_id,
	trace_id,
	operation,
	start_at,
	finish_at
) VALUES `)
	arg := 1
	arg = writeValuePlaceholders(buf, arg, 1, 5)
	buf.WriteString(" RETURNING NOTHING;")

	if spanCtx.BaggagesLength() > 0 {
		buf.WriteString(`
INSERT INTO baggages (
	span_id,
	key,
	string_value,
	numeric_value
) VALUES `)
		arg = writeValuePlaceholders(buf, arg, spanCtx.BaggagesLength(), 4)
		buf.WriteString(" RETURNING NOTHING;")
	}

	if spanRef.ChildOfLength() > 0 {
		buf.WriteString(`
INSERT INTO span_reference_child_ofs(
	span_id,
	parent_span_id
) VALUES `)
		arg = writeValuePlaceholders(buf, arg, spanRef.ChildOfLength(), 2)
		buf.WriteString(" RETURNING NOTHING;")
	}

	if spanRef.FollowsFromLength() > 0 {
		buf.WriteString(`
INSERT INTO span_reference_follows_froms(
	span_id,
	parent_span_id
) VALUES `)
		arg = writeValuePlaceholders(buf, arg, spanRef.FollowsFromLength(), 2)
		buf.WriteString(" RETURNING NOTHING;")
	}

	if span.TagsLength() > 0 {
		buf.WriteString(`
INSERT INTO tags (
	span_id,
	key,
	string_value,
	numeric_value
) VALUES `)
		arg = writeValuePlaceholders(buf, arg, span.TagsLength(), 4)
		buf.WriteString(" RETURNING NOTHING;")
	}

	return func() []interface{} {

		spanArgs := []interface{}{
			spanCtx.SpanId(),
			spanCtx.TraceId(),
			span.Operation(),
			time.Unix(startAt.UnixSecond(), startAt.UnixNanosecond()),
			time.Unix(finishAt.UnixSecond(), finishAt.UnixNanosecond()),
		}

		if spanCtx.BaggagesLength() > 0 {
			baggagesArgs := make([]interface{}, 0, spanCtx.BaggagesLength())
			kv := new(tracing.KeyValue)
			for i := 0; spanCtx.Baggages(kv, i); i++ {
				baggagesArgs = append(baggagesArgs,
					kv.Key(),
					kv.StringValue(),
					kv.NumericValue(),
				)
			}
			spanArgs = append(spanArgs, baggagesArgs...)
		}

		if spanRef.ChildOfLength() > 0 {
			childOfArgs := make([]interface{}, 0, spanRef.ChildOfLength())
			for i := 0; i < spanRef.ChildOfLength(); i++ {
				childOfArgs = append(childOfArgs,
					spanCtx.SpanId(),
					spanRef.ChildOf(i),
				)
			}
			spanArgs = append(spanArgs, childOfArgs...)
		}

		if spanRef.FollowsFromLength() > 0 {
			followsFromArgs := make([]interface{}, 0, spanRef.FollowsFromLength())
			for i := 0; i < spanRef.FollowsFromLength(); i++ {
				followsFromArgs = append(followsFromArgs,
					spanCtx.SpanId(),
					spanRef.FollowsFrom(i),
				)
			}
			spanArgs = append(spanArgs, followsFromArgs...)
		}

		if span.TagsLength() > 0 {
			tagsArgs := make([]interface{}, 0, span.TagsLength())
			kv := new(tracing.KeyValue)
			for i := 0; span.Tags(kv, i); i++ {
				tagsArgs = append(tagsArgs,
					kv.Key(),
					kv.StringValue(),
					kv.NumericValue(),
				)
			}
			spanArgs = append(spanArgs, tagsArgs...)
		}

		return spanArgs
	}
}

func writeValuePlaceholders(buf *bytes.Buffer, start, rows, args int) int {
	if start <= 0 {
		start = 1
	}
	for row := 0; row < rows; row++ {
		if row == 0 {
			buf.WriteString("(")
		} else {
			buf.WriteString(", (")
		}
		for arg := 0; arg < args; arg++ {
			if arg != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(strconv.Itoa(start))
			buf.WriteRune('$')
		}
		buf.WriteRune(')')
		start++
	}
	return start
}
