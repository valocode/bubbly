module github.com/valocode/bubbly

go 1.16

require (
	entgo.io/contrib v0.0.0-20210701194530-6b9b6b0bd76c
	entgo.io/ent v0.9.0
	github.com/99designs/gqlgen v0.13.0
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/containerd/continuity v0.1.0 // indirect
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-git/go-git/v5 v5.4.2
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/hcl/v2 v2.10.1
	github.com/hashicorp/terraform v0.15.3
	github.com/jackc/pgproto3/v2 v2.1.0 // indirect
	github.com/jackc/pgx/v4 v4.11.0
	github.com/mattn/go-sqlite3 v1.14.8
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/opencontainers/runc v1.0.0-rc95 // indirect
	github.com/ory/dockertest/v3 v3.7.0
	github.com/rs/cors v1.8.0
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/swag v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/zclconf/go-cty v1.9.0
	github.com/zclconf/go-cty-yaml v1.0.2
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)

// replace entgo.io/contrib => github.com/jlarfors/contrib v0.0.0-20210728113018-b9b57a221a03

// For local development of contrib...
replace entgo.io/contrib => /Users/jacoblarfors/work/ent/contrib
