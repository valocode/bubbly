module github.com/valocode/bubbly

go 1.16

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/Masterminds/squirrel v1.5.0
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/appleboy/gofight/v2 v2.1.2
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/cockroachdb/cockroach-go/v2 v2.1.0
	github.com/containerd/containerd v1.5.0 // indirect
	github.com/cornelk/hashmap v1.0.1
	github.com/dchest/siphash v1.2.2 // indirect
	github.com/docker/docker v20.10.6+incompatible
	github.com/fatih/color v1.10.0
	github.com/go-git/go-git/v5 v5.2.0
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/google/uuid v1.2.0
	github.com/graphql-go/graphql v0.7.9
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.10.0
	github.com/hashicorp/terraform v0.15.3
	github.com/imdario/mergo v0.3.11
	github.com/jackc/pgx/v4 v4.10.1
	github.com/labstack/echo/v4 v4.2.1
	github.com/lib/pq v1.9.0 // indirect
	github.com/likexian/gokit v0.20.15
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.9
	github.com/nats-io/nats.go v1.10.0
	github.com/ory/dockertest/v3 v3.6.3
	github.com/rs/zerolog v1.20.0
	github.com/ryanuber/columnize v2.1.2+incompatible
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/echo-swagger v1.1.0
	github.com/swaggo/swag v1.7.0
	github.com/zclconf/go-cty v1.8.3
	github.com/zclconf/go-cty-yaml v1.0.2
	github.com/ziflex/lecho/v2 v2.1.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210603125802-9665404d3644 // indirect
	golang.org/x/tools v0.1.2 // indirect
	gopkg.in/h2non/gock.v1 v1.0.16
)

replace github.com/hashicorp/hcl/v2 => github.com/valocode/hcl/v2 v2.10.0-patch-1
