"use strict";
const common_vendor = require("../common/vendor.js");
const store_modules_user = require("./modules/user.js");
const store_modules_region = require("./modules/region.js");
const store = common_vendor.createStore({
  modules: {
    user: store_modules_user.user,
    region: store_modules_region.region
  }
});
exports.store = store;
//# sourceMappingURL=../../.sourcemap/mp-weixin/store/index.js.map
