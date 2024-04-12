package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type recurrenceValidator struct{}

var _ validator.String = (*recurrenceValidator)(nil)

func NewRecurrenceValidator() validator.String {
	return &recurrenceValidator{}
}

func (u *recurrenceValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// valid recurrences
	validRecurrences := []string{
		"lifetime",
		"yearly",
		"monthly",
		"weekly",
	}

	val := req.ConfigValue.ValueString()

	if val == "" {
		resp.Diagnostics.AddError("recurrence cannot be empty", "recurrence provided is empty")
		return
	}

	valid := false
	for _, r := range validRecurrences {
		if r == val {
			valid = true
			break
		}
	}

	if !valid {
		resp.Diagnostics.AddError("invalid recurrence", "recurrence must be one of lifetime, yearly, monthly, weekly")
		return
	}
}

func (u *recurrenceValidator) Description(context.Context) string {
	return "Validates a recurrence"
}

func (u *recurrenceValidator) MarkdownDescription(context.Context) string {
	return "Validates a recurrence"
}
