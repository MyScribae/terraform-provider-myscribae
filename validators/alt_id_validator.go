package validators

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type altIdValidator struct {
	required bool
}

var _ validator.String = (*altIdValidator)(nil)

func NewAltIdValidator(required bool) validator.String {
	return &altIdValidator{
		required: required,
	}
}

func (u *altIdValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	valPtr := req.ConfigValue.ValueStringPointer()

	if valPtr == nil {
		if !u.required {
			resp.Diagnostics.AddError("alt_id cannot be empty", "alt_id provided is empty")
			return
		}

		return
	}
	val := *valPtr

	if val == "" {
		resp.Diagnostics.AddError("alt_id cannot be empty", "alt_id provided is empty")
		return
	}

	lowerSnakeCaseRegex := regexp.MustCompile(`^[a-z0-9_]+$`)
	if !lowerSnakeCaseRegex.MatchString(val) {
		resp.Diagnostics.AddError("invalid alt_id", "alt_id must be lower snake case")
		return
	}

	if len(val) > 50 {
		resp.Diagnostics.AddError("invalid alt_id", "alt_id must be less than 50 characters")
		return
	}

	if len(val) < 1 {
		resp.Diagnostics.AddError("invalid alt_id", "alt_id must be at least 1 character")
		return
	}

	if val[0] == '_' {
		resp.Diagnostics.AddError("invalid alt_id", "alt_id cannot start with an underscore")
		return
	}

	if val[len(val)-1] == '_' {
		resp.Diagnostics.AddError("invalid alt_id", "alt_id cannot end with an underscore")
		return
	}
}

func (u *altIdValidator) Description(context.Context) string {
	return "Validates an alt_id"
}

func (u *altIdValidator) MarkdownDescription(context.Context) string {
	return "Validates an alt_id"
}
