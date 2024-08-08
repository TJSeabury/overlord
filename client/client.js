((d, w) => {
  /**
   * @typedef {Object} ErrorDetails
   * @property {string} errortext
   * @property {string} url
   * @property {string} line
   * @property {string} datetime
   * @property {string} os
   * @property {string} browser
   */

  /**
   * @property
   */
  class ShadowWatcher {
    constructor(token) {}
    /**
     *
     * @param {ErrorDetails} details
     */
    SendLog(details) {}
  }

  w.onerror = (error, url, line) => {
    controller.sendLog({
      acc: "error",
      data: "ERR:" + error + " URL:" + url + " L:" + line,
    });
  };
})(document, window);
