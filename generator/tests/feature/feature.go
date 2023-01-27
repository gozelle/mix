package tests_feature

import (
	"context"
	"github.com/gozelle/mix/generator/tests/basic"
	"github.com/gozelle/mix/generator/tests/stringer"
	"regexp"
)

type FeatureAPI interface {
	/* Test 测试2 */
	Test(ctx context.Context, in *Feature) (out Feature, err error) // 测试2
	// 测试3
	Test2(ctx context.Context) (err error)
}

type Feature struct {
	a regexp.Regexp
	tests_basic.Basic
	*tests_stringer.Stringer
	Link           *Feature // 递归定义
	LinkArray      []*Feature
	LinkArrayArray [][]*Feature
}
