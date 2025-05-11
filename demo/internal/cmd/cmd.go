package cmd

import (
	"context"
	"net/http"
	"strings"
	"time"

	backendLogic "demo/internal/logic/backend"
	bottomTabLogic "demo/internal/logic/bottom_tab"
	browseHistoryLogic "demo/internal/logic/browse_history"
	butlerLogic "demo/internal/logic/butler"
	clientLogic "demo/internal/logic/client"
	commentLogic "demo/internal/logic/comment"
	contentLogic "demo/internal/logic/content"
	contentClientLogic "demo/internal/logic/content/client"
	dbToolLogic "demo/internal/logic/dbtool"
	exchangeLogic "demo/internal/logic/exchange"
	favoriteLogic "demo/internal/logic/favorite"
	logLogic "demo/internal/logic/log"
	mallLogic "demo/internal/logic/mall"
	orderLogic "demo/internal/logic/order"
	packageLogic "demo/internal/logic/package"
	paymentLogic "demo/internal/logic/payment"
	_ "demo/internal/logic/product" // 自动初始化商品逻辑
	regionLogic "demo/internal/logic/region"
	settingsLogic "demo/internal/logic/settings"
	statsLogic "demo/internal/logic/stats"
	storageLogic "demo/internal/logic/storage"
	userLogic "demo/internal/logic/user"
	wechatPayLogic "demo/internal/logic/wechat_pay"

	"demo/internal/middleware"
	"demo/internal/service"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"

	backend "demo/internal/controller/backend"
	backendController "demo/internal/controller/backend"
	bottomTabController "demo/internal/controller/bottom_tab"
	browseHistoryController "demo/internal/controller/browse_history"
	butlerController "demo/internal/controller/butler"
	"demo/internal/controller/client"
	clientController "demo/internal/controller/client"
	commentController "demo/internal/controller/comment"
	contentController "demo/internal/controller/content"
	contentClientController "demo/internal/controller/content/client"
	exchangeController "demo/internal/controller/exchange"
	licenseController "demo/internal/controller/license"
	logController "demo/internal/controller/log"
	mallController "demo/internal/controller/mall"
	orderController "demo/internal/controller/order"
	package_controller "demo/internal/controller/package"
	paymentController "demo/internal/controller/payment"
	clientOrderController "demo/internal/controller/payment/client"
	wechatController "demo/internal/controller/payment/wechat"
	publisherController "demo/internal/controller/publisher"
	regionController "demo/internal/controller/region"
	settingsController "demo/internal/controller/settings"
	storageController "demo/internal/controller/storage"
	userController "demo/internal/controller/user"
	wxController "demo/internal/controller/wx"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 设置全局时区为中国时区
			loc, err := time.LoadLocation("Asia/Shanghai")
			if err != nil {
				glog.Fatal(ctx, "加载时区失败:", err)
				return err
			}
			time.Local = loc

			s := g.Server()

			// 初始化数据库
			glog.Info(ctx, "正在检查数据库状态...")
			if err := dbToolLogic.InitDatabase(ctx); err != nil {
				glog.Fatal(ctx, "数据库初始化失败:", err)
				return err
			}

			// 初始化用户服务
			service.SetUser(userLogic.New())

			// 初始化日志服务
			service.SetLog(logLogic.New())

			// 初始化客户服务
			service.SetClient(clientLogic.New())

			// 初始化存储服务
			service.SetStorage(storageLogic.New())

			// 初始化支付服务
			service.SetPayment(paymentLogic.New())

			// 初始化订单服务
			service.SetOrder(orderLogic.New())

			// 初始化奖励过期清理定时任务
			settingsLogic.InitRewardCleanTask()

			// 初始化微信支付服务
			service.SetWechatPay(wechatPayLogic.New())

			// 初始化后端服务
			service.SetBackend(backendLogic.New())

			// 初始化内容管理服务
			service.SetContent(contentLogic.New())
			service.SetCategory(contentLogic.NewCategory())

			// 初始化客户端内容服务
			service.SetContentClient(contentClientLogic.New())

			// 初始化收藏服务
			service.SetFavorite(favoriteLogic.New())

			// 初始化套餐服务
			service.SetPackage(packageLogic.New())

			// 初始化地区服务
			service.SetRegion(regionLogic.New())

			// 初始化评论服务
			service.SetComment(commentLogic.New())

			// 初始化浏览历史记录服务
			service.SetBrowseHistory(browseHistoryLogic.New())

			// 初始化底部导航栏服务
			service.SetBottomTab(bottomTabLogic.New())

			// 初始化兑换记录服务
			service.RegisterExchange(exchangeLogic.New())

			// 初始化商城分类服务
			service.RegisterMall(mallLogic.New())

			// 初始化专属管家服务
			service.SetButler(butlerLogic.New())

			// 初始化统计数据服务
			service.RegisterStats(statsLogic.New())

			// 检查是否启用API文档
			apiEnabled := g.Cfg().MustGet(ctx, "server.apiEnabled", true).Bool()

			// 配置OpenAPI文档 - 根据配置决定是否启用API文档
			if apiEnabled {
				s.SetOpenApiPath("/api.json")
				s.SetSwaggerPath("/swagger")
				g.Log().Info(ctx, "API文档已启用，访问路径: /api.json, /swagger")
			} else {
				// 禁用API文档访问
				s.SetOpenApiPath("")
				s.SetSwaggerPath("")
				g.Log().Info(ctx, "API文档已禁用")
			}

			// 添加Swagger中间件用于在OpenAPI Json中注入安全定义
			s.Use(middleware.SwaggerMiddleware)

			// 添加CORS中间件进行域名访问控制
			s.Use(middleware.CorsMiddleware)

			// 添加授权验证中间件，确保所有API请求都经过授权验证
			s.Use(middleware.LicenseAuth)

			// 添加XSS防护中间件，对请求参数进行安全过滤
			// 可以通过配置文件 manifest/config/config.yaml 中的 security.xss 配置项来控制
			if g.Cfg().MustGet(ctx, "security.xss.enabled", true).Bool() {
				glog.Info(ctx, "XSS防护中间件已启用")
				s.Use(middleware.XssMiddleware)
			}

			// 注册中间件
			s.Use(ghttp.MiddlewareHandlerResponse)

			// 记录当前工作目录和执行文件目录
			pwd := gfile.Pwd()
			execPath := gfile.SelfDir()
			g.Log().Info(ctx, "当前工作目录:", pwd)
			g.Log().Info(ctx, "执行文件目录:", execPath)

			// 设置OpenAPI和Swagger路径
			if apiEnabled {
				swaggerPath := "resource/public/swagger"
				if gfile.Exists(swaggerPath) {
					s.AddStaticPath("/swagger", swaggerPath)
				} else if gfile.Exists(gfile.Join(pwd, swaggerPath)) {
					s.AddStaticPath("/swagger", gfile.Join(pwd, swaggerPath))
				} else if gfile.Exists(gfile.Join(execPath, swaggerPath)) {
					s.AddStaticPath("/swagger", gfile.Join(execPath, swaggerPath))
				} else {
					g.Log().Warning(ctx, "静态资源目录 'resource/public/swagger' 不存在，已跳过")
				}
			} else {
				g.Log().Info(ctx, "API文档已禁用，跳过添加Swagger静态资源路径")
			}

			// 设置导出文件的静态目录
			exportPath := "resource/export"
			if gfile.Exists(exportPath) {
				s.AddStaticPath("/export", exportPath)
			} else if gfile.Exists(gfile.Join(pwd, exportPath)) {
				s.AddStaticPath("/export", gfile.Join(pwd, exportPath))
			} else if gfile.Exists(gfile.Join(execPath, exportPath)) {
				s.AddStaticPath("/export", gfile.Join(execPath, exportPath))
			} else {
				g.Log().Warning(ctx, "静态资源目录 'resource/export' 不存在，已跳过")
			}

			// 设置前端资源目录
			assetsPath := "resource/dist/assets"
			if gfile.Exists(assetsPath) {
				s.AddStaticPath("/assets", assetsPath)
			} else if gfile.Exists(gfile.Join(pwd, assetsPath)) {
				s.AddStaticPath("/assets", gfile.Join(pwd, assetsPath))
			} else if gfile.Exists(gfile.Join(execPath, assetsPath)) {
				s.AddStaticPath("/assets", gfile.Join(execPath, assetsPath))
			} else {
				g.Log().Warning(ctx, "静态资源目录 'resource/dist/assets' 不存在，已跳过")
			}

			// 设置前端静态资源目录
			distPath := "resource/dist"
			if gfile.Exists(distPath) {
				s.AddStaticPath("/", distPath)
			} else if gfile.Exists(gfile.Join(pwd, distPath)) {
				s.AddStaticPath("/", gfile.Join(pwd, distPath))
			} else if gfile.Exists(gfile.Join(execPath, distPath)) {
				s.AddStaticPath("/", gfile.Join(execPath, distPath))
			} else {
				g.Log().Warning(ctx, "静态资源目录 'resource/dist' 不存在，已跳过")
			}

			// 创建控制器实例
			userController := userController.New().V1()
			logController := logController.New().V1()
			clientController := clientController.New().V1()
			storageController := storageController.New().V1()
			paymentController := paymentController.New().V1()
			orderController := orderController.New()
			wechatPayController := wechatController.New()
			backendController := backendController.New()
			contentController := contentController.New().V1()
			contentClientController := contentClientController.New().V1()
			regionController := regionController.New().V1()
			commentController := commentController.New()
			browseHistoryController := browseHistoryController.New().V1()
			publisherController := publisherController.New().V1()
			settingsController := settingsController.New().V1()
			bottomTabController := bottomTabController.New().V1()
			productController := mallController.New()
			wxProductController := wxController.New()
			exchangeRecordController := &exchangeController.Controller{}
			shopCategoryController := &mallController.Controller{}
			licenseController := licenseController.NewController()
			statsController := &backend.StatsController{}

			// 授权验证API路由组
			s.Group("/license", func(group *ghttp.RouterGroup) {
				// 添加响应处理中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				// 注册授权验证路由
				licenseController.Register(group)
			})

			// API路由组 - 公开接口（登录、验证码和客户端内容）
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 登录和验证码接口
				group.GET("/captcha", userController.GetCaptcha)
				group.POST("/login", userController.Login)

				// 微信小程序登录接口
				group.POST("/wx/wxapp-login", clientController.WxappLogin)

				// 客户端内容公开接口
				group.GET("/wx/client/content/categories", contentClientController.CategoryList)
				group.GET("/wx/client/content/region/list", contentClientController.RegionContentList)
				group.GET("/wx/client/content/region/idle/list", contentClientController.RegionIdleList)
				group.GET("/wx/client/content/public/detail", contentClientController.ContentPublicDetail)
				group.GET("/wx/client/package/list", contentClientController.WxClientPackageList)

				// 客户端地区接口（公开）
				group.GET("/wx/client/region/list", regionController.WxClientList)

				// 微信客户端-活动区域接口（公开）
				group.GET("/wx/activity-area/get", settingsController.WxGetActivityArea)

				// 微信客户端-导航小程序接口（公开）
				group.GET("/wx/mini-program/list", settingsController.WxGetMiniProgramList)

				// 微信客户端-轮播图接口（公开）
				group.GET("/wx/banner/list", settingsController.WxGetBannerList)

				// 微信客户端-内页轮播图接口（公开）
				group.GET("/wx/inner-banner/list", settingsController.WxGetInnerBannerList)

				// 微信客户端-底部导航栏接口（公开）
				group.GET("/wx/bottom-tab/list", bottomTabController.WxGetBottomTabList)

				// 微信客户端-协议接口（公开）
				group.GET("/wx/agreement", settingsController.WxGetAgreement)

				// 微信客户端-小程序基础设置接口（公开）
				group.GET("/wx/mini-program/base/settings", settingsController.WxGetMiniProgramBaseSettings)

				// 微信客户端-广告设置接口（公开）
				group.GET("/wx/ad/settings", settingsController.WxGetAdSettings)

				// 微信客户端-分享设置接口（公开）
				group.GET("/wx/share/settings", settingsController.WxGetShareSettings)

				// 微信客户端-兑换记录公开接口（公开）
				group.GET("/wx/exchange-record/list", exchangeRecordController.WxGetPublicList)
			})

			// API路由组 - 用户认证相关接口（需要认证）
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 退出登录接口
				group.POST("/logout", userController.Logout)

				// 获取用户个人信息接口
				group.GET("/user/info", userController.GetUserInfo)
			})

			// API路由组 - 用户管理接口（需要认证）
			s.Group("/user", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 用户管理接口
				group.GET("/list", userController.List)
				group.POST("/create", userController.Create)
				group.PUT("/update", userController.Update)
				group.DELETE("/delete", userController.Delete)
			})

			// API路由组 - 支付设置接口（需要认证）
			s.Group("/payment", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 添加响应处理中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// 获取支付配置
				group.Bind(
					paymentController.GetConfig,
					paymentController.SaveConfig,
				)
			})

			// API路由组 - 客户管理接口（需要认证）
			s.Group("/client", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 客户管理接口
				group.GET("/list", clientController.List)
				group.POST("/create", clientController.Create)
				group.PUT("/update", clientController.Update)
				group.DELETE("/delete", clientController.Delete)

				// 客户时长接口
				group.GET("/duration/list", clientController.DurationList)
				group.GET("/duration/detail", clientController.DurationDetail)
				group.POST("/duration/create", clientController.DurationCreate)
				group.PUT("/duration/update", clientController.DurationUpdate)
				group.DELETE("/duration/delete", clientController.DurationDelete)
			})

			// API路由组 - 微信小程序配置接口（需要认证）
			s.Group("/wxapp", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 添加响应处理中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// 绑定微信小程序配置接口（修复重复前缀问题）
				group.Bind(
					service.Client().GetWxappConfig,
					service.Client().SaveWxappConfig,
				)
			})

			// API路由组 - 客户接口（需要客户认证）
			s.Group("/wx", func(group *ghttp.RouterGroup) {
				// 添加客户JWT认证中间件
				group.Middleware(middleware.ClientAuth)

				// 客户信息接口
				group.GET("/client/info", clientController.Info)
				// 更新个人信息接口
				group.PUT("/client/update-profile", clientController.UpdateProfile)
				// 获取剩余时长接口
				group.GET("/client/duration/remaining", clientController.WxGetRemainingDuration)

				// 浏览历史记录接口
				group.GET("/client/browse-history/list", browseHistoryController.List)
				group.POST("/client/browse-history/clear", browseHistoryController.Clear)
				group.POST("/client/browse-history/add", browseHistoryController.Add)
				group.GET("/client/browse-history/count", browseHistoryController.Count)

				// 收藏相关接口
				favoriteController := contentClientController.Favorite()
				group.Bind(
					favoriteController.Add,
					favoriteController.Cancel,
					favoriteController.GetStatus,
					favoriteController.GetList,
					favoriteController.GetCount,
				)

				// 广告观看完成接口
				group.POST("/ad/reward/viewed", settingsController.WxRewardAdViewed)

				// 微信客户端兑换记录接口
				group.GET("/client/exchange-record/page", exchangeRecordController.WxGetPage)
				group.POST("/client/exchange-record/create", exchangeRecordController.WxCreate)

				// 微信客户端商城分类接口
				group.GET("/client/shop-category/list", shopCategoryController.WxGetCategoryList)

				// 微信客户端我的发布接口
				group.GET("/client/publish/list", contentClientController.WxMyPublishList)
				group.GET("/client/publish/count", contentClientController.WxMyPublishCount)

				// 微信客户端消息相关接口
				messageController := new(client.MessageController)
				messageController.Register(group)
			})

			// 发布人信息相关接口组
			s.Group("/wx/publisher", func(group *ghttp.RouterGroup) {
				// 添加响应处理中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// 发布人信息查询接口（不需要认证）
				group.GET("/info", publisherController.PublisherInfo)

				// 需要认证的发布人接口子组
				authGroup := group.Group("/")
				authGroup.Middleware(middleware.ClientAuth)

				// 发布人关注/取消关注接口（需要认证）
				authGroup.POST("/follow", publisherController.FollowPublisher)
				authGroup.POST("/unfollow", publisherController.UnfollowPublisher)

				// 获取关注状态接口（需要认证）
				authGroup.GET("/follow/status", publisherController.FollowStatus)

				// 获取关注列表和总数（需要认证）
				authGroup.GET("/following/list", publisherController.FollowingList)
				authGroup.GET("/following/count", publisherController.FollowingCount)
			})

			// 微信客户端内容发布接口（需要客户认证）
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 添加客户JWT认证中间件
				group.Middleware(middleware.ClientAuth)

				// 内容发布接口
				group.POST("/wx/client/content/idle/create", contentClientController.WxIdleCreate)
				group.POST("/wx/client/content/info/create", contentClientController.WxInfoCreate)

				// 图片上传接口
				group.POST("/wx/upload/image", storageController.Client.WxUploadImage)
			})

			// API路由组 - 存储配置接口（需要认证）
			s.Group("/storage", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 添加响应处理中间件
				group.Middleware(ghttp.MiddlewareHandlerResponse)

				// 绑定存储配置接口
				group.Bind(
					storageController.GetConfig,
					storageController.SaveConfig,
				)
			})

			// 注册套餐管理路由
			package_controller.Package.Register(s.Group("/"))

			// 注册订单管理接口
			orderController.Register(s.Group("/"))

			// 注册微信支付接口
			wechatPayController.Register(s.Group("/"))

			// 注册客户端订单接口
			clientOrderController := clientOrderController.NewAPI()
			clientOrderController.Register(s.Group("/"))

			// 注册评论管理接口
			commentController.Register(s.Group("/"))

			// API路由组 - 城市系统基础设置接口（需要认证）
			s.Group("/content", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 导航小程序相关接口
				group.GET("/mini-program/list", settingsController.GetMiniProgramList)
				group.POST("/mini-program/create", settingsController.CreateMiniProgram)
				group.PUT("/mini-program/update", settingsController.UpdateMiniProgram)
				group.DELETE("/mini-program/delete", settingsController.DeleteMiniProgram)
				group.PUT("/mini-program/status/update", settingsController.UpdateMiniProgramStatus)
				group.PUT("/mini-program/global-status/update", settingsController.UpdateMiniProgramGlobalStatus)

				// 轮播图相关接口
				group.GET("/banner/list", settingsController.GetBannerList)
				group.POST("/banner/create", settingsController.CreateBanner)
				group.PUT("/banner/update", settingsController.UpdateBanner)
				group.DELETE("/banner/delete", settingsController.DeleteBanner)
				group.PUT("/banner/status/update", settingsController.UpdateBannerStatus)
				group.PUT("/banner/global-status/update", settingsController.UpdateBannerGlobalStatus)

				// 内页轮播图相关接口
				group.GET("/inner-banner/list", settingsController.GetInnerBannerList)
				group.POST("/inner-banner/create", settingsController.CreateInnerBanner)
				group.PUT("/inner-banner/update", settingsController.UpdateInnerBanner)
				group.DELETE("/inner-banner/delete", settingsController.DeleteInnerBanner)
				group.PUT("/inner-banner/status/update", settingsController.UpdateInnerBannerStatus)
				group.PUT("/inner-banner/home/global-status/update", settingsController.UpdateHomeInnerBannerGlobalStatus)
				group.PUT("/inner-banner/idle/global-status/update", settingsController.UpdateIdleInnerBannerGlobalStatus)

				// 活动区域相关接口
				group.GET("/activity-area/get", settingsController.GetActivityArea)
				group.POST("/activity-area/save", settingsController.SaveActivityArea)
				group.PUT("/activity-area/global-status/update", settingsController.UpdateActivityAreaGlobalStatus)

				// 底部导航栏相关接口
				group.GET("/bottom-tab/list", bottomTabController.GetBottomTabList)
				group.POST("/bottom-tab/create", bottomTabController.CreateBottomTab)
				group.PUT("/bottom-tab/update", bottomTabController.UpdateBottomTab)
				group.DELETE("/bottom-tab/delete", bottomTabController.DeleteBottomTab)
				group.PUT("/bottom-tab/status/update", bottomTabController.UpdateBottomTabStatus)
			})

			// 小程序基础设置接口（需要认证）
			s.Group("/settings", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 小程序基础设置接口
				group.GET("/mini-program/base/settings", settingsController.GetMiniProgramBaseSettings)
				group.POST("/mini-program/base/settings/save", settingsController.SaveMiniProgramBaseSettings)

				// 广告设置接口
				group.GET("/ad/settings", settingsController.GetAdSettings)
				group.POST("/ad/settings/save", settingsController.SaveAdSettings)

				// 奖励设置接口
				group.GET("/reward/settings", settingsController.GetRewardSettings)
				group.POST("/reward/settings/save", settingsController.SaveRewardSettings)

				// 协议设置接口
				group.GET("/agreement/settings", settingsController.GetAgreementSettings)
				group.POST("/agreement/settings/save", settingsController.SaveAgreementSettings)

				// 分享设置接口
				group.GET("/share/settings", settingsController.GetShareSettings)
				group.POST("/share/settings/save", settingsController.SaveShareSettings)
			})

			// API路由组 - 日志管理接口（需要认证）
			s.Group("/log", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 操作日志接口
				group.GET("/list", logController.List)
				group.DELETE("/delete", logController.Delete)
				group.GET("/export", logController.Export)
			})

			// 注册路由组
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Auth)

				// 文件管理接口
				group.Bind(
					backendController.V1().FileUpload,
					backendController.V1().FileList,
					backendController.V1().FileDetail,
					backendController.V1().FileDelete,
					backendController.V1().FileBatchDelete,
					backendController.V1().FileUpdatePublic,
				)

				// 内容管理接口
				group.Bind(
					// 内容相关接口
					contentController.List,
					contentController.Detail,
					contentController.Create,
					contentController.Update,
					contentController.Delete,
					contentController.UpdateStatus,
					contentController.UpdateRecommend,
					contentController.GetCategories,

					// 首页分类相关接口
					contentController.HomeList,
					contentController.HomeCreate,
					contentController.HomeUpdate,
					contentController.HomeDelete,

					// 闲置分类相关接口
					contentController.IdleList,
					contentController.IdleCreate,
					contentController.IdleUpdate,
					contentController.IdleDelete,
				)

				// 地区管理接口
				group.Bind(
					regionController.List,
					regionController.Detail,
					regionController.Create,
					regionController.Update,
					regionController.Delete,
				)

				// 兑换记录接口
				group.GET("/exchange-record/get", exchangeRecordController.Get)
				group.POST("/exchange-record/create", exchangeRecordController.Create)
				group.PUT("/exchange-record/update", exchangeRecordController.Update)
				group.DELETE("/exchange-record/delete", exchangeRecordController.Delete)
				group.PUT("/exchange-record/status/update", exchangeRecordController.UpdateStatus)
				group.GET("/exchange-record/list", exchangeRecordController.GetList)

				// 统计数据接口
				group.GET("/stats/data", statsController.GetStatsData)
				group.GET("/stats/trend", statsController.GetStatsTrend)

				// 商城分类接口
				group.GET("/shop-category/get", shopCategoryController.GetCategory)
				group.POST("/shop-category/create", shopCategoryController.CreateCategory)
				group.PUT("/shop-category/update", shopCategoryController.UpdateCategory)
				group.DELETE("/shop-category/delete", shopCategoryController.DeleteCategory)
				group.PUT("/shop-category/status/update", shopCategoryController.UpdateCategoryStatus)
				group.GET("/shop-category/list", shopCategoryController.GetCategoryList)
			})

			// 注册商品管理接口（需要认证）
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 添加JWT认证中间件
				group.Middleware(middleware.Auth)

				// 商品管理接口
				group.GET("/product/list", productController.List)
				group.GET("/product/detail", productController.Detail)
				group.POST("/product/create", productController.Create)
				group.POST("/product/update", productController.Update)
				group.POST("/product/delete", productController.Delete)
				group.POST("/product/status", productController.UpdateStatus)
			})

			// 注册微信小程序端商品接口（公开）
			s.Group("/wx", func(group *ghttp.RouterGroup) {
				// 微信小程序商品接口
				group.GET("/product/list", wxProductController.List)
				group.GET("/product/detail", wxProductController.Detail)
			})

			// 注册专属管家API
			s.Group("/", func(group *ghttp.RouterGroup) {
				// 客户端接口 - 获取专属管家图片（公开接口，无需认证）
				group.GET("/wx/client/butler/image", butlerController.Butler.GetImage)

				// 管理端接口组（需要认证）
				adminGroup := group.Group("/")
				adminGroup.Middleware(middleware.Auth) // 添加JWT认证中间件

				// 注册管理端API
				adminGroup.GET("/butler", butlerController.Butler.Get)
				adminGroup.POST("/butler/image/save", butlerController.Butler.SaveImage)
			})

			// 添加操作日志中间件，放在全局使用，但在路由注册后添加
			s.Use(middleware.OperationLog)

			// 处理前端路由，将所有未匹配的路由请求返回前端index.html
			s.BindHandler("/*", func(r *ghttp.Request) {
				// 如果请求的路径不存在且没有文件扩展名，返回index.html
				if r.StaticFile != nil {
					return
				}

				// 检查是否是静态资源请求
				path := r.URL.Path
				if strings.HasSuffix(path, ".js") ||
					strings.HasSuffix(path, ".css") ||
					strings.HasSuffix(path, ".png") ||
					strings.HasSuffix(path, ".jpg") ||
					strings.HasSuffix(path, ".jpeg") ||
					strings.HasSuffix(path, ".gif") ||
					strings.HasSuffix(path, ".svg") ||
					strings.HasSuffix(path, ".ico") {
					r.Response.WriteStatus(http.StatusNotFound)
					return
				}

				// 所有非静态资源请求都指向前端路由
				indexFile := "resource/dist/index.html"
				if gfile.Exists(indexFile) {
					r.Response.ServeFile(indexFile)
					return
				}

				// 尝试从当前工作目录查找
				pwd := gfile.Pwd()
				indexFromPwd := gfile.Join(pwd, indexFile)
				if gfile.Exists(indexFromPwd) {
					r.Response.ServeFile(indexFromPwd)
					return
				}

				// 尝试从执行文件所在目录查找
				execPath := gfile.SelfDir()
				indexFromExec := gfile.Join(execPath, indexFile)
				if gfile.Exists(indexFromExec) {
					r.Response.ServeFile(indexFromExec)
					return
				}

				// 前端文件不存在时，返回404状态码
				g.Log().Notice(r.Context(), "前端文件不存在，请确认资源目录中包含index.html文件")
				g.Log().Debug(r.Context(), "尝试查找的路径:", indexFile, indexFromPwd, indexFromExec)
				r.Response.WriteStatus(http.StatusNotFound)
			})

			// 启动服务
			s.Run()
			return nil
		},
	}
)
