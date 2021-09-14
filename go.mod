module github.com/valocode/bubbly

go 1.16

require (
	entgo.io/contrib v0.1.1-0.20210909110117-6470c1c113ab
	entgo.io/ent v0.9.2-0.20210821141344-368a8f7a2e9a
	github.com/99designs/gqlgen v0.14.0
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/MakeNowJust/heredoc/v2 v2.0.1
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/bytecodealliance/wasmtime-go v0.29.0 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/docker/docker v20.10.8+incompatible // indirect
	github.com/fatih/color v1.12.0
	github.com/go-bindata/go-bindata v1.0.1-0.20190711162640-ee3c2418e368 // indirect
	github.com/go-chi/chi v3.3.2+incompatible // indirect
	github.com/go-git/go-git/v5 v5.4.2
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gofrs/uuid v3.3.0+incompatible // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.10.1 // indirect
	github.com/jackc/pgproto3/v2 v2.1.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/labstack/echo/v4 v4.5.0
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/moby/buildkit v0.9.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/open-policy-agent/conftest v0.27.0
	github.com/open-policy-agent/opa v0.31.0
	github.com/opencontainers/runc v1.0.0-rc95 // indirect
	github.com/ory/dockertest/v3 v3.7.0
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rs/zerolog v1.23.0
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tmccombs/hcl2json v0.3.3 // indirect
	github.com/vektah/gqlparser/v2 v2.2.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/xiaoqidun/entps v0.0.0-20210811073333-c1af85abbd8d
	github.com/zclconf/go-cty v1.9.1 // indirect
	github.com/ziflex/lecho/v2 v2.5.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210817190340-bfb29a6856f2 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	golang.org/x/text v0.3.7 // indirect
)

// replace entgo.io/contrib => github.com/jlarfors/contrib v0.0.0-20210728113018-b9b57a221a03

// For local development of contrib...
//replace entgo.io/contrib => /Users/jacoblarfors/work/ent/contrib
