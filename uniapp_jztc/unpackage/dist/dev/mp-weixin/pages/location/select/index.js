"use strict";
const common_vendor = require("../../../common/vendor.js");
const utils_pinyin = require("../../../utils/pinyin.js");
const _sfc_main = {
  data() {
    return {
      currentLocation: "",
      searchKeyword: "",
      recentLocations: [],
      scrollTop: 0,
      letterPositions: {},
      showLetterIndicator: false,
      currentLetter: "",
      touchStartY: 0,
      touchStartTime: 0,
      isTouching: false,
      letterIndexRect: null
      // 字母索引容器的位置信息
    };
  },
  computed: {
    ...common_vendor.mapState("region", ["regionList", "loading", "error"]),
    ...common_vendor.mapGetters("region", ["hotRegions"]),
    // 过滤后的城市列表，分组显示
    filteredCityList() {
      const groupedList = this.$store.getters["region/groupedRegionList"];
      if (!this.searchKeyword) {
        return groupedList;
      }
      const keyword = this.searchKeyword.toLowerCase();
      const result = {};
      for (const letter in groupedList) {
        const filteredCities = groupedList[letter].filter((city) => {
          if (city.name.toLowerCase().includes(keyword)) {
            return true;
          }
          if (utils_pinyin.getStringFirstLetter(city.name).toLowerCase() === keyword.charAt(0).toLowerCase()) {
            return true;
          }
          const fullPinyin = utils_pinyin.getPinyin(city.name, { separator: "" }).toLowerCase();
          if (fullPinyin.includes(keyword)) {
            return true;
          }
          const pinyinInitials = utils_pinyin.getPinyin(city.name, { pattern: "first", separator: "" }).toLowerCase();
          if (pinyinInitials.includes(keyword)) {
            return true;
          }
          return false;
        });
        if (filteredCities.length > 0) {
          result[letter] = filteredCities;
        }
      }
      return result;
    },
    // 索引字母列表
    indexLetters() {
      return Object.keys(this.filteredCityList).sort();
    }
  },
  onLoad() {
    const savedLocation = common_vendor.index.getStorageSync("currentLocation");
    if (savedLocation) {
      this.currentLocation = savedLocation;
      common_vendor.index.__f__("log", "at pages/location/select/index.vue:206", "使用已保存的位置:", savedLocation);
    } else {
      this.currentLocation = "加载中...";
    }
    const recentLocations = common_vendor.index.getStorageSync("recentLocations");
    if (recentLocations) {
      this.recentLocations = JSON.parse(recentLocations);
    }
    this.fetchRegionList();
  },
  onReady() {
    setTimeout(() => {
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select(".letter-index").boundingClientRect((rect) => {
        this.letterIndexRect = rect;
      }).exec();
    }, 300);
  },
  methods: {
    ...common_vendor.mapActions("region", {
      getRegionList: "getRegionList"
    }),
    // 获取地区列表
    fetchRegionList() {
      this.getRegionList({ status: 0 });
      this.$nextTick(() => {
        let watchExecuted = false;
        const checkRegionList = () => {
          if (this.regionList && this.regionList.length > 0 && !watchExecuted) {
            watchExecuted = true;
            const savedLocation = common_vendor.index.getStorageSync("currentLocation");
            if (!savedLocation || savedLocation === "未知地区" || savedLocation === "加载中..." || savedLocation === "请选择地区") {
              if (this.regionList.length > 0) {
                this.currentLocation = this.regionList[0].name;
                common_vendor.index.setStorageSync("currentLocation", this.currentLocation);
                common_vendor.index.__f__("log", "at pages/location/select/index.vue:255", "位置选择页更新默认位置:", this.currentLocation);
              }
            }
            setTimeout(() => {
              this.getLetterPositions();
            }, 500);
          }
        };
        checkRegionList();
        if (!watchExecuted) {
          const timer = setInterval(() => {
            checkRegionList();
            if (watchExecuted || !this.regionList) {
              clearInterval(timer);
            }
          }, 500);
          setTimeout(() => {
            clearInterval(timer);
          }, 3e4);
        }
      });
    },
    handleSearch(e) {
    },
    // 选择城市
    selectLocation(cityName) {
      let selectedRegion = null;
      this.regionList.forEach((region) => {
        if (region.name === cityName) {
          selectedRegion = region;
        }
      });
      if (!selectedRegion) {
        common_vendor.index.__f__("error", "at pages/location/select/index.vue:304", "未找到选择的城市对象:", cityName);
        selectedRegion = {
          id: 0,
          name: cityName
        };
      }
      common_vendor.index.__f__("log", "at pages/location/select/index.vue:312", "选择城市:", selectedRegion);
      common_vendor.index.setStorageSync("currentLocation", selectedRegion.name);
      common_vendor.index.setStorageSync("currentLocationId", selectedRegion.id);
      common_vendor.index.__f__("log", "at pages/location/select/index.vue:319", "保存地区ID到本地存储:", selectedRegion.id);
      let recent = this.recentLocations.filter((item) => item !== selectedRegion.name);
      recent.unshift(selectedRegion.name);
      if (recent.length > 3) {
        recent = recent.slice(0, 3);
      }
      this.recentLocations = recent;
      common_vendor.index.setStorageSync("recentLocations", JSON.stringify(recent));
      const eventChannel = this.getOpenerEventChannel();
      if (eventChannel && eventChannel.emit) {
        eventChannel.emit("locationSelected", {
          id: selectedRegion.id,
          name: selectedRegion.name
        });
        common_vendor.index.__f__("log", "at pages/location/select/index.vue:337", "通过事件通道发送地区信息:", selectedRegion.id, selectedRegion.name);
      } else {
        common_vendor.index.__f__("error", "at pages/location/select/index.vue:339", "无法获取事件通道");
      }
      common_vendor.index.navigateBack();
      common_vendor.index.showToast({
        title: "已切换到" + selectedRegion.name,
        icon: "success"
      });
    },
    // 获取所有字母索引的位置信息
    getLetterPositions() {
      const query = common_vendor.index.createSelectorQuery().in(this);
      this.indexLetters.forEach((letter) => {
        query.select("#letter-" + letter).boundingClientRect();
      });
      query.select(".content-scroll").boundingClientRect();
      query.exec((res) => {
        if (res && res.length > 0) {
          const scrollViewInfo = res[res.length - 1];
          this.letterPositions = {};
          for (let i = 0; i < this.indexLetters.length; i++) {
            const letter = this.indexLetters[i];
            const position = res[i];
            if (position) {
              this.letterPositions[letter] = position.top - scrollViewInfo.top;
            }
          }
        }
      });
    },
    // 滚动到指定字母区域
    scrollToLetter(letter) {
      if (this.letterPositions && this.letterPositions[letter] !== void 0) {
        this.scrollTop = this.letterPositions[letter];
        if (common_vendor.index.vibrateShort) {
          common_vendor.index.vibrateShort();
        }
        return;
      }
      const query = common_vendor.index.createSelectorQuery().in(this);
      query.select("#letter-" + letter).boundingClientRect((data) => {
        if (data) {
          const scrollView = common_vendor.index.createSelectorQuery().in(this).select(".content-scroll");
          scrollView.boundingClientRect((scrollData) => {
            if (scrollData) {
              const scrollTop = data.top - scrollData.top;
              this.scrollTop = scrollTop;
              if (common_vendor.index.vibrateShort) {
                common_vendor.index.vibrateShort();
              }
            }
          }).exec();
        }
      }).exec();
    },
    // 处理触摸开始事件
    handleTouchStart(e) {
      this.touchStartY = e.touches[0].clientY;
      this.touchStartTime = Date.now();
      this.isTouching = true;
      const letter = this.getLetterFromTouch(e.touches[0].clientY);
      if (letter) {
        this.currentLetter = letter;
        this.showLetterIndicator = true;
        this.scrollToLetter(letter);
      }
    },
    // 处理触摸移动事件
    handleTouchMove(e) {
      if (!this.isTouching)
        return;
      const letter = this.getLetterFromTouch(e.touches[0].clientY);
      if (letter && letter !== this.currentLetter) {
        this.currentLetter = letter;
        this.scrollToLetter(letter);
      }
    },
    handleTouchEnd(e) {
      this.showLetterIndicator = false;
    },
    // 根据触摸位置获取对应的字母
    getLetterFromTouch(touchY) {
      if (!this.letterIndexRect)
        return null;
      const offsetY = touchY - this.letterIndexRect.top;
      const itemHeight = this.letterIndexRect.height / this.indexLetters.length;
      let index = Math.floor(offsetY / itemHeight);
      if (index < 0)
        index = 0;
      if (index >= this.indexLetters.length)
        index = this.indexLetters.length - 1;
      return this.indexLetters[index];
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
    a: common_vendor.p({
      type: "search",
      size: "16",
      color: "#666666"
    }),
    b: common_vendor.o([($event) => $data.searchKeyword = $event.detail.value, (...args) => $options.handleSearch && $options.handleSearch(...args)]),
    c: $data.searchKeyword,
    d: common_vendor.p({
      type: "location",
      size: "18",
      color: "#fc3e2b"
    }),
    e: common_vendor.t($data.currentLocation),
    f: common_vendor.o(($event) => $options.selectLocation($data.currentLocation)),
    g: $data.recentLocations.length > 0
  }, $data.recentLocations.length > 0 ? {
    h: common_vendor.f($data.recentLocations, (item, index, i0) => {
      return {
        a: common_vendor.t(item),
        b: index,
        c: common_vendor.o(($event) => $options.selectLocation(item), index)
      };
    })
  } : {}, {
    i: common_vendor.f(_ctx.hotRegions, (item, k0, i0) => {
      return {
        a: common_vendor.t(item.name),
        b: item.id,
        c: common_vendor.o(($event) => $options.selectLocation(item.name), item.id)
      };
    }),
    j: common_vendor.f($options.filteredCityList, (section, letter, i0) => {
      return {
        a: common_vendor.t(letter),
        b: "letter-" + letter,
        c: common_vendor.f(section, (city, k1, i1) => {
          return {
            a: common_vendor.t(city.name),
            b: city.id,
            c: common_vendor.o(($event) => $options.selectLocation(city.name), city.id)
          };
        }),
        d: letter
      };
    }),
    k: $data.scrollTop,
    l: common_vendor.f($options.indexLetters, (letter, k0, i0) => {
      return {
        a: common_vendor.t(letter),
        b: letter,
        c: letter,
        d: common_vendor.o(($event) => $options.scrollToLetter(letter), letter)
      };
    }),
    m: common_vendor.o((...args) => $options.handleTouchStart && $options.handleTouchStart(...args)),
    n: common_vendor.o((...args) => $options.handleTouchMove && $options.handleTouchMove(...args)),
    o: common_vendor.o((...args) => $options.handleTouchEnd && $options.handleTouchEnd(...args)),
    p: $data.showLetterIndicator
  }, $data.showLetterIndicator ? {
    q: common_vendor.t($data.currentLetter)
  } : {}, {
    r: _ctx.loading
  }, _ctx.loading ? {} : {}, {
    s: _ctx.error
  }, _ctx.error ? {
    t: common_vendor.t(_ctx.error),
    v: common_vendor.o((...args) => $options.fetchRegionList && $options.fetchRegionList(...args))
  } : {});
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../../.sourcemap/mp-weixin/pages/location/select/index.js.map
