"use strict";
const common_vendor = require("../../../common/vendor.js");
const apis_index = require("../../../apis/index.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      activeType: 0,
      // 0:全部, 1:普通信息, 2:闲置
      isRefreshing: false,
      loading: false,
      publishList: [],
      page: 1,
      pageSize: 10,
      total: 0,
      hasMore: true
    };
  },
  onLoad() {
    this.loadPublishList();
  },
  methods: {
    // 切换选项卡
    switchTab(type) {
      if (this.activeType === type)
        return;
      this.activeType = type;
      this.resetAndLoadList();
    },
    // 重置列表并重新加载
    resetAndLoadList() {
      this.publishList = [];
      this.page = 1;
      this.hasMore = true;
      this.loadPublishList();
    },
    // 加载发布列表
    async loadPublishList() {
      if (!this.hasMore || this.loading)
        return;
      this.loading = true;
      try {
        const params = {
          page: this.page,
          pageSize: this.pageSize,
          type: this.activeType
        };
        const res = await apis_index.publish.getMyPublishList(params);
        if (res.code === 0 && res.data) {
          const { list, total, page, pages } = res.data;
          this.publishList = this.page === 1 ? list : [...this.publishList, ...list];
          this.total = total;
          this.hasMore = page < pages;
          if (this.hasMore) {
            this.page++;
          }
        } else {
          common_vendor.index.showToast({
            title: res.message || "加载失败",
            icon: "none"
          });
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/publish/index.vue:160", "加载发布列表失败", error);
        common_vendor.index.showToast({
          title: "加载失败",
          icon: "none"
        });
      } finally {
        this.loading = false;
      }
    },
    // 下拉刷新
    async onRefresh() {
      this.isRefreshing = true;
      this.resetAndLoadList();
      setTimeout(() => {
        this.isRefreshing = false;
      }, 1e3);
    },
    // 上拉加载更多
    loadMore() {
      this.loadPublishList();
    },
    // 获取状态样式
    getStatusClass(status) {
      switch (status) {
        case "已发布":
          return "published";
        case "已下架":
          return "offline";
        case "已售出":
          return "sold";
        case "审核中":
          return "pending";
        default:
          return "default";
      }
    },
    // 点击内容项
    handleItemClick(item) {
      common_vendor.index.navigateTo({
        url: `/pages/content/detail?id=${item.id}`
      });
    },
    // 编辑
    handleEdit(item) {
      common_vendor.index.navigateTo({
        url: `/pages/publish/edit?id=${item.id}`
      });
    },
    // 删除
    handleDelete(item, index) {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要删除该发布吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              const result = await apis_index.publish.deleteContent(item.id);
              if (result.code === 0) {
                this.publishList.splice(index, 1);
                common_vendor.index.showToast({
                  title: "删除成功",
                  icon: "success"
                });
              } else {
                common_vendor.index.showToast({
                  title: result.message || "删除失败",
                  icon: "none"
                });
              }
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/my/publish/index.vue:232", "删除失败", error);
              common_vendor.index.showToast({
                title: "删除失败",
                icon: "none"
              });
            }
          }
        }
      });
    },
    // 去发布
    handlePublish() {
      common_vendor.index.switchTab({
        url: "/pages/publish/index"
      });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.activeType === 0 ? 1 : "",
    b: common_vendor.o(($event) => $options.switchTab(0)),
    c: $data.activeType === 1 ? 1 : "",
    d: common_vendor.o(($event) => $options.switchTab(1)),
    e: $data.activeType === 2 ? 1 : "",
    f: common_vendor.o(($event) => $options.switchTab(2)),
    g: $data.publishList.length > 0
  }, $data.publishList.length > 0 ? {
    h: common_vendor.f($data.publishList, (item, index, i0) => {
      return {
        a: common_vendor.t(item.status),
        b: common_vendor.n($options.getStatusClass(item.status)),
        c: common_vendor.t(item.category),
        d: common_vendor.t(item.title),
        e: common_vendor.t(item.publishedAt),
        f: common_vendor.o(($event) => $options.handleEdit(item), item.id),
        g: common_vendor.o(($event) => $options.handleDelete(item, index), item.id),
        h: item.id,
        i: common_vendor.o(($event) => $options.handleItemClick(item), item.id)
      };
    })
  } : {}, {
    i: $data.publishList.length > 0 && $data.hasMore
  }, $data.publishList.length > 0 && $data.hasMore ? {} : {}, {
    j: $data.publishList.length > 0 && !$data.hasMore
  }, $data.publishList.length > 0 && !$data.hasMore ? {} : {}, {
    k: $data.publishList.length === 0 && !$data.loading
  }, $data.publishList.length === 0 && !$data.loading ? {
    l: common_assets._imports_0$1,
    m: common_vendor.o((...args) => $options.handlePublish && $options.handlePublish(...args))
  } : {}, {
    n: $data.isRefreshing,
    o: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    p: common_vendor.o((...args) => $options.loadMore && $options.loadMore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/publish/index.js.map
