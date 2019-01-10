resource "random_pet" "my_pet" {
}

output "pet_name" {
  value = "${random_pet.my_pet.id}"
}