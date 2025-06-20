module github.com/HollyEllmo/my-first-go-project

go 1.24.2

require (
	github.com/HollyEllmo/my-proto-repo/gen/go/filter v0.0.0
	github.com/HollyEllmo/my-proto-repo/gen/go/prod_service v0.0.0-20250619092020-52c6c1eb3d32
	github.com/Masterminds/squirrel v1.5.4
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgx/v4 v4.18.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.9
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.11.1
	github.com/sirupsen/logrus v1.9.3
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.8.1
	golang.org/x/sync v0.15.0
	google.golang.org/grpc v1.73.0
)

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.26.3 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/HollyEllmo/my-proto-repo/gen/go/prod_service => github.com/HollyEllmo/my_proto_repo/gen/go/prod_service v0.0.0-20250619092020-52c6c1eb3d32

replace github.com/HollyEllmo/my-proto-repo/gen/go/filter => github.com/HollyEllmo/my_proto_repo/gen/go/filter v0.0.0-20250619092020-52c6c1eb3d32
