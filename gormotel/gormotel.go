package gormotel

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	gormSpanKey        = "gorm_span"
	callBackBeforeName = "trace:before"
	callBackAfterName  = "trace:after"
)

// tracer is the global tracer used by the GORM plugin.
var tracer = otel.Tracer(gormSpanKey)

type TracePlugin struct{}

var _ gorm.Plugin = &TracePlugin{}

// Name
// @Description: 返回插件名称
// @Auth shigx 2023-08-10 17:47:48
// @return string
func (op *TracePlugin) Name() string {
	return "TracePlugin"
}

// Initialize
// @Description: 初始化
// @Auth shigx 2023-08-10 17:48:06
// @param db
// @return err
func (op *TracePlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	_ = db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	_ = db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	_ = db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	_ = db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	_ = db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	_ = db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	_ = db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	_ = db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	_ = db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	_ = db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	_ = db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	_ = db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

// before
// @Description: 前置钩子执行方法
// @Auth shigx 2023-08-10 17:48:21
// @param db
func before(db *gorm.DB) {
	dbCtx := db.Statement.Context
	if !trace.SpanFromContext(dbCtx).IsRecording() {
		return
	}

	ctx, span := tracer.Start(dbCtx, "GORM SQL")
	span.SetAttributes(
		attribute.String("db.system", "gorm"),
	)

	db.InstanceSet(gormSpanKey, ctx)
	return
}

// after
// @Description: 后置钩子执行方法
// @Auth shigx 2023-08-10 17:49:30
// @param db
func after(db *gorm.DB) {
	ctx, ok := db.InstanceGet(gormSpanKey)
	if !ok {
		return
	}
	span := trace.SpanFromContext(ctx.(context.Context))
	defer span.End()

	// sql
	span.SetAttributes(
		attribute.Int64("db.rows", db.RowsAffected),
		attribute.String("db.sql", db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)),
	)

	// error
	if db.Error != nil && !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		span.RecordError(db.Error)
		span.SetStatus(codes.Error, db.Error.Error())
	}

	return
}
