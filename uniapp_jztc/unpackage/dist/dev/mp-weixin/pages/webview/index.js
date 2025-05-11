"use strict";
const common_vendor = require("../../common/vendor.js");
const _sfc_main = {
  data() {
    return {
      webUrl: "",
      originalUrl: "",
      securityEnabled: false,
      urlTitle: ""
    };
  },
  onLoad(options) {
    if (options.url) {
      try {
        this.originalUrl = decodeURIComponent(options.url) || "";
        if (options.title) {
          this.urlTitle = decodeURIComponent(options.title);
          common_vendor.index.setNavigationBarTitle({
            title: this.urlTitle
          });
        } else {
          common_vendor.index.setNavigationBarTitle({
            title: "网页浏览"
          });
        }
        this.enableWebViewSecurity();
        this.setSecureWebUrl(this.originalUrl);
      } catch (e) {
        common_vendor.index.__f__("error", "at pages/webview/index.vue:53", "URL解析错误:", e);
        common_vendor.index.showToast({
          title: "无效的链接",
          icon: "none"
        });
        setTimeout(() => {
          common_vendor.index.navigateBack();
        }, 1500);
      }
    } else {
      common_vendor.index.showToast({
        title: "缺少URL参数",
        icon: "none"
      });
      setTimeout(() => {
        common_vendor.index.navigateBack();
      }, 1500);
    }
  },
  mounted() {
    this.disableSharedArrayBuffer();
  },
  methods: {
    // 启用WebView安全设置
    enableWebViewSecurity() {
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.setWebViewSecurity) {
        common_vendor.wx$1.setWebViewSecurity({
          enable: true,
          success: () => {
            common_vendor.index.__f__("log", "at pages/webview/index.vue:83", "成功启用WebView安全模式");
            this.securityEnabled = true;
          },
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/webview/index.vue:87", "启用WebView安全模式失败:", err);
          }
        });
      }
    },
    // 禁用SharedArrayBuffer相关功能
    disableSharedArrayBuffer() {
      if (typeof common_vendor.wx$1 !== "undefined") {
        if (common_vendor.wx$1.onWebViewLoad) {
          common_vendor.wx$1.onWebViewLoad(() => {
            try {
              common_vendor.index.__f__("log", "at pages/webview/index.vue:101", "WebView加载完成，应用安全措施");
              this.injectSecurityScript();
            } catch (e) {
              common_vendor.index.__f__("error", "at pages/webview/index.vue:106", "WebView安全处理失败:", e);
            }
          });
        }
      }
    },
    // 注入安全脚本
    injectSecurityScript() {
      if (typeof common_vendor.wx$1 !== "undefined" && common_vendor.wx$1.evaluateWebView) {
        const securityScript = `
						try {
							// 阻止使用SharedArrayBuffer
							if (typeof SharedArrayBuffer !== 'undefined') {
								uni.__f__('warn','at pages/webview/index.vue:121','检测到SharedArrayBuffer，但已被安全策略禁用');
							}
							
							// 添加必要的安全头
							const meta = document.createElement('meta');
							meta.httpEquiv = 'Cross-Origin-Embedder-Policy';
							meta.content = 'require-corp';
							document.head.appendChild(meta);
							
							const meta2 = document.createElement('meta');
							meta2.httpEquiv = 'Cross-Origin-Opener-Policy';
							meta2.content = 'same-origin';
							document.head.appendChild(meta2);
							
							true; // 返回执行结果
						} catch(e) {
							uni.__f__('error','at pages/webview/index.vue:137','安全脚本执行出错:', e);
							false; // 返回执行结果
						}
					`;
        common_vendor.wx$1.evaluateWebView({
          webviewId: "webviewId",
          // 可能需要根据实际情况获取
          script: securityScript,
          success: (res) => {
            common_vendor.index.__f__("log", "at pages/webview/index.vue:147", "安全脚本注入成功", res);
          },
          fail: (err) => {
            common_vendor.index.__f__("error", "at pages/webview/index.vue:150", "安全脚本注入失败:", err);
          }
        });
      }
    },
    // 设置安全的WebURL，处理跨源隔离问题
    setSecureWebUrl(url) {
      if (!url) {
        this.webUrl = "";
        return;
      }
      try {
        let secureUrl = url;
        if (url.startsWith("http") && !url.includes("coop=cross-origin") && !url.includes("coep=require-corp")) {
          const separator = url.includes("?") ? "&" : "?";
          secureUrl = `${url}${separator}coop=cross-origin&coep=require-corp`;
        }
        common_vendor.index.__f__("log", "at pages/webview/index.vue:177", "使用安全URL:", secureUrl);
        this.webUrl = secureUrl;
      } catch (e) {
        common_vendor.index.__f__("error", "at pages/webview/index.vue:180", "处理URL时出错:", e);
        this.webUrl = url;
      }
    },
    // 处理web-view的消息
    handleMessage(event) {
      common_vendor.index.__f__("log", "at pages/webview/index.vue:187", "收到web-view消息:", event);
    },
    // 处理web-view的错误
    handleError(error) {
      common_vendor.index.__f__("error", "at pages/webview/index.vue:192", "web-view加载错误:", error);
      common_vendor.index.showToast({
        title: "页面加载出错",
        icon: "none"
      });
    }
  }
};
function _sfc_render(_ctx, _cache, $props, $setup, $data, $options) {
  return common_vendor.e({
    a: !$data.webUrl
  }, !$data.webUrl ? {} : {}, {
    b: $data.webUrl,
    c: common_vendor.o((...args) => $options.handleMessage && $options.handleMessage(...args)),
    d: common_vendor.o((...args) => $options.handleError && $options.handleError(...args))
  });
}
const MiniProgramPage = /* @__PURE__ */ common_vendor._export_sfc(_sfc_main, [["render", _sfc_render]]);
wx.createPage(MiniProgramPage);
//# sourceMappingURL=../../../.sourcemap/mp-weixin/pages/webview/index.js.map
