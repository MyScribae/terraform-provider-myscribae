package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type urlValidator struct{}

func NewUrlValidator() validator.String {
	return &urlValidator{}
}

func (u *urlValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	val := req.ConfigValue.ValueString()

	if val == "" {
		resp.Diagnostics.AddError("url cannot be empty", "url provided is empty")
		return
	}

	_, err := url.ParseRequestURI(val)
	if err != nil {
		resp.Diagnostics.AddError("invalid url", fmt.Sprintf("invalid url: %s", err.Error()))
		return
	}
}

func (u *urlValidator) Description(context.Context) string {
	return "Validates a url"
}

func (u *urlValidator) MarkdownDescription(context.Context) string {
	return "Validates a url"
}
