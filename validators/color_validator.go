package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type colorValidator struct{}

var _ validator.String = (*colorValidator)(nil)

func NewColorValidator() validator.String {
	return &colorValidator{}
}

func (u *colorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	val := req.ConfigValue.ValueString()

	if val == "" {
		resp.Diagnostics.AddError("color cannot be empty", "color provided is empty")
		return
	}

	if len(val) != 7 {
		resp.Diagnostics.AddError("invalid color", "color must be 7 characters long")
		return
	}

	if val[0] != '#' {
		resp.Diagnostics.AddError("invalid color", "color must start with #")
		return
	}

	for i, c := range val[1:] {
		if c < '0' || c > '9' {
			if c < 'a' || c > 'f' {
				resp.Diagnostics.AddError("invalid color", "color must be a valid hex color")
				return
			}
		}
		if i > 5 {
			resp.Diagnostics.AddError("invalid color", "color must be 7 characters long")
			return
		}
	}
}

func (u *colorValidator) Description(context.Context) string {
	return "Validates a hex color"
}

func (u *colorValidator) MarkdownDescription(context.Context) string {
	return "Validates a hex color"
}
