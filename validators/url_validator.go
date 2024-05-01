package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type urlValidator struct {
	Required bool
}

func NewUrlValidator(required bool) validator.String {
	return &urlValidator{
		Required: required,
	}
}

func (u *urlValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	val := req.ConfigValue.ValueStringPointer()

	if val != nil && *val != "" {
		_, err := url.ParseRequestURI(*val)
		if err != nil {
			resp.Diagnostics.AddError("invalid url", fmt.Sprintf("invalid url: %s", err.Error()))
			return
		}
	} else if u.Required {
		resp.Diagnostics.AddError("url cannot be empty", "url provided is empty")
	}
}

func (u *urlValidator) Description(context.Context) string {
	return "Validates a url"
}

func (u *urlValidator) MarkdownDescription(context.Context) string {
	return "Validates a url"
}
