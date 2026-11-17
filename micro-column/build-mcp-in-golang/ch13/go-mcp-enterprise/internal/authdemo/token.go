package authdemo

import (
	"context"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
)

type Identity struct {
	UserID   string
	TenantID string
	Scopes   []string
}

// demoTokens 仅用于本地教学和测试。
// 生产环境应替换为企业 IdP 或 OAuth Authorization Server 的 Token 验证逻辑。
var demoTokens = map[string]Identity{
	"demo-alice-token": {
		UserID:   "alice",
		TenantID: "acme",
		Scopes:   []string{"orders:read", "orders:write"},
	},
	"demo-bob-token": {
		UserID:   "bob",
		TenantID: "acme",
		Scopes:   []string{"orders:read"},
	},
	"demo-charlie-token": {
		UserID:   "charlie",
		TenantID: "globex",
		Scopes:   []string{"orders:read"},
	},
}

func VerifyToken(
	ctx context.Context,
	token string,
	req *http.Request,
) (*auth.TokenInfo, error) {
	identity, ok := demoTokens[token]
	if !ok {
		return nil, auth.ErrInvalidToken
	}

	return &auth.TokenInfo{
		UserID:     identity.UserID,
		Scopes:     identity.Scopes,
		Expiration: time.Now().Add(time.Hour),
		Extra: map[string]any{
			"tenant_id": identity.TenantID,
		},
	}, nil
}
