//go:build contract

package client_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bigwhite/shortlink/pkg/client"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceClient_Contract(t *testing.T) {
	// 1. 创建 Pact Mock Server
	// 修正: 使用 MockHTTPProviderConfig 而不是 V2PactConfig
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "ShortlinkService",
		Provider: "UserService",
		Host:     "127.0.0.1",
	})
	assert.NoError(t, err)

	// 2. 定义交互 (Interaction)
	err = mockProvider.
		AddInteraction().
		Given("user 123 exists and has creation permission").
		UponReceiving("a request to check user 123's permission").
		WithRequest(http.MethodGet, "/users/123/permissions").
		WillRespondWith(http.StatusOK, func(b *consumer.V4ResponseBuilder) {
			// 修正: 使用 matchers 包中的匹配器，而不是 dsl 包
			b.Header("Content-Type", matchers.String("application/json"))
			// 修正: 使用 matchers.MapMatcher 和 matchers.Like
			b.JSONBody(matchers.MapMatcher{
				"can_create": matchers.Like(true),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			// 创建客户端并执行测试
			userClient := &client.UserServiceClient{
				BaseURL:    fmt.Sprintf("http://%s:%d", config.Host, config.Port),
				HTTPClient: &http.Client{},
			}

			canCreate, err := userClient.CanCreateLink(context.Background(), "123")

			assert.NoError(t, err)
			assert.True(t, canCreate)

			return nil
		})

	assert.NoError(t, err)
}
