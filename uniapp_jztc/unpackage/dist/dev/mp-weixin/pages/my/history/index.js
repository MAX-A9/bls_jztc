"use strict";
const common_vendor = require("../../../common/vendor.js");
const apis_content = require("../../../apis/content.js");
const common_assets = require("../../../common/assets.js");
const _sfc_main = {
  data() {
    return {
      currentFilter: 0,
      isRefreshing: false,
      dateFilters: ["全部", "今天", "昨天", "本周", "本月"],
      timeTypeMap: {
        0: "all",
        1: "today",
        2: "yesterday",
        3: "this_week",
        4: "this_month"
      },
      historyList: [],
      page: 1,
      pageSize: 10,
      hasMore: true,
      loading: false,
      total: 0
    };
  },
  created() {
    this.loadBrowseHistory();
  },
  computed: {
    groupedHistory() {
      const grouped = {};
      this.historyList.forEach((item) => {
        const date = item.browseTime.split(" ")[0];
        if (!grouped[date]) {
          grouped[date] = [];
        }
        grouped[date].push(item);
      });
      return grouped;
    }
  },
  methods: {
    switchFilter(index) {
      if (this.currentFilter === index)
        return;
      this.currentFilter = index;
      this.page = 1;
      this.historyList = [];
      this.hasMore = true;
      this.loadBrowseHistory();
    },
    // 加载浏览历史数据
    async loadBrowseHistory() {
      if (this.loading)
        return;
      try {
        this.loading = true;
        const params = {
          timeType: this.timeTypeMap[this.currentFilter],
          page: this.page,
          pageSize: this.pageSize
        };
        const result = await apis_content.getBrowseHistoryList(params);
        if (result.code === 0 && result.data) {
          const { list = [], total = 0, page = 1 } = result.data;
          const formattedList = list.map((item) => {
            return {
              id: item.id,
              contentId: item.contentId,
              title: item.contentTitle || "未知标题",
              image: item.contentCover || "",
              price: item.price || 0,
              contentType: item.contentType,
              browseTime: item.browseTime,
              time: item.browseTime.split(" ")[1],
              category: item.category || "未分类",
              contentStatus: item.contentStatus
            };
          });
          if (this.page === 1) {
            this.historyList = formattedList;
          } else {
            this.historyList = [...this.historyList, ...formattedList];
          }
          this.total = total;
          this.hasMore = this.historyList.length < total;
          this.page++;
        } else {
          throw new Error(result.message || "获取浏览历史失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at pages/my/history/index.vue:203", "加载浏览历史失败:", error);
        common_vendor.index.showToast({
          title: "加载失败，请重试",
          icon: "none",
          duration: 2e3
        });
      } finally {
        this.loading = false;
        if (this.isRefreshing) {
          this.isRefreshing = false;
        }
      }
    },
    async onRefresh() {
      this.isRefreshing = true;
      this.page = 1;
      this.hasMore = true;
      await this.loadBrowseHistory();
      common_vendor.index.showToast({
        title: "刷新成功",
        icon: "success"
      });
    },
    handleItemClick(item) {
      if (item.contentStatus === 0 || item.contentStatus === 2) {
        common_vendor.index.showToast({
          title: "该内容已下架",
          icon: "none",
          duration: 2e3
        });
        return;
      }
      let url = "";
      switch (item.contentType) {
        case "idle":
          url = `/pages/community/detail?id=${item.contentId}`;
          break;
        case "article":
        case "info":
        default:
          url = `/pages/content/detail?id=${item.contentId}`;
          break;
      }
      common_vendor.index.navigateTo({ url });
    },
    async handleClearHistory() {
      common_vendor.index.showModal({
        title: "提示",
        content: "确定要清空浏览记录吗？",
        success: async (res) => {
          if (res.confirm) {
            try {
              const timeType = this.timeTypeMap[this.currentFilter];
              const result = await apis_content.clearBrowseHistory(timeType);
              if (result.code === 0) {
                this.historyList = [];
                this.page = 1;
                this.hasMore = false;
                common_vendor.index.showToast({
                  title: "已清空记录",
                  icon: "success"
                });
              } else {
                throw new Error(result.message || "清空失败");
              }
            } catch (error) {
              common_vendor.index.__f__("error", "at pages/my/history/index.vue:279", "清空浏览记录失败:", error);
              common_vendor.index.showToast({
                title: "清空失败，请重试",
                icon: "none"
              });
            }
          }
        }
      });
    },
    handleExplore() {
      common_vendor.index.switchTab({
        url: "/pages/index/index"
      });
    },
    // 加载更多
    loadMore() {
      if (this.hasMore && !this.loading) {
        this.loadBrowseHistory();
      }
    },
    handleImageError(item) {
      item.image = "";
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
    a: common_vendor.f($data.dateFilters, (item, index, i0) => {
      return {
        a: common_vendor.t(item),
        b: index,
        c: $data.currentFilter === index ? 1 : "",
        d: common_vendor.o(($event) => $options.switchFilter(index), index)
      };
    }),
    b: common_vendor.p({
      type: "trash",
      size: "16",
      color: "#666666"
    }),
    c: common_vendor.t($data.currentFilter > 0 ? $data.dateFilters[$data.currentFilter] : "全部"),
    d: common_vendor.o((...args) => $options.handleClearHistory && $options.handleClearHistory(...args)),
    e: $data.loading && $data.historyList.length === 0
  }, $data.loading && $data.historyList.length === 0 ? {} : {}, {
    f: $data.historyList.length > 0
  }, $data.historyList.length > 0 ? common_vendor.e({
    g: common_vendor.f($options.groupedHistory, (group, date, i0) => {
      return {
        a: common_vendor.t(date),
        b: common_vendor.f(group, (item, index, i1) => {
          return common_vendor.e({
            a: item.contentStatus === 0 || item.contentStatus === 2
          }, item.contentStatus === 0 || item.contentStatus === 2 ? {} : {}, {
            b: item.image
          }, item.image ? {
            c: item.image,
            d: common_vendor.o(($event) => $options.handleImageError(item), index)
          } : {}, {
            e: common_vendor.t(item.title),
            f: item.contentStatus === 0 || item.contentStatus === 2 ? 1 : "",
            g: item.price > 0
          }, item.price > 0 ? {
            h: common_vendor.t(item.price)
          } : {}, {
            i: common_vendor.t(item.time),
            j: common_vendor.t(item.category),
            k: item.contentStatus === 0 || item.contentStatus === 2
          }, item.contentStatus === 0 || item.contentStatus === 2 ? {} : {}, {
            l: item.contentStatus === 0 || item.contentStatus === 2 ? 1 : "",
            m: !item.image ? 1 : "",
            n: index,
            o: common_vendor.o(($event) => $options.handleItemClick(item), index),
            p: item.contentStatus === 0 || item.contentStatus === 2 ? 1 : ""
          });
        }),
        c: date
      };
    }),
    h: $data.historyList.length > 0
  }, $data.historyList.length > 0 ? common_vendor.e({
    i: $data.loading && $data.historyList.length > 0
  }, $data.loading && $data.historyList.length > 0 ? {} : $data.hasMore ? {} : {}, {
    j: $data.hasMore
  }) : {}) : {}, {
    k: !$data.loading && $data.historyList.length === 0
  }, !$data.loading && $data.historyList.length === 0 ? {
    l: common_assets._imports_0$3,
    m: common_vendor.o((...args) => $options.handleExplore && $options.handleExplore(...args))
  } : {}, {
    n: $data.isRefreshing,
    o: common_vendor.o((...args) => $options.onRefresh && $options.onRefresh(...args)),
    p: common_vendor.o((...args) => $options.loadMore && $options.loadMore(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/my/history/index.js.map
