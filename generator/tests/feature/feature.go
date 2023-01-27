package tests_feature

import (
	"context"
	"github.com/gozelle/mix/generator/tests/basic"
	"github.com/gozelle/mix/generator/tests/stringer"
	"io"
	"regexp"
)

type FeatureAPI interface {
	tests_basic.BasicAPI
	Sub
	Test(ctx context.Context, in *Feature) (out Feature, err error) // 测试2
}

type Sub interface {
	Download(ctx context.Context) (io.Reader, error)
	Query() <-chan any
}

type Feature struct {
	a regexp.Regexp
	tests_basic.Basic
	*tests_stringer.Stringer
	Link           *Feature // 递归定义
	LinkArray      []*Feature
	LinkArrayArray [][]*Feature
}
