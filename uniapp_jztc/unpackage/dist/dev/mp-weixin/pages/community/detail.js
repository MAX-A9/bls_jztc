"use strict";
const common_vendor = require("../../common/vendor.js");
const mixins_deviceAdapter = require("../../mixins/device-adapter.js");
const mixins_share = require("../../mixins/share.js");
const apis_content = require("../../apis/content.js");
const utils_date = require("../../utils/date.js");
const utils_request = require("../../utils/request.js");
const utils_share = require("../../utils/share.js");
const ActionBar = () => "../../components/action-bar/index.js";
const _sfc_main = {
  mixins: [mixins_deviceAdapter.deviceAdapter, mixins_share.shareMixin],
  components: {
    ActionBar
  },
  onLoad(options) {
    const id = options.id;
    common_vendor.index.__f__("log", "at pages/community/detail.vue:193", "商品ID:", id);
    if (id) {
      this.goodsId = id;
      this.loadGoodsData(id);
      this.loadBannerData();
    } else {
      common_vendor.index.showToast({
        title: "缺少商品ID",
        icon: "none"
      });
    }
  },
  data() {
    return {
      goodsId: null,
      goodsData: {
        id: "",
        title: "",
        price: "0",
        originalPrice: "0",
        condition: "",
        location: "",
        detailLocation: "",
        tradeMethod: "线下交易",
        description: "",
        images: [],
        publisher_avatar: "",
        publisher_name: "",
        publisher_id: "",
        publish_count: 0,
        likes: 0,
        viewCount: 0,
        comments: [],
        cover: ""
        // 用于分享图片
      },
      loading: false,
      loadError: false,
      errorMsg: "",
      isFollowed: false,
      isCollected: false,
      currentSwiperIndex: 0,
      commentLoading: false,
      defaultLocation: "江西南昌-青山湖区-高新大道1888号",
      bannerList: [],
      currentBannerIndex: 0,
      isContentPage: true,
      // 标记为内容页面，用于分享功能
      shareData: null
    };
  },
  computed: {
    pageTitle() {
      return "闲置详情";
    },
    pageStyle() {
      return {
        "--nav-height": `${this.layoutSize.navHeight}px`,
        "--content-padding": `${this.layoutSize.contentPadding}rpx`
      };
    },
    // 处理HTML内容，添加类名以适配小程序
    processedDescription() {
      if (!this.goodsData.description)
        return "";
      let content = this.goodsData.description.replace(/<p/g, '<p class="rich-paragraph"').replace(/<span/g, '<span class="rich-text"').replace(/<div/g, '<div class="rich-paragraph"');
      return content;
    }
  },
  methods: {
    handleBack() {
      common_vendor.index.navigateBack();
    },
    previewImage(index) {
      common_vendor.index.previewImage({
        current: index,
        urls: this.goodsData.images
      });
    },
    async handleFollow() {
      try {
        const publisherId = this.goodsData.publisher_id;
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
          await this.loadFollowStatus(String(publisherId));
          common_vendor.index.showToast({
            title: this.isFollowed ? "已关注" : "已取消关注",
            icon: "success"
          });
        } else {
          throw new Error((result == null ? void 0 : result.message) || "操作失败");
        }
      } catch (error) {
        common_vendor.index.hideLoading();
        common_vendor.index.__f__("error", "at pages/community/detail.vue:329", "关注操作失败:", error);
        common_vendor.index.showToast({
          title: error.message || "操作失败，请重试",
          icon: "none"
        });
      }
    },
    async handleCollect() {
      try {
        common_vendor.index.showLoading({
          title: this.isCollected ? "取消收藏中..." : "收藏中...",
          mask: true
        });
        let result;
        if (this.isCollected) {
          result = await apis_content.cancelFavorite(this.goodsId);
        } else {
          result = await apis_content.addFavorite(this.goodsId);
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
        common_vendor.index.__f__("error", "at pages/community/detail.vue:369", "收藏操作失败:", error);
        common_vendor.index.showToast({
          title: error.message || "操作失败，请重试",
          icon: "none"
        });
      }
    },
    handleMessage() {
      const publisherId = this.goodsData.publisher_id || "";
      if (!publisherId) {
        common_vendor.index.showToast({
          title: "无法获取发布者信息",
          icon: "none"
        });
        return;
      }
      common_vendor.index.navigateTo({
        url: `/pages/chat/detail?userId=${publisherId}&userName=${this.goodsData.publisher_name}`,
        success: () => {
          common_vendor.index.__f__("log", "at pages/community/detail.vue:391", "跳转到私信页面成功");
        },
        fail: (err) => {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:394", "跳转到私信页面失败:", err);
          common_vendor.index.showToast({
            title: "无法打开私信",
            icon: "none"
          });
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
    // 使用API加载商品数据
    async loadGoodsData(id) {
      this.loading = true;
      try {
        const result = await apis_content.getPublicContentDetail(id);
        if (result.code === 0 && result.data) {
          const data = result.data;
          const filteredContent = this.filterImgTags(data.content);
          const formattedTime = this.formatPublishTime(data.publishTime);
          this.goodsData = {
            id: data.id,
            title: data.title || "",
            price: data.price || 0,
            originalPrice: data.originalPrice || 0,
            description: filteredContent || "",
            images: data.images || [],
            condition: data.condition || "",
            location: data.location || "",
            detailLocation: data.tradePlace || this.defaultLocation,
            tradeMethod: data.tradeMethod || "线下交易",
            publisher_avatar: data.publisherAvatar || "",
            publisher_name: data.publisher || "匿名用户",
            publisher_id: data.publisherId || data.publisher_id || "",
            publish_count: data.publishCount || 0,
            likes: data.likes || 0,
            viewCount: data.views || 0,
            comments: [],
            cover: data.cover || ""
            // 更新cover字段
          };
          this.loadComments();
          this.loadFavoriteStatus();
          if (this.goodsData.publisher_id) {
            this.loadPublisherInfo();
          }
          this.updateShareData();
        } else {
          throw new Error(result.message || "获取商品详情失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:479", "加载商品详情失败:", error);
        common_vendor.index.showToast({
          title: "加载商品详情失败",
          icon: "none"
        });
      } finally {
        this.loading = false;
      }
    },
    // 加载收藏状态
    async loadFavoriteStatus() {
      if (!this.goodsId)
        return;
      try {
        const result = await apis_content.getFavoriteStatus(this.goodsId);
        if (result.code === 0 && result.data) {
          this.isCollected = result.data.isFavorite;
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:499", "获取收藏状态失败:", error);
      }
    },
    // 加载评论列表
    async loadComments() {
      if (!this.goodsData.id)
        return;
      this.commentLoading = true;
      try {
        const result = await apis_content.getCommentList(this.goodsData.id, {
          page: 1,
          pageSize: 10
        });
        if (result.code === 0 && result.data) {
          this.goodsData.comments = (result.data.list || []).map((item) => ({
            id: item.id,
            avatar: item.avatarUrl || "",
            name: item.realName || "匿名用户",
            content: item.comment || "",
            time: this.formatPublishTime(item.createdAt)
          }));
          if (result.data.total !== void 0) {
            this.goodsData.commentCount = result.data.total;
          }
        } else {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:530", "获取评论列表失败:", result.message);
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:533", "加载评论失败:", error);
      } finally {
        this.commentLoading = false;
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
                contentId: this.goodsData.id
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
              common_vendor.index.__f__("error", "at pages/community/detail.vue:587", "提交评论失败:", error);
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
    handleSwiperChange(e) {
      this.currentSwiperIndex = e.detail.current;
    },
    handleBannerClick(item) {
      if (!item || !item.linkUrl)
        return;
      if (item.linkType === "page") {
        common_vendor.index.navigateTo({
          url: "/" + item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/community/detail.vue:613", "页面跳转失败", item.linkUrl, err);
          }
        });
      } else if (item.linkType === "webview") {
        common_vendor.index.navigateTo({
          url: "/pages/webview/index?url=" + encodeURIComponent(item.linkUrl),
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/community/detail.vue:620", "网页跳转失败", item.linkUrl, err);
          }
        });
      } else if (item.linkType === "miniprogram") {
        common_vendor.index.navigateToMiniProgram({
          appId: item.linkUrl,
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/community/detail.vue:628", "小程序跳转失败", item.linkUrl, err);
            common_vendor.index.showToast({
              title: "跳转失败",
              icon: "none"
            });
          }
        });
      }
    },
    handleBannerSwiperChange(e) {
      this.currentBannerIndex = e.detail.current;
    },
    copyLocation() {
      common_vendor.index.setClipboardData({
        data: this.goodsData.detailLocation || this.defaultLocation,
        success: () => {
          common_vendor.index.showToast({
            title: "地址已复制",
            icon: "success"
          });
        }
      });
    },
    // 加载发布者信息
    async loadPublisherInfo() {
      try {
        const publisherId = this.goodsData.publisher_id;
        if (!publisherId) {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:657", "无法加载发布者信息: publisher_id不存在");
          return;
        }
        const publisherIdStr = String(publisherId);
        const result = await apis_content.getPublisherInfo(publisherIdStr);
        if (result && result.code === 0 && result.data) {
          const publisherData = result.data;
          if (publisherData.real_name) {
            this.goodsData.publisher_name = publisherData.real_name;
          }
          if (publisherData.avatar_url) {
            this.goodsData.publisher_avatar = publisherData.avatar_url;
          }
          if (publisherData.publish_count !== void 0) {
            this.goodsData.publish_count = publisherData.publish_count;
          }
          await this.loadFollowStatus(publisherIdStr);
        } else {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:688", "获取发布者信息失败:", (result == null ? void 0 : result.message) || "未知错误");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:691", "加载发布者信息异常:", error);
      }
    },
    // 加载关注状态
    async loadFollowStatus(publisherId) {
      try {
        if (!publisherId) {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:698", "无法获取关注状态: publisher_id不存在");
          return;
        }
        const result = await apis_content.getPublisherFollowStatus(publisherId);
        if (result && result.code === 0 && result.data) {
          this.isFollowed = !!result.data.is_followed;
        } else {
          common_vendor.index.__f__("error", "at pages/community/detail.vue:708", "获取关注状态失败:", (result == null ? void 0 : result.message) || "未知错误");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:711", "加载关注状态异常:", error);
      }
    },
    // 加载轮播图数据
    async loadBannerData() {
      try {
        const result = await utils_request.get("/wx/inner-banner/list", {
          bannerType: "idle"
        });
        if (result && result.code === 0 && result.data) {
          if (result.data.isGlobalEnabled && result.data.list && result.data.list.length > 0) {
            this.bannerList = result.data.list.filter((item) => item.isEnabled).sort((a, b) => a.order - b.order);
          }
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:730", "获取轮播图数据失败", error);
      }
    },
    // 更新分享数据
    async updateShareData() {
      try {
        const settings = await utils_share.getShareSettings();
        if (settings && this.goodsData) {
          const shareTitle = this.goodsData.title || settings.content_share_text || "查看闲置商品";
          let shareImage = "";
          if (this.goodsData.images && this.goodsData.images.length > 0) {
            shareImage = this.goodsData.images[0];
          } else if (settings.content_share_image) {
            shareImage = settings.content_share_image;
          }
          this.goodsData.cover = shareImage;
          this.shareData = {
            title: shareTitle,
            imageUrl: shareImage,
            path: `/pages/community/detail?id=${this.goodsId}`
          };
          common_vendor.index.__f__("log", "at pages/community/detail.vue:765", "闲置商品分享数据已更新:", this.shareData);
        } else {
          this.initShareData();
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/community/detail.vue:771", "更新分享数据失败:", error);
        this.initShareData();
      }
    },
    // 添加onShow生命周期
    onShow() {
      this.updateShareData();
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
    f: !$data.loading
  }, !$data.loading ? common_vendor.e({
    g: common_vendor.f($data.goodsData.images, (img, index, i0) => {
      return {
        a: img,
        b: common_vendor.o(($event) => $options.previewImage(index), index),
        c: index
      };
    }),
    h: common_vendor.o((...args) => $options.handleSwiperChange && $options.handleSwiperChange(...args)),
    i: $data.goodsData.images && $data.goodsData.images.length > 0
  }, $data.goodsData.images && $data.goodsData.images.length > 0 ? {
    j: common_vendor.t($data.currentSwiperIndex + 1),
    k: common_vendor.t($data.goodsData.images.length)
  } : {}, {
    l: common_vendor.t($data.goodsData.price),
    m: $data.goodsData.originalPrice
  }, $data.goodsData.originalPrice ? {
    n: common_vendor.t($data.goodsData.originalPrice)
  } : {}, {
    o: common_vendor.t($data.goodsData.likes || 0),
    p: common_vendor.t($data.goodsData.viewCount),
    q: common_vendor.t($data.goodsData.title),
    r: $options.processedDescription,
    s: common_vendor.t($data.goodsData.detailLocation || $data.defaultLocation),
    t: common_vendor.o((...args) => $options.copyLocation && $options.copyLocation(...args)),
    v: common_vendor.t($data.goodsData.tradeMethod),
    w: $data.goodsData.publisher_avatar,
    x: common_vendor.t($data.goodsData.publisher_name),
    y: common_vendor.t($data.goodsData.publish_count),
    z: common_vendor.t($data.isFollowed ? "已关注" : "关注"),
    A: common_vendor.o((...args) => $options.handleFollow && $options.handleFollow(...args)),
    B: $data.bannerList.length > 0
  }, $data.bannerList.length > 0 ? common_vendor.e({
    C: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: item.image,
        b: index,
        c: common_vendor.o(($event) => $options.handleBannerClick(item), index)
      };
    }),
    D: common_vendor.o((...args) => $options.handleBannerSwiperChange && $options.handleBannerSwiperChange(...args)),
    E: $data.bannerList.length > 1
  }, $data.bannerList.length > 1 ? {
    F: common_vendor.f($data.bannerList, (item, index, i0) => {
      return {
        a: index,
        b: $data.currentBannerIndex === index ? 1 : ""
      };
    })
  } : {}) : {}, {
    G: common_vendor.t($data.goodsData.commentCount || 0),
    H: common_vendor.p({
      type: "chat",
      size: "14",
      color: "#007AFF"
    }),
    I: common_vendor.o((...args) => $options.handleComment && $options.handleComment(...args)),
    J: $data.goodsData.comments && $data.goodsData.comments.length
  }, $data.goodsData.comments && $data.goodsData.comments.length ? {
    K: common_vendor.f($data.goodsData.comments, (comment, index, i0) => {
      return {
        a: comment.avatar,
        b: common_vendor.t(comment.name),
        c: common_vendor.t(comment.content),
        d: common_vendor.t(comment.time),
        e: index
      };
    })
  } : {}) : {}, {
    L: _ctx.navigationBarHeight + "px",
    M: common_vendor.o($options.handleComment),
    N: common_vendor.o($options.handleCollect),
    O: common_vendor.o($options.handleShare),
    P: common_vendor.o($options.handleMessage),
    Q: common_vendor.p({
      ["is-collected"]: $data.isCollected,
      publisher: {
        id: $data.goodsData.publisher_id,
        name: $data.goodsData.publisher_name
      }
    }),
    R: common_vendor.s($options.pageStyle)
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/community/detail.js.map
