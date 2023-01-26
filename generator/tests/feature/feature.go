package tests_feature

import (
	"context"
	"github.com/gozelle/mix/generator/tests/basic"
	"github.com/gozelle/mix/generator/tests/stringer"
)

type FeatureAPI interface {
	Test(ctx context.Context, in *Feature) (out Feature, err error)
}

type Feature struct {
	tests_basic.Basic
	*tests_stringer.Stringer
	Link           *Feature
	LinkArray      []*Feature
	LinkArrayArray [][]*Feature
}
