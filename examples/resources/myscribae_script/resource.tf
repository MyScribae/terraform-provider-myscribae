resource "myscribae_script" "example" {
  provider_id        = myscribae_provider.example.id
  script_group_id    = myscribae_script_group.example.id
  alt_id             = "example_script_group"
  name               = "Example Group"
  description        = "Example group is a group of scripts"
  price_in_cents     = 1000
  sla_sec            = 3600
  token_lifetime_sec = 1800
  recurrence         = "monthly"
  public             = false
}