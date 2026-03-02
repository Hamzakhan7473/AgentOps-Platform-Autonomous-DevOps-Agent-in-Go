module github.com/agentops/platform

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.32.6
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.7.0
	github.com/sashabaranov/go-openai v1.27.0
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/exporters/prometheus v0.45.2
	go.opentelemetry.io/otel/sdk/metric v1.32.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.12
)
