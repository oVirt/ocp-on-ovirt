provider "google" {
  project = "openshift-gce-devel"
  region  = "us-central1"
  zone    = "us-central1-c"
}
resource "google_compute_network" "rhv-network"{
    name = "ocp-rhv44-network"
}

resource "google_compute_subnetwork" "rhv-subnetwork"{
    name = "ocp-rhv44-subnetwork"
    network = "${google_compute_network.rhv-network.self_link}"
    ip_cidr_range = "10.0.0.0/24"
    private_ip_google_access = true
}

#firewall rules from all sources
resource "google_compute_firewall" "ocp-rhv44-firewall-all" {
  name    = "ocp-rhv44-firewall-all"
  network = "${google_compute_network.rhv-network.self_link}"

  allow {
    protocol = "tcp"
    ports    = ["6443", "443", "22","35357","9696","6642"]
  }

  source_ranges = [ "0.0.0.0/0" ]
}

#allow all communications internally
resource "google_compute_firewall" "ocp-rhv44-firewall-internal" {
  name    = "ocp-rhv44-firewall-internal"
  network = "${google_compute_network.rhv-network.self_link}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    #ports    = ["6443", "443", "80"]
  }

  allow {
    protocol = "udp"
  }

  source_ranges = [ "10.0.0.0/24" ]
}


resource "google_compute_instance" "engine-instance" {
  name         = "${var.rhv-engine-name}"
  machine_type = "custom-${var.rhv-engine-vcpu}-${var.rhv-engine-memory}"
  boot_disk {
    initialize_params {
      image = "${var.rhv-base-image}"
      size = "${var.rhv-engine-disk-size}"
    }
  }
network_interface {
    # A default network is created for all GCP projects
    network       = "${google_compute_network.rhv-network.self_link}"
    subnetwork = "${google_compute_subnetwork.rhv-subnetwork.self_link}"
    network_ip = "10.0.0.10"
  }

 labels = {
      rhv_role = "rhv-engine"
  }
 metadata = {
    ssh-keys =  "${var.gce-ssh-user}:${file(var.gce-ssh-pub-key-file)}"
  }
}


