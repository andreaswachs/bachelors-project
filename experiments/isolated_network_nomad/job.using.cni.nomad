job "docs" {
  datacenters = ["dc1"]

    group "example" {
      count = 1

        network {
          port "http" {
            to = "8181"
          }
        }

      task "server" {
        driver = "docker"

        resources {
          cpu    = 600
          memory = 128
        }

        config {
          image = "hashicorp/http-echo"
          network_mode = "ingress"
          ipv4_address = "192.168.137.232"
          ports = ["http"]
          args = [
          "-listen",
          ":8181",
          "-text",
          "hello world",
          ]
        }

      }
    }
}

