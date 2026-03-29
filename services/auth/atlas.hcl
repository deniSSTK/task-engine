variable "dev_db_url" {
  type    = string
  default = "docker://postgres/18/dev?search_path=public"
}

variable "local_db_url" {
  type    = string
  default = "postgres://admin:password@localhost:5432/auth_db?sslmode=disable"
}

data "external_schema" "ent" {
  program = [
    "go",
    "run",
    "entgo.io/ent/cmd/ent",
    "schema",
    "./schema",
    "--dialect",
    "postgres",
  ]
}

env "local" {
  src = data.external_schema.ent.url

  migration {
    dir = "file://migrations"
  }

  dev = var.dev_db_url
  url = var.local_db_url

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
