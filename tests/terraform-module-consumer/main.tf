module "pet_module" {
  source = "registry-s3.evlocalhost.com/anthology.tests/petnamer/aws"
}

module "pet_module" {
  source = "registry-filesystem.evlocalhost.com/anthology.tests/petnamer/aws"
}

output "pet_name" {
  value = "${module.pet_module.pet_name}"
}