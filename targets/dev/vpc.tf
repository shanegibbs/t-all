module "vpc" {
    azs = ["a" ,"b", "c"]
    source = "../../modules/vpc"
    name = "primary"
}