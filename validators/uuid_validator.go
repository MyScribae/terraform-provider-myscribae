package validators

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type uuidValidator struct {
	Required bool
}

var _ validator.String = (*uuidValidator)(nil)

func NewUuidValidator(required bool) validator.String {
	return &uuidValidator{
		Required: required,
	}
}

func (u *uuidValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	val := req.ConfigValue.ValueStringPointer()

	if val != nil && *val != "" {
		_, err := uuid.Parse(*val)
		if err != nil {
			resp.Diagnostics.AddError("invalid uuid", fmt.Sprintf("invalid uuid: %s", err.Error()))
			return
		}
	} else if u.Required {
		resp.Diagnostics.AddError("uuid cannot be empty", "uuid provided is empty")
	}
}

func (u *uuidValidator) Description(context.Context) string {
	return "Validates a uuid"
}

func (u *uuidValidator) MarkdownDescription(context.Context) string {
	return "Validates a uuid"
}
