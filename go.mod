module cartloom

go 1.21

require (
	github.com/aws/aws-sdk-go-v2 v1.30.0
	github.com/aws/aws-sdk-go-v2/config v1.25.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.30.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/segmentio/kafka-go v0.4.35
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.22
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.8
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.12
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.12
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.13
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.14
	github.com/aws/aws-sdk-go-v2/service/sso v1.22.0
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.26.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.30.0
	github.com/aws/smithy-go v1.22.0
	github.com/beorn7/perks v1.0.1
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/dghubble/oauth1 v0.7.3
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f
	github.com/jmespath/go-jmespath v0.4.0
	github.com/klauspost/compress v1.17.9
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/pierrec/lz4/v4 v4.1.15
	github.com/prometheus/client_golang v1.20.4
	github.com/prometheus/client_model v0.6.1
	github.com/prometheus/common v0.55.0
	github.com/prometheus/procfs v0.15.1
	golang.org/x/sys v0.22.0
	golang.org/x/time v0.7.0
	google.golang.org/protobuf v1.34.2
)

require github.com/joho/godotenv v1.5.1 // indirect
