"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_deviceAdapter = require("../../mixins/device-adapter.js");
const mixins_share = require("../../mixins/share.js");
const apis_content = require("../../apis/content.js");
const utils_date = require("../../utils/date.js");
const utils_request = require("../../utils/request.js");
const apis_message = require("../../apis/message.js");
const utils_share = require("../../utils/share.js");
const ActionBar = () => "../../components/action-bar/index.js";
const _sfc_main = {
  mixins: [mixins_deviceAdapter.deviceAdapter, mixins_share.shareMixin],
  components: {
    ActionBar
  },
  data() {
    return {
      contentId: null,
      contentData: {
        id: "",
        title: "",
        content: "",
        category: "",
        publisher: "",
        publishTime: "",
        publisher_name: "",
        publisher_id: "",
        publisher_avatar: "",
        images: [],
        isTop: false,
        publish_count: 0,
        comment_count: 0,
        comments: [],
        cover: ""
        // 添加封面图属性，用于分享
      },
      loading: false,
      loadError: false,
      errorMsg: "",
      isFollowed: false,
      isCollected: false,
      commentLoading: false,
      bannerList: [],
      currentBannerIndex: 0,
      isContentPage: true
      // 标记为内容页面，用于分享功能
    };
  },
  computed: {
    pageTitle() {
      var _a;
      return ((_a = this.contentData) == null ? void 0 : _a.title) || "内容详情";
    },
    pageStyle() {
      return {
        "--nav-height": `${this.layoutSize.navHeight}px`,
        "--content-padding": `${this.layoutSize.contentPadding}rpx`
      };
    }
  },
  onLoad(options) {
    if (options.id) {
      this.contentId = options.id;
      this.loadDetail();
      this.loadBannerData();
    } else {
      this.loadError = true;
      this.errorMsg = "参数错误，缺少内容ID";
    }
  },
  onShow() {
    this.updateShareData();
  },
  methods: {
    // 过滤HTML内容中的图片标签
    filterImgTags(htmlContent) {
      if (!htmlContent)
        return "";
      return htmlContent.replace(/<img[^>]*>/g, "");
    },
    // 格式化发布时间为"多久前"
    formatPublishTime(time) {
      if (!time)
        return "";
      return utils_date.formatTimeAgo(time);
    },
    // 加载内容详情
    async loadDetail() {
      if (!this.contentId) {
        this.loadError = true;
        this.errorMsg = "内容ID无效";
        return;
      }
      this.loading = true;
      try {
        const result = await apis_content.getPublicContentDetail(this.contentId);
        if (result.code === 0 && result.data) {
          const data = result.data;
          const filteredContent = this.filterImgTags(data.content);
          const formattedTime = this.formatPublishTime(data.publishTime);
          this.contentData = {
            id: data.id,
            title: data.title || "",
            content: filteredContent,
            category: data.category || "",
            publishTime: formattedTime,
            publisher: data.publisher || "",
            publisher_name: data.publisher || "",
            publisher_id: data.publisher_id || "",
            publisher_avatar: "",
            isTop: data.isTop || false,
            images: data.images || [],
            publish_count: 0,
            comment_count: 0,
            comments: [],
            cover: data.cover || (data.images && data.images.length > 0 ? data.images[0] : "")
            // 优先使用封面图，其次使用第一张图片
          };
          this.loadComments();
          this.loadFavoriteStatus();
          if (this.contentData.publisher_id) {
            this.loadPublisherInfo();
          }
          this.recordBrowseHistory();
          this.updateShareData();
        } else {
          throw new Error(result.message || "获取内容详情失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:297", "加载内容详情失败:", error);
        this.loadError = true;
        this.errorMsg = error.message || "加载失败，请重试";
      } finally {
        this.loading = false;
      }
    },
    // 加载收藏状态
    async loadFavoriteStatus() {
      try {
        const result = await apis_content.getFavoriteStatus(this.contentId);
        if (result.code === 0 && result.data) {
          this.isCollected = result.data.isFavorite;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:313", "获取收藏状态失败:", error);
      }
    },
    // 加载评论列表
    async loadComments() {
      if (!this.contentData.id)
        return;
      this.commentLoading = true;
      try {
        const result = await apis_content.getCommentList(this.contentData.id, {
          page: 1,
          pageSize: 10
        });
        if (result.code === 0 && result.data) {
          this.contentData.comments = (result.data.list || []).map((item) => ({
            id: item.id,
            avatar: item.avatarUrl || "",
            name: item.realName || "匿名用户",
            content: item.comment || "",
            time: this.formatPublishTime(item.createdAt)
          }));
          if (result.data.total !== void 0) {
            this.contentData.comment_count = result.data.total;
          }
        } else {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:344", "获取评论列表失败:", result.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:347", "加载评论失败:", error);
      } finally {
        this.commentLoading = false;
      }
    },
    handleBack() {
      common_vendor.index.navigateBack();
    },
    previewImage(index) {
      common_vendor.index.previewImage({
        current: index,
        urls: this.contentData.images,
        success: () => {
          common_vendor.index.__f__("log", "at pages/content/detail.vue:360", "图片预览成功");
        },
        fail: (err) => {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:363", "图片预览失败:", err);
        }
      });
    },
    async handleFollow() {
      try {
        const publisherId = this.contentData.publisher_id;
        if (!publisherId) {
          common_vendor.index.showToast({
            title: "无法获取发布者信息",
            icon: "none"
          });
          return;
        }
        common_vendor.index.showLoading({
          title: this.isFollowed ? "取消关注中..." : "关注中...",
          mask: true
        });
        let result;
        if (this.isFollowed) {
          result = await apis_content.unfollowUser({
            publisher_id: publisherId
          });
        } else {
          result = await apis_content.followUser({
            publisher_id: publisherId
          });
        }
        common_vendor.index.hideLoading();
        if (result && result.code === 0) {
          this.isFollowed = !this.isFollowed;
          common_vendor.index.showToast({
            title: this.isFollowed ? "已关注" : "已取消关注",
            icon: "success"
          });
        } else {
          throw new Error((result == null ? void 0 : result.message) || "操作失败");
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/content/detail.vue:418", "关注操作失败:", error);
        common_vendor.index.showToast({
          title: error.message || "操作失败，请重试",
          icon: "none"
        });
      }
    },
    handleComment() {
      common_vendor.index.showModal({
        title: "留言",
        editable: true,
        placeholderText: "请输入留言内容",
        success: async (res) => {
          if (res.confirm && res.content) {
            let loadingShown = false;
            try {
              const commentData = {
                comment: res.content,
                contentId: this.contentId
              };
              common_vendor.index.showLoading({
                title: "提交中..."
              });
              loadingShown = true;
              const result = await apis_content.createComment(commentData);
              common_vendor.index.hideLoading();
              loadingShown = false;
              if (result.code === 0) {
                common_vendor.index.showToast({
                  title: "留言成功",
                  icon: "success"
                });
                this.loadComments();
              } else {
                throw new Error(result.message || "评论失败");
              }
            } catch (error) {
              if (loadingShown) {
                common_vendor.index.hideLoading();
                loadingShown = false;
              }
              common_vendor.index.__f__("error", "at pages/content/detail.vue:473", "提交评论失败:", error);
              common_vendor.index.showToast({
                title: error.message || "留言失败，请重试",
                icon: "none"
              });
            } finally {
              if (loadingShown) {
                common_vendor.index.hideLoading();
              }
            }
          }
        }
      });
    },
    handleShare() {
      this.updateShareData();
      common_vendor.index.showShareMenu({
        withShareTicket: true,
        menus: ["shareAppMessage", "shareTimeline"]
      });
    },
    async handleMessage() {
      const publisherId = this.contentData.publisher_id || "";
      if (!publisherId) {
        common_vendor.index.showToast({
          title: "无法获取发布者信息",
          icon: "none"
        });
        return;
      }
      try {
        const result = await apis_message.createConversation({
          targetId: parseInt(publisherId)
        });
        if (result && result.code === 0) {
          common_vendor.index.navigateTo({
            url: `/pages/chat/detail?id=${publisherId}&nickName=${encodeURIComponent(this.contentData.publisher_name || "用户")}`,
            fail: (err) => {
              common_vendor.index.__f__("error", "at pages/content/detail.vue:519", "跳转到聊天页面失败:", err);
              common_vendor.index.showToast({
                title: "无法打开聊天页面",
                icon: "none"
              });
            }
          });
        } else {
          throw new Error((result == null ? void 0 : result.message) || "创建会话失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:530", "私信功能错误:", error);
        common_vendor.index.showToast({
          title: "无法发起私信，请重试",
          icon: "none"
        });
      }
    },
    // 处理收藏
    async handleCollect() {
      try {
        common_vendor.index.showLoading({
          title: this.isCollected ? "取消收藏中..." : "收藏中...",
          mask: true
        });
        let result;
        if (this.isCollected) {
          result = await apis_content.cancelFavorite(this.contentId);
        } else {
          result = await apis_content.addFavorite(this.contentId);
        }
        common_vendor.index.hideLoading();
        if (result.code === 0) {
          this.isCollected = !this.isCollected;
          common_vendor.index.showToast({
            title: this.isCollected ? "已收藏" : "已取消收藏",
            icon: "success"
          });
        } else {
          throw new Error(result.message || "操作失败");
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/content/detail.vue:571", "收藏操作失败:", error);
        common_vendor.index.showToast({
          title: error.message || "操作失败，请重试",
          icon: "none"
        });
      }
    },
    handleBannerClick(item) {
      if (!item || !item.linkUrl)
        return;
      if (item.linkType === "page") {
        common_vendor.index.navigateTo({
          url: "/" + item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/content/detail.vue:586", "页面跳转失败", item.linkUrl, err);
          }
        });
      } else if (item.linkType === "webview") {
        if (!this.isValidUrl(item.linkUrl)) {
          common_vendor.index.showToast({
            title: "无效的URL",
            icon: "none"
          });
          return;
        }
        if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.setWebViewSecurity) {
          common_vendor.wx$1.setWebViewSecurity({
            enable: true,
            success: () => {
              common_vendor.index.__f__("log", "at pages/content/detail.vue:605", "成功设置WebView安全模式");
              this.navigateToWebView(item.linkUrl);
            },
            fail: (err) => {
              common_vendor.index.__f__("error", "at pages/content/detail.vue:609", "设置WebView安全模式失败:", err);
              this.navigateToWebView(item.linkUrl);
            }
          });
        } else {
          this.navigateToWebView(item.linkUrl);
        }
      } else if (item.linkType === "miniprogram") {
        common_vendor.index.navigateToMiniProgram({
          appId: item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/content/detail.vue:623", "小程序跳转失败", item.linkUrl, err);
            common_vendor.index.showToast({
              title: "跳转失败",
              icon: "none"
            });
          }
        });
      }
    },
    // 检查URL是否有效
    isValidUrl(url) {
      if (!url)
        return false;
      try {
        new URL(url);
        return true;
      } catch (e) {
        return false;
      }
    },
    // 导航到WebView页面
    navigateToWebView(url) {
      const separator = url.includes("?") ? "&" : "?";
      const secureParams = `${separator}coop=cross-origin&coep=require-corp`;
      let secureUrl = url;
      if (url.startsWith("http://") || url.startsWith("https://")) {
        secureUrl = url + secureParams;
      }
      common_vendor.index.navigateTo({
        url: "/pages/webview/index?url=" + encodeURIComponent(secureUrl),
        fail: (err) => {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:659", "网页跳转失败", url, err);
          common_vendor.index.showToast({
            title: "网页打开失败",
            icon: "none"
          });
        }
      });
    },
    // 记录浏览历史 - 简化错误处理
    async recordBrowseHistory() {
      try {
        const contentId = this.contentData.id || this.contentId;
        if (!contentId) {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:673", "记录浏览历史失败: 缺少有效的内容ID");
          return;
        }
        await apis_content.addBrowseRecord(contentId, "article");
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:680", "记录浏览历史失败:", error);
      }
    },
    // 加载发布者信息
    async loadPublisherInfo() {
      try {
        const publisherId = this.contentData.publisher_id;
        if (!publisherId) {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:689", "无法加载发布者信息: publisher_id不存在");
          return;
        }
        const publisherIdStr = String(publisherId);
        const result = await apis_content.getPublisherInfo(publisherIdStr);
        if (result && result.code === 0 && result.data) {
          const publisherData = result.data;
          if (publisherData.real_name) {
            this.contentData.publisher_name = publisherData.real_name;
          }
          if (publisherData.avatar_url) {
            this.contentData.publisher_avatar = publisherData.avatar_url;
          }
          if (publisherData.publish_count !== void 0) {
            this.contentData.publish_count = publisherData.publish_count;
          }
          await this.loadFollowStatus(publisherIdStr);
        } else {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:720", "获取发布者信息失败:", (result == null ? void 0 : result.message) || "未知错误");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:723", "加载发布者信息异常:", error);
      }
    },
    // 加载关注状态
    async loadFollowStatus(publisherId) {
      try {
        if (!publisherId) {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:730", "无法获取关注状态: publisher_id不存在");
          return;
        }
        const result = await apis_content.getPublisherFollowStatus(publisherId);
        if (result && result.code === 0 && result.data) {
          this.isFollowed = !!result.data.is_followed;
        } else {
          common_vendor.index.__f__("error", "at pages/content/detail.vue:740", "获取关注状态失败:", (result == null ? void 0 : result.message) || "未知错误");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:743", "加载关注状态异常:", error);
      }
    },
    // 加载轮播图数据
    async loadBannerData() {
      try {
        const result = await utils_request.get("/wx/inner-banner/list", {
          bannerType: "home"
        });
        if (result && result.code === 0 && result.data) {
          if (result.data.isGlobalEnabled && result.data.list && result.data.list.length > 0) {
            this.bannerList = result.data.list.filter((item) => item.isEnabled).sort((a, b) => a.order - b.order);
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:762", "获取轮播图数据失败", error);
      }
    },
    handleSwiperChange(e) {
      this.currentBannerIndex = e.detail.current;
    },
    // 更新分享数据
    async updateShareData() {
      try {
        const settings = await utils_share.getShareSettings();
        if (settings && this.contentData) {
          const shareTitle = this.contentData.title || settings.content_share_text || "查看内容详情";
          let shareImage = "";
          if (this.contentData.cover) {
            shareImage = this.contentData.cover;
          } else if (this.contentData.images && this.contentData.images.length > 0) {
            shareImage = this.contentData.images[0];
          } else if (settings.content_share_image) {
            shareImage = settings.content_share_image;
          }
          this.shareData = {
            title: shareTitle,
            imageUrl: shareImage,
            path: `/pages/content/detail?id=${this.contentId}`
          };
          common_vendor.index.__f__("log", "at pages/content/detail.vue:800", "内容分享数据已更新:", this.shareData);
        } else {
          this.initShareData();
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/content/detail.vue:806", "更新分享数据失败:", error);
        this.initShareData();
      }
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  const _component_action_bar = common_vendor.resolveComponent("action-bar");
  (_easycom_uni_icons2 + _component_action_bar)();
}
const _easycom_uni_icons = () => "../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.p({
      type: "left",
      size: "20",
      color: "#333333"
    }),
    b: common_vendor.o((...args) => $options.handleBack && $options.handleBack(...args)),
    c: common_vendor.t($options.pageTitle),
    d: _ctx.statusBarHeight + "px",
    e: $data.loading
  }, $data.loading ? {} : {}, {
    f: $data.loadError
  }, $data.loadError ? {
    g: common_vendor.t($data.errorMsg),
    h: common_vendor.o((...args) => $options.loadDetail && $options.loadDetail(...args))
  } : {}, {
    i: !$data.loading && !$data.loadError
  }, !$data.loading && !$data.loadError ? common_vendor.e({
    j: $data.contentData.isTop
  }, $data.contentData.isTop ? {} : {}, {
    k: common_vendor.t($data.contentData.category),
    l: common_vendor.t($data.contentData.title),
    m: common_vendor.t($data.contentData.publisher),
    n: common_vendor.t($data.contentData.publishTime),
    o: $data.contentData.content
  }, $data.contentData.content ? {
    p: $data.contentData.content
  } : {}, {
    q: $data.contentData.images && $data.contentData.images.length
  }, $data.contentData.images && $data.contentData.images.length ? {
    r: common_vendor.f($data.contentData.images, (img, index, i0) => {
      return {
        a: index,
        b: img,
        c: common_vendor.o(($event) => $options.previewImage(index), index)
      };
    })
  } : {}, {
    s: $data.contentData.publisher_avatar,
    t: common_vendor.t($data.contentData.publisher_name),
    v: common_vendor.t($data.contentData.publish_count),
    w: common_vendor.t($data.isFollowed ? "已关注" : "关注"),
    x: common_vendor.o((...args) => $options.handleFollow && $options.handleFollow(...args)),
    y: $data.bannerList.length > 0
  }, $data.bannerList.length > 0 ? common_vendor.e({
    z: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: item.image,
        b: index,
        c: common_vendor.o(($event) => $options.handleBannerClick(item), index)
      };
    }),
    A: common_vendor.o((...args) => $options.handleSwiperChange && $options.handleSwiperChange(...args)),
    B: $data.bannerList.length > 1
  }, $data.bannerList.length > 1 ? {
    C: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: index,
        b: $data.currentBannerIndex === index ? 1 : ""
      };
    })
  } : {}) : {}, {
    D: common_vendor.t($data.contentData.comment_count),
    E: common_vendor.p({
      type: "chat",
      size: "14",
      color: "#007AFF"
    }),
    F: common_vendor.o((...args) => $options.handleComment && $options.handleComment(...args)),
    G: $data.contentData.comments && $data.contentData.comments.length
  }, $data.contentData.comments && $data.contentData.comments.length ? {
    H: common_vendor.f($data.contentData.comments, (comment, index, i0) => {
      return {
        a: comment.avatar,
        b: common_vendor.t(comment.name),
        c: common_vendor.t(comment.content),
        d: common_vendor.t(comment.time),
        e: index
      };
    })
  } : {}) : {}, {
    I: _ctx.navigationBarHeight + "px",
    J: common_vendor.o($options.handleComment),
    K: common_vendor.o($options.handleCollect),
    L: common_vendor.o($options.handleShare),
    M: common_vendor.o($options.handleMessage),
    N: common_vendor.p({
      ["is-collected"]: $data.isCollected,
      publisher: {
        id: $data.contentData.publisher_id,
        name: $data.contentData.publisher_name
      }
    }),
    O: common_vendor.s($options.pageStyle)
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/content/detail.js.map
