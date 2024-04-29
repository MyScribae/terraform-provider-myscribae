resource "myscribae_script_group" "example" {
    provider_id = myscribae_provider.example.id
    alt_id      = "example_script_group"
    name        = "Example Group"
    description = "Example group is a group of scripts"
    public      = false
}