"use strict";
const common_vendor = require("../../common/vendor.js");
const apis_content = require("../../apis/content.js");
const utils_date = require("../../utils/date.js");
const common_assets = require("../../common/assets.js");
const _sfc_main = {
  name: "ContentList",
  props: {
    // 可选参数，允许从父组件传入
    categoryId: {
      type: [Number, String],
      default: ""
    }
  },
  data() {
    return {
      contentList: [],
      loading: false,
      hasMore: true,
      page: 1,
      pageSize: 10,
      regionId: 0,
      // 当前区域ID
      currentCategory: "",
      // 当前分类
      isInitialized: false
      // 添加初始化标志
    };
  },
  created() {
    const currentLocationId = common_vendor.index.getStorageSync("currentLocationId") || 0;
    if (currentLocationId) {
      this.regionId = currentLocationId;
      common_vendor.index.__f__("log", "at components/content-list/index.vue:86", "组件创建时，从本地存储初始化区域ID:", this.regionId);
    }
    common_vendor.index.$on("locationChanged", this.handleLocationChanged);
    common_vendor.index.__f__("log", "at components/content-list/index.vue:91", "已注册locationChanged事件监听");
  },
  mounted() {
    this.$nextTick(() => {
      this.init();
    });
  },
  beforeDestroy() {
    common_vendor.index.$off("locationChanged", this.handleLocationChanged);
  },
  methods: {
    // 初始化方法，用于组件挂载后或页面显示时调用
    init() {
      if (this.isInitialized && this.contentList.length > 0) {
        common_vendor.index.__f__("log", "at components/content-list/index.vue:108", "内容列表已初始化，跳过重复初始化");
        return;
      }
      const currentLocationId = common_vendor.index.getStorageSync("currentLocationId") || 0;
      if ((!this.regionId || this.regionId === 0) && currentLocationId) {
        this.regionId = currentLocationId;
        common_vendor.index.__f__("log", "at components/content-list/index.vue:118", "初始化时，从本地存储更新区域ID:", this.regionId);
      }
      if (this.regionId && this.regionId !== 0) {
        common_vendor.index.__f__("log", "at components/content-list/index.vue:123", "初始化时，开始加载内容列表数据...");
        this.resetAndLoad();
        this.isInitialized = true;
      } else {
        common_vendor.index.__f__("log", "at components/content-list/index.vue:128", "初始化时，未找到有效的区域ID");
      }
    },
    // 格式化发布时间
    formatPublishTime(time) {
      return utils_date.formatTimeAgo(time);
    },
    handleItemClick(item) {
      common_vendor.index.navigateTo({
        url: `/pages/content/detail?id=${item.id}&category=${item.category}`,
        success: () => {
          common_vendor.index.__f__("log", "at components/content-list/index.vue:141", "跳转到详情页:", item);
        },
        fail: (err) => {
          common_vendor.index.__f__("error", "at components/content-list/index.vue:144", "跳转失败:", err);
        }
      });
    },
    // 处理分类切换
    handleCategoryChange(categoryId) {
      this.currentCategory = categoryId;
      this.resetAndLoad();
    },
    // 处理位置变更
    handleLocationChanged(data) {
      common_vendor.index.__f__("log", "at components/content-list/index.vue:158", "位置变更事件触发，接收到的regionId:", data.regionId);
      const newRegionId = data.regionId || 0;
      if (!newRegionId) {
        common_vendor.index.__f__("error", "at components/content-list/index.vue:163", "接收到无效的regionId:", newRegionId);
        return;
      }
      if (this.regionId !== newRegionId) {
        common_vendor.index.__f__("log", "at components/content-list/index.vue:169", "地区ID已更新:", this.regionId, "->", newRegionId);
        this.regionId = newRegionId;
        this.resetAndLoad();
      } else {
        common_vendor.index.__f__("log", "at components/content-list/index.vue:174", "地区ID未变化，保持当前值:", this.regionId);
      }
    },
    // 重置数据并重新加载
    resetAndLoad() {
      common_vendor.index.__f__("log", "at components/content-list/index.vue:180", "重置内容列表，当前区域ID:", this.regionId);
      this.page = 1;
      this.contentList = [];
      this.hasMore = true;
      this.loadData();
    },
    // 加载更多数据
    loadMore() {
      if (this.loading || !this.hasMore)
        return;
      this.page++;
      this.loadData();
    },
    // 加载数据
    async loadData() {
      if (this.loading)
        return;
      try {
        this.loading = true;
        const currentLocationId = common_vendor.index.getStorageSync("currentLocationId") || 0;
        common_vendor.index.__f__("log", "at components/content-list/index.vue:203", "从本地存储获取的区域ID:", currentLocationId);
        if (!this.regionId && currentLocationId) {
          common_vendor.index.__f__("log", "at components/content-list/index.vue:207", "使用本地存储区域ID更新组件regionId:", currentLocationId);
          this.regionId = currentLocationId;
        }
        if (!this.regionId) {
          common_vendor.index.showToast({
            title: "请先选择地区",
            icon: "none"
          });
          throw new Error("未选择地区");
        }
        const params = {
          regionId: this.regionId,
          page: this.page,
          pageSize: this.pageSize
        };
        if (this.currentCategory) {
          params.category = this.currentCategory;
        }
        common_vendor.index.__f__("log", "at components/content-list/index.vue:232", "加载内容列表，参数:", params);
        const result = await apis_content.getRegionContentList(params);
        if (result.code === 0 && result.data) {
          const { list = [], total = 0 } = result.data;
          let formattedList = list.map((item) => {
            return {
              ...item,
              // 格式化发布时间为"多久前"的格式
              publishTimeFormatted: this.formatPublishTime(item.publishTime)
            };
          });
          formattedList.sort((a, b) => {
            if (a.isTop && !b.isTop)
              return -1;
            if (!a.isTop && b.isTop)
              return 1;
            let timeA = 0, timeB = 0;
            try {
              if (a.publishTime && typeof a.publishTime === "string" && a.publishTime.includes(" ")) {
                const aDateParts = a.publishTime.split(" ")[0].split("-");
                const aTimeParts = a.publishTime.split(" ")[1].split(":");
                const aDate = new Date(
                  parseInt(aDateParts[0]),
                  // 年
                  parseInt(aDateParts[1]) - 1,
                  // 月(0-11)
                  parseInt(aDateParts[2]),
                  // 日
                  parseInt(aTimeParts[0]),
                  // 时
                  parseInt(aTimeParts[1]),
                  // 分
                  parseInt(aTimeParts[2])
                  // 秒
                );
                timeA = aDate.getTime();
              } else {
                timeA = utils_date.getTimestamp(a.publishTime);
              }
              if (b.publishTime && typeof b.publishTime === "string" && b.publishTime.includes(" ")) {
                const bDateParts = b.publishTime.split(" ")[0].split("-");
                const bTimeParts = b.publishTime.split(" ")[1].split(":");
                const bDate = new Date(
                  parseInt(bDateParts[0]),
                  // 年
                  parseInt(bDateParts[1]) - 1,
                  // 月(0-11)
                  parseInt(bDateParts[2]),
                  // 日
                  parseInt(bTimeParts[0]),
                  // 时
                  parseInt(bTimeParts[1]),
                  // 分
                  parseInt(bTimeParts[2])
                  // 秒
                );
                timeB = bDate.getTime();
              } else {
                timeB = utils_date.getTimestamp(b.publishTime);
              }
            } catch (e) {
              common_vendor.index.__f__("error", "at components/content-list/index.vue:297", "日期解析错误:", e, a.publishTime, b.publishTime);
              timeA = timeB = Date.now();
            }
            return timeB - timeA;
          });
          if (this.page === 1) {
            this.contentList = formattedList;
          } else {
            this.contentList = [...this.contentList, ...formattedList];
          }
          this.hasMore = this.contentList.length < total;
          common_vendor.index.__f__("log", "at components/content-list/index.vue:316", `加载了${list.length}条数据，总共${total}条，使用的区域ID:${this.regionId}`);
          this.$parent && this.$parent.refreshing === true && this.notifyRefreshComplete();
        } else {
          throw new Error(result.message || "获取内容列表失败");
        }
      } catch (error) {
        common_vendor.index.__f__("error", "at components/content-list/index.vue:324", "加载内容列表失败:", error);
        common_vendor.index.showToast({
          title: "加载失败，请重试",
          icon: "none",
          duration: 2e3
        });
      } finally {
        this.loading = false;
      }
    },
    // 通知父组件刷新完成
    notifyRefreshComplete() {
      this.$emit("refreshComplete");
      if (this.$parent && this.$parent.refreshing === true) {
        this.$parent.refreshing = false;
        common_vendor.index.__f__("log", "at components/content-list/index.vue:343", "内容加载完成，关闭刷新状态");
      }
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: $data.loading && $data.contentList.length === 0
  }, $data.loading && $data.contentList.length === 0 ? {} : {}, {
    b: common_vendor.f($data.contentList, (item, index, i0) => {
      return common_vendor.e({
        a: item.isTop
      }, item.isTop ? {} : {}, {
        b: common_vendor.t(item.category),
        c: common_vendor.t(item.title),
        d: common_vendor.t(item.publisher),
        e: common_vendor.t(item.publishTimeFormatted),
        f: common_vendor.n({
          "full-width": !item.image
        }),
        g: item.image
      }, item.image ? {
        h: item.image
      } : {}, {
        i: index,
        j: common_vendor.o(($event) => $options.handleItemClick(item), index)
      });
    }),
    c: $data.contentList.length > 0
  }, $data.contentList.length > 0 ? common_vendor.e({
    d: $data.hasMore && !$data.loading
  }, $data.hasMore && !$data.loading ? {} : {}, {
    e: $data.loading && $data.contentList.length > 0
  }, $data.loading && $data.contentList.length > 0 ? {} : {}, {
    f: !$data.hasMore && $data.contentList.length > 0
  }, !$data.hasMore && $data.contentList.length > 0 ? {} : {}) : {}, {
    g: !$data.loading && $data.contentList.length === 0
  }, !$data.loading && $data.contentList.length === 0 ? {
    h: common_assets._imports_0$1
  } : {});
}
const Component = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createComponent(Component);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/components/content-list/index.js.map
