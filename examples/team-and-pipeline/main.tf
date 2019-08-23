provider "concourse" {
  target = "local"
}

resource "concourse_team" "avengers" {
  name = "avengers"
}

resource "concourse_pipeline" "ironman" {
  team = concourse_team.avengers.name
  name = "ironman"
}