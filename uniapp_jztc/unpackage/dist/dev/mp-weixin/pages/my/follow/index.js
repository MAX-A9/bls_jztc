"use strict";
const common_vendor = require("../../../common/vendor.js");
const apis_content = require("../../../apis/content.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      isRefreshing: false,
      loading: false,
      userList: [],
      currentPage: 1,
      pageSize: 10,
      hasMore: true,
      total: 0
    };
  },
  onLoad() {
    this.loadFollowingList();
  },
  methods: {
    // 加载关注列表
    async loadFollowingList() {
      if (this.loading)
        return;
      this.loading = true;
      try {
        const result = await apis_content.getFollowingList({
          page: this.currentPage,
          size: this.pageSize
        });
        if (result.code === 0 && result.data) {
          if (this.currentPage === 1) {
            this.userList = result.data.list || [];
          } else {
            this.userList = [...this.userList, ...result.data.list || []];
          }
          this.total = result.data.total || 0;
          this.hasMore = this.userList.length < this.total;
        } else {
          common_vendor.index.showToast({
            title: result.message || "加载失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/follow/index.vue:107", "加载关注列表失败:", error);
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
      this.currentPage = 1;
      await this.loadFollowingList();
      if (this.isRefreshing) {
        this.isRefreshing = false;
      }
    },
    // 加载更多
    loadMore() {
      if (this.loading || !this.hasMore)
        return;
      this.currentPage++;
      this.loadFollowingList();
    },
    // 点击用户进入详情
    handleUserClick(user) {
      if (!user || !user.client_id)
        return;
      common_vendor.index.navigateTo({
        url: `/pages/user/detail?id=${user.client_id}`
      });
    },
    // 取消关注
    async toggleFollow(user, index) {
      if (!user || !user.client_id)
        return;
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要取消关注吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              common_vendor.index.showLoading({
                title: "处理中...",
                mask: true
              });
              const result = await apis_content.unfollowUser({
                publisher_id: user.client_id
              });
              common_vendor.index.hideLoading();
              if (result && result.code === 0) {
                this.userList.splice(index, 1);
                this.total--;
                common_vendor.index.showToast({
                  title: "已取消关注",
                  icon: "success"
                });
              } else {
                throw new Error((result == null ? void 0 : result.message) || "操作失败");
              }
            } catch (error) {
              common_vendor.index.hideLoading();
              common_vendor.index.__f__("error", "at pages/my/follow/index.vue:187", "取消关注失败:", error);
              common_vendor.index.showToast({
                title: error.message || "操作失败，请重试",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 跳转到发现页
    handleExplore() {
      common_vendor.index.switchTab({
        url: "/pages/index/index"
      });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.loading && $data.userList.length === 0
  }, $data.loading && $data.userList.length === 0 ? {} : {}, {
    b: $data.userList.length > 0
  }, $data.userList.length > 0 ? common_vendor.e({
    c: common_vendor.f($data.userList, (item, index, i0) => {
      return {
        a: item.avatar_url || "/static/avatar/default.png",
        b: common_vendor.t(item.real_name),
        c: common_vendor.t(item.publish_count || 0),
        d: common_vendor.o(($event) => $options.toggleFollow(item, index), item.client_id),
        e: item.client_id,
        f: common_vendor.o(($event) => $options.handleUserClick(item), item.client_id)
      };
    }),
    d: $data.hasMore && !$data.loading
  }, $data.hasMore && !$data.loading ? {} : {}, {
    e: $data.loading && $data.userList.length > 0
  }, $data.loading && $data.userList.length > 0 ? {} : {}, {
    f: !$data.hasMore && $data.userList.length > 0
  }, !$data.hasMore && $data.userList.length > 0 ? {} : {}) : {}, {
    g: !$data.loading && $data.userList.length === 0
  }, !$data.loading && $data.userList.length === 0 ? {
    h: common_assets._imports_0$1,
    i: common_vendor.o((...args) => $options.handleExplore && $options.handleExplore(...args))
  } : {}, {
    j: $data.isRefreshing,
    k: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    l: common_vendor.o((...args) => $options.loadMore && $options.loadMore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/follow/index.js.map
