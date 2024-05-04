package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type xorValidator struct {
	requiresOne bool
	fields      []string
}

var _ validator.String = (*xorValidator)(nil)

func NewXorValidator(fields []string, requiresOne bool) *xorValidator {
	return &xorValidator{
		requiresOne: requiresOne,
		fields:      fields,
	}
}

func (u *xorValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var count int
	for _, field := range u.fields {
		var strVal types.String
		diags := req.Config.GetAttribute(ctx, path.Root(field), &strVal)
		if diags.HasError() {
			continue
		}
		if strVal.IsNull() || strVal.IsUnknown() || strVal.ValueString() == "" {
			continue
		}
		count++
	}

	if count == 0 {
		resp.Diagnostics.AddError("at least one field is required", fmt.Sprintf("at least one field is required from %v", u.fields))
		return
	}

	if u.requiresOne && count > 1 {
		resp.Diagnostics.AddError("only one field is allowed", fmt.Sprintf("only one field is allowed from %v", u.fields))
		return
	}
}

func (u *xorValidator) Description(context.Context) string {
	return "XOR Validator"
}

func (u *xorValidator) MarkdownDescription(context.Context) string {
	return "XOR Validator"
}
