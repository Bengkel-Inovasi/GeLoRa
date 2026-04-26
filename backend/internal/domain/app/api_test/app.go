package domainappapitest

import (
	"context"

	"github.com/Bengkel-Inovasi/GeLoRa/backend/internal/config"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
)

const path = "domain/app/api_test"

type ApiTest struct {
	infrastructure *Infrastructure
	wiring         *Wiring
}

func Run() int {
	const tag = path + "/Run"
	a := &ApiTest{}
	ctx := context.Background()

	if err := config.LoadEnv(); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to load .env", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := a.NewInfrastructure(ctx); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create infrastructure", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := a.NewWiring(ctx); err != nil {
		logTmp(nil).Error(ctx, tag, "Failed to create wiring", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	if err := a.test(ctx); err != nil {
		a.wiring.log.Error(ctx, tag, "API test failed", domainmodel.LogMeta{"error": err.Error()})
		return 1
	}

	a.wiring.log.Info(ctx, tag, "All API tests passed", nil)
	return 0
}

func (a *ApiTest) test(ctx context.Context) error {
	const tag = path + "/test"

	// Sign in as super user
	authResp, err := a.wiring.appself.PostAuthSignIn(ctx, "super", "super@gelora")
	if err != nil {
		a.wiring.log.Error(ctx, tag, "PostAuthSignIn failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.appself.SetAccessToken(authResp.AccessToken)
	a.wiring.appself.SetRefreshToken(authResp.RefreshToken)
	a.wiring.log.Info(ctx, tag, "PostAuthSignIn passed", nil)

	// Get own user info
	userMe, err := a.wiring.appself.GetUserInfoMe(ctx)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "GetUserInfoMe failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.log.Info(ctx, tag, "GetUserInfoMe passed", domainmodel.LogMeta{"id": userMe.Id, "username": userMe.Username})

	// List users
	usersList, err := a.wiring.appself.GetUsersList(ctx, 1, 10, nil, nil, nil)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "GetUsersList failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.log.Info(ctx, tag, "GetUsersList passed", domainmodel.LogMeta{"total": usersList.Pagination.Total})

	// List nodes
	nodesList, err := a.wiring.appself.GetNodesList(ctx, 1, 10, nil, nil, nil)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "GetNodesList failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.log.Info(ctx, tag, "GetNodesList passed", domainmodel.LogMeta{"total": nodesList.Pagination.Total})

	// List sessions
	sessionsList, err := a.wiring.appself.GetSessionsList(ctx, 1, 10, nil, nil, nil, nil)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "GetSessionsList failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.log.Info(ctx, tag, "GetSessionsList passed", domainmodel.LogMeta{"total": sessionsList.Pagination.Total})

	// List records
	records, err := a.wiring.appself.GetRecords(ctx, nil, nil, nil, nil, nil)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "GetRecords failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.log.Info(ctx, tag, "GetRecords passed", domainmodel.LogMeta{"count": len(records.Data)})

	// Refresh token
	refreshResp, err := a.wiring.appself.PostAuthRefresh(ctx, authResp.RefreshToken)
	if err != nil {
		a.wiring.log.Error(ctx, tag, "PostAuthRefresh failed", domainmodel.LogMeta{"error": err.Error()})
		return err
	}
	a.wiring.appself.SetAccessToken(refreshResp.AccessToken)
	a.wiring.appself.SetRefreshToken(refreshResp.RefreshToken)
	a.wiring.log.Info(ctx, tag, "PostAuthRefresh passed", nil)

	return nil
}
