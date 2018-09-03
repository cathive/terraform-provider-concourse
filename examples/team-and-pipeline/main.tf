provider "concourse" {
  target = "local"
}

resource "concourse_team" "avengers" {
  name = "avengers"
}

resource "concourse_pipeline" "batman" {
  team = "${concourse_team.avengers.name}"
  name = "batman"
}