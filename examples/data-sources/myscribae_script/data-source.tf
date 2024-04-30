data "myscribae_script" "example" {
  provider_id     = myscribae_provider.id
  script_group_id = myscribae_script_group.example.id
  alt_id          = "example_script"
}
