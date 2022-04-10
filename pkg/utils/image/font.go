package image

import (
	"context"

	"github.com/golang/freetype/truetype"

	"idas/pkg/utils/sets"
)

func LoadSystemFonts(ctx context.Context, fontNames ...string) (*truetype.Font, error) {
	newFontNames := sets.New[string](fontNames...)
	return loadSystemFonts(ctx, newFontNames)
}
