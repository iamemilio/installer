variable "cidr_block" {
  type = string
}

variable "cluster_id" {
  type = string
}

variable "external_network" {
  description = "Name of the external network providing Floating IP addresses."
  type        = string
  default     = ""
}

variable "external_network_id" {
  description = "UUID of the external network providing Floating IP addresses."
  type        = string
  default     = ""
}

variable "lb_floating_ip" {
  description = "(optional) Existing floating IP address to attach to the load balancer created by the installer."
  type        = string
  default     = ""
}

variable "enable_bootstrap_floating_ip" {
  description = "(optional) If true the bootstrap machine gets a floating IP address that will be used to collect logs."
  type        = bool
  default     = true
}

variable "masters_count" {
  type = string
}

variable "external_dns" {
  type = list(string)
}

variable "trunk_support" {
  type = string
}

variable "octavia_support" {
  type = "string"
}
