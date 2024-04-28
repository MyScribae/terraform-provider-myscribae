package provider_test

// import (
// 	"testing"
// )

// func TestProviderDatasource(t *testing.T) {
// 	acctest.Test(t, acctest.TestCase{
// 		Providers: map[string]acctest.ProviderFactory{
// 			"myscribae": func() *schema.Provider {
// 				return myscribae.Provider()
// 			},
// 		},
// 		Steps: []acctest.Step{
// 			{
// 				Config: `
// 				data "myscribae_provider" "provider" {
// 					alt_id = "provider"

// 					category = "provider"

// 					name = "provider"

// 					description = "provider"

// 					logo_url = "provider"

// 					banner_url = "provider"

// 					url = "provider"

// 					color = "provider"

// 					public = true

// 					account_service = true
// 				}
// 				`,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "alt_id", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "category", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "name", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "description", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "logo_url", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "banner_url", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "url", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "color", "provider"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "public", "true"),
// 					resource.TestCheckResourceAttr("data.myscribae_provider.provider", "account_service", "true"),
// 				),
// 			},
// 		},
// 	})
// }
