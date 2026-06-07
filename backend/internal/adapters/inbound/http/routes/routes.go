package adaptersinboundhttproutes

import (
	"net/http"

	_ "github.com/Bengkel-Inovasi/GeLoRa/backend/docs/swagger"
	adaptersinboundhttphandler "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/handler"
	adaptersinboundhttpmiddlewareauth "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/auth"
	adaptersinboundhttpmiddlewarelogger "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/logger"
	adaptersinboundhttpmiddlewarerolerule "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/adapters/inbound/http/middleware/rolerule"
	domainmodel "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/domain/model"
	portsinboundhttp "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/inbound/http"
	portsoutboundlogging "github.com/Bengkel-Inovasi/GeLoRa/backend/internal/ports/outbound/logging"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func Route(
	router *gin.RouterGroup,
	log portsoutboundlogging.Generic,
	svcUser portsinboundhttp.User,
	hdlUser *adaptersinboundhttphandler.User,
	hdlNode *adaptersinboundhttphandler.Node,
	hdlSession *adaptersinboundhttphandler.Session,
	hdlRecord *adaptersinboundhttphandler.Record,
	hdlAlert *adaptersinboundhttphandler.Alert,
) {
	docs := router.Group("/docs")
	{
		docs.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/docs/swagger/index.html")
		})
		docs.GET("/swagger", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/docs/swagger/index.html")
		})
		docs.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	auth := router.Group(
		"/auth",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"Auth request succeeded",
			"Auth request failed",
			"Auth server error",
		),
	)
	{
		auth.POST(
			"/sign-in",
			hdlUser.PostSignIn,
		)
		auth.POST(
			"/sign-up",
			hdlUser.PostSignUp,
		)
		auth.POST(
			"/refresh",
			hdlUser.PostRefresh,
		)
	}

	users := router.Group(
		"/users",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"User request succeeded",
			"User request failed",
			"User server error",
		),
		adaptersinboundhttpmiddlewareauth.Jwt(),
	)
	{
		users.GET(
			"/me",
			hdlUser.GetUserInfoMe,
		)
		users.PATCH(
			"/me",
			hdlUser.PatchUserInfoMe,
		)
		users.PUT(
			"/me/password",
			hdlUser.PutUserPasswordMe,
		)
		users.GET(
			"",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlUser.GetUsersList,
		)
		users.GET(
			"/:id",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlUser.GetUserInfoById,
		)
		users.PATCH(
			"/:id",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleClaimsRoleHigherThanUserIdPath("id", svcUser),
				},
			}),
			hdlUser.PatchUserInfoById,
		)
		users.PUT(
			"/:id/password",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleClaimsRoleHigherThanUserIdPath("id", svcUser),
				},
			}),
			hdlUser.PutUserPasswordById,
		)
		users.PUT(
			"/:id/role",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleClaimsRoleHigherThanUserIdPath("id", svcUser),
				},
			}),
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleClaimsRoleHigherThanBodyRole("role"),
				},
			}),
			hdlUser.PutUserRoleById,
		)
		users.DELETE(
			"/:id",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleClaimsRoleHigherThanUserIdPath("id", svcUser),
				},
			}),
			hdlUser.DeleteUserById,
		)
	}

	nodes := router.Group(
		"/nodes",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"Node request succeeded",
			"Node request failed",
			"Node server error",
		),
		adaptersinboundhttpmiddlewareauth.Jwt(),
	)
	{
		nodes.POST(
			"/register",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlNode.PostRegisterClimber,
		)
		nodes.GET(
			"",
			hdlNode.GetNodesList,
		)
		nodes.GET(
			"/:id",
			hdlNode.GetNodeInfoById,
		)
		nodes.PATCH(
			"/:id",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlNode.PatchNodeInfoById,
		)
		nodes.PUT(
			"/:id/validate",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlNode.PutNodeValidate,
		)
		nodes.DELETE(
			"/:id",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{
					Role: domainmodel.UserRoleSuper,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
				{
					Role: domainmodel.UserRoleAdmin,
					Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
				},
			}),
			hdlNode.DeleteNodeById,
		)
	}

	sessions := router.Group(
		"/sessions",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"Session request succeeded",
			"Session request failed",
			"Session server error",
		),
		adaptersinboundhttpmiddlewareauth.Jwt(),
		adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
			{
				Role: domainmodel.UserRoleSuper,
				Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
			},
			{
				Role: domainmodel.UserRoleAdmin,
				Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
			},
			{
				Role: domainmodel.UserRoleClient,
				Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
			},
		}),
	)
	{
		sessions.POST(
			"",
			hdlSession.PostSession,
		)
		sessions.GET(
			"",
			hdlSession.GetSessionsList,
		)
		sessions.GET(
			"/:id",
			hdlSession.GetSessionById,
		)
		sessions.PUT(
			"/end",
			hdlSession.PutSessionEnd,
		)
		sessions.DELETE(
			"/:id",
			hdlSession.DeleteSessionById,
		)
	}

	alerts := router.Group(
		"/alerts",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"Alert request succeeded",
			"Alert request failed",
			"Alert server error",
		),
		adaptersinboundhttpmiddlewareauth.Jwt(),
	)
	{
		// All authenticated users can report an emergency
		alerts.POST("", hdlAlert.PostAlert)

		// Only admin/super can see and acknowledge alerts
		alerts.GET(
			"",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{Role: domainmodel.UserRoleSuper, Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll()},
				{Role: domainmodel.UserRoleAdmin, Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll()},
			}),
			hdlAlert.GetAlertsList,
		)
		alerts.PUT(
			"/:id/acknowledge",
			adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
				{Role: domainmodel.UserRoleSuper, Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll()},
				{Role: domainmodel.UserRoleAdmin, Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll()},
			}),
			hdlAlert.PutAlertAcknowledge,
		)
	}

	records := router.Group(
		"/records",
		adaptersinboundhttpmiddlewarelogger.Logger(
			log,
			"Record request succeeded",
			"Record request failed",
			"Record server error",
		),
		adaptersinboundhttpmiddlewareauth.Jwt(),
		adaptersinboundhttpmiddlewarerolerule.RoleRule([]adaptersinboundhttpmiddlewarerolerule.RoleRuler{
			{
				Role: domainmodel.UserRoleSuper,
				Rule: adaptersinboundhttpmiddlewarerolerule.RoleAllowAll(),
			},
			{
				Role: domainmodel.UserRoleAdmin,
				Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
			},
			{
				Role: domainmodel.UserRoleClient,
				Rule: adaptersinboundhttpmiddlewarerolerule.RuleAllowAll(),
			},
		}),
	)
	{
		records.GET(
			"",
			hdlRecord.GetRecords,
		)
	}
}
