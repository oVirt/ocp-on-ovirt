variable "rhv-base-image" {
    type = "string"
    description = "image familt"
    default = "centos-cloud/centos-8"
}


variable "rhv-engine-name" {
    type = "string"
    description = "rhv engine VM instance name"
    default = "ocp-rhv44-vm-engine"
}

variable "rhv_host_count" {
    type = number
    description = "Number of host instances to launch"
    default = 2
}

variable "rhv-engine-vcpu" {
    type = number
    description = "virtual cpu count"
    default = 6
}

variable "rhv-engine-memory" {
    type = number
    description = "memory in mega"
    default =16384
}

variable "rhv-engine-disk-size" {
    type = number
    description = "rhev disk size in Giga"
    default = "60"
}

variable "gce-ssh-user" {
    type = "string"
    description = "ssh username"
    default = "centos"
}


variable "gce-ssh-pub-key-file" {
    type = "string"
    description = "ssh public key file"
    default = "~/.ssh/id_rsa.pub"
}
