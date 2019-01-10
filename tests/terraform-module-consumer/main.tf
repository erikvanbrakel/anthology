module "pet_module" {
  source = "registry-s3.anthology.localtest.me/anthology/petnamer/aws"
}

module "pet_module" {
  source = "registry-filesystem.anthology.localtest.me/anthology/petnamer/aws"
}

output "pet_name" {
  value = "${module.pet_module.pet_name}"
}
