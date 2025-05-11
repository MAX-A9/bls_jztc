"use strict";
const common_vendor = require("../../../common/vendor.js");
const apis_content = require("../../../apis/content.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      currentTab: 0,
      tabList: ["闲置", "帖子"],
      isRefreshing: false,
      loading: false,
      favoriteList: [],
      postList: [],
      favoriteParams: {
        page: 1,
        pageSize: 10,
        type: 2
      },
      postParams: {
        page: 1,
        pageSize: 10,
        type: 1
      },
      hasMore: true,
      favoriteTotal: 0,
      postTotal: 0
    };
  },
  onLoad() {
    this.loadFavoriteData();
  },
  methods: {
    // 切换选项卡
    switchTab(index) {
      if (this.currentTab === index)
        return;
      this.currentTab = index;
      if (index === 0 && this.favoriteList.length === 0) {
        this.loadFavoriteData();
      } else if (index === 1 && this.postList.length === 0) {
        this.loadPostData();
      }
    },
    // 加载闲置收藏数据
    async loadFavoriteData() {
      if (this.loading)
        return;
      this.loading = true;
      try {
        const result = await apis_content.getFavoriteList({
          page: this.favoriteParams.page,
          pageSize: this.favoriteParams.pageSize,
          type: this.favoriteParams.type
        });
        if (result.code === 0 && result.data) {
          if (this.favoriteParams.page === 1) {
            this.favoriteList = result.data.list || [];
          } else {
            this.favoriteList = [...this.favoriteList, ...result.data.list || []];
          }
          this.favoriteTotal = result.data.total || 0;
          this.hasMore = this.favoriteList.length < this.favoriteTotal;
        } else {
          common_vendor.index.showToast({
            title: result.message || "加载失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.showToast({
          title: "网络异常，请稍后重试",
          icon: "none"
        });
      } finally {
        this.loading = false;
        if (this.isRefreshing) {
          this.isRefreshing = false;
        }
      }
    },
    // 加载帖子收藏数据
    async loadPostData() {
      if (this.loading)
        return;
      this.loading = true;
      try {
        const result = await apis_content.getFavoriteList({
          page: this.postParams.page,
          pageSize: this.postParams.pageSize,
          type: this.postParams.type
        });
        if (result.code === 0 && result.data) {
          if (this.postParams.page === 1) {
            this.postList = result.data.list || [];
          } else {
            this.postList = [...this.postList, ...result.data.list || []];
          }
          this.postTotal = result.data.total || 0;
          this.hasMore = this.postList.length < this.postTotal;
        } else {
          common_vendor.index.showToast({
            title: result.message || "加载失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.showToast({
          title: "网络异常，请稍后重试",
          icon: "none"
        });
      } finally {
        this.loading = false;
        if (this.isRefreshing) {
          this.isRefreshing = false;
        }
      }
    },
    // 下拉刷新
    async onRefresh() {
      this.isRefreshing = true;
      if (this.currentTab === 0) {
        this.favoriteParams.page = 1;
        await this.loadFavoriteData();
      } else {
        this.postParams.page = 1;
        await this.loadPostData();
      }
    },
    // 加载更多
    loadMore() {
      if (this.loading || !this.hasMore)
        return;
      if (this.currentTab === 0) {
        this.favoriteParams.page++;
        this.loadFavoriteData();
      } else {
        this.postParams.page++;
        this.loadPostData();
      }
    },
    // 点击收藏项
    handleItemClick(item) {
      if (this.currentTab === 0 || item.type === 2) {
        common_vendor.index.navigateTo({
          url: `/pages/community/detail?id=${item.id}`
        });
      } else {
        common_vendor.index.navigateTo({
          url: `/pages/content/detail?id=${item.id}`
        });
      }
    },
    // 取消收藏
    async toggleFavorite(contentId, type, index) {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要取消收藏吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              const result = await apis_content.cancelFavorite(contentId);
              if (result.code === 0) {
                if (type === 1) {
                  this.postList.splice(index, 1);
                } else {
                  this.favoriteList.splice(index, 1);
                }
                common_vendor.index.showToast({
                  title: "已取消收藏",
                  icon: "success"
                });
              } else {
                common_vendor.index.showToast({
                  title: result.message || "操作失败",
                  icon: "none"
                });
              }
            } catch (error) {
              common_vendor.index.showToast({
                title: "网络异常，请稍后重试",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 去发现
    handleExplore() {
      common_vendor.index.switchTab({
        url: "/pages/index/index"
      });
    }
  }
};
if (!Array) {
  const _easycom_uni_icons2 = common_vendor.resolveComponent("uni-icons");
  _easycom_uni_icons2();
}
const _easycom_uni_icons = () => "../../../node-modules/@dcloudio/uni-ui/lib/uni-icons/uni-icons.js";
if (!Math) {
  _easycom_uni_icons();
}
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: common_vendor.f($data.tabList, (item, index, i0) => {
      return {
        a: common_vendor.t(item),
        b: index,
        c: $data.currentTab === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchTab(index), index)
      };
    }),
    b: $data.currentTab === 0 && $data.favoriteList.length > 0
  }, $data.currentTab === 0 && $data.favoriteList.length > 0 ? common_vendor.e({
    c: common_vendor.f($data.favoriteList, (item, index, i0) => {
      return common_vendor.e({
        a: item.image,
        b: common_vendor.t(item.title),
        c: item.price > 0
      }, item.price > 0 ? {
        d: common_vendor.t(item.price)
      } : {}, {
        e: common_vendor.t(item.publisher),
        f: "0496ab34-0-" + i0,
        g: common_vendor.o(($event) => $options.toggleFavorite(item.id, 2, index), item.id),
        h: item.id,
        i: common_vendor.o(($event) => $options.handleItemClick(item), item.id)
      });
    }),
    d: common_vendor.p({
      type: "star-filled",
      size: "20",
      color: "#fc3e2b"
    }),
    e: $data.hasMore && !$data.loading
  }, $data.hasMore && !$data.loading ? {} : {}, {
    f: $data.loading
  }, $data.loading ? {} : {}, {
    g: !$data.hasMore && $data.favoriteList.length > 0
  }, !$data.hasMore && $data.favoriteList.length > 0 ? {} : {}) : {}, {
    h: $data.currentTab === 1 && $data.postList.length > 0
  }, $data.currentTab === 1 && $data.postList.length > 0 ? common_vendor.e({
    i: common_vendor.f($data.postList, (item, index, i0) => {
      return common_vendor.e({
        a: item.image
      }, item.image ? {
        b: item.image
      } : {}, {
        c: common_vendor.t(item.title),
        d: common_vendor.t(item.publisher),
        e: "0496ab34-1-" + i0,
        f: common_vendor.o(($event) => $options.toggleFavorite(item.id, 1, index), item.id),
        g: !item.image ? 1 : "",
        h: item.id,
        i: common_vendor.o(($event) => $options.handleItemClick(item), item.id)
      });
    }),
    j: common_vendor.p({
      type: "star-filled",
      size: "20",
      color: "#fc3e2b"
    }),
    k: $data.hasMore && !$data.loading
  }, $data.hasMore && !$data.loading ? {} : {}, {
    l: $data.loading
  }, $data.loading ? {} : {}, {
    m: !$data.hasMore && $data.postList.length > 0
  }, !$data.hasMore && $data.postList.length > 0 ? {} : {}) : {}, {
    n: $data.currentTab === 0 && $data.favoriteList.length === 0 && !$data.loading || $data.currentTab === 1 && $data.postList.length === 0 && !$data.loading
  }, $data.currentTab === 0 && $data.favoriteList.length === 0 && !$data.loading || $data.currentTab === 1 && $data.postList.length === 0 && !$data.loading ? {
    o: common_assets._imports_0$1,
    p: common_vendor.o((...args) => $options.handleExplore && $options.handleExplore(...args))
  } : {}, {
    q: $data.isRefreshing,
    r: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    s: common_vendor.o((...args) => $options.loadMore && $options.loadMore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/favorite/index.js.map
