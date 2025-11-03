// runtime.go
package runtime

// ClientRuntime 返回客戶端 runtime JavaScript 代碼
// 這個腳本會自動綁定所有帶有 data-gvd-handler 屬性的元素事件
func ClientRuntime() string {
	return `
(function() {
  // 初始化 __gvd 全域物件
  window.__gvd = window.__gvd || {};
  window.__gvd.handlers = window.__gvd.handlers || {};

  // 綁定所有帶有 data-gvd-handler 屬性的元素
  function bindHandlers() {
    // 查找所有帶有 data-gvd-handler 屬性的元素
    var elements = document.querySelectorAll('[data-gvd-handler]');

    elements.forEach(function(el) {
      var handlerAttr = el.getAttribute('data-gvd-handler');
      if (!handlerAttr) return;

      // 解析 handlerID|eventType 格式
      var parts = handlerAttr.split('|');
      if (parts.length !== 2) return;

      var handlerID = parts[0];
      var eventType = parts[1];

      // 檢查是否為命名 handler（格式：named:functionName）
      if (handlerID.startsWith('named:')) {
        var fnName = handlerID.substring(6); // 移除 "named:" 前綴

        // 嘗試從全域作用域獲取函數
        var fn = window[fnName];
        if (typeof fn === 'function') {
          el.addEventListener(eventType, function(evt) {
            fn.call(el, evt, el);
          });
        } else {
          console.warn('Named handler not found: ' + fnName);
        }
      } else {
        // 內聯 JSAction handler
        var handler = window.__gvd.handlers[handlerID];
        if (handler && typeof handler.fn === 'function') {
          el.addEventListener(eventType, function(evt) {
            handler.fn.call(el, evt, el);
          });
        } else {
          console.warn('Handler not found for ID: ' + handlerID);
        }
      }
    });
  }

  // 在 DOM 準備就緒時綁定所有 handler
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function() {
      bindHandlers();
    });
  } else {
    // DOM 已經加載完成，立即綁定
    bindHandlers();
  }

  // 提供一個公開的方法來重新綁定 handler（用於動態內容）
  window.__gvd.rebind = function() {
    bindHandlers();
  };
})();
`
}
