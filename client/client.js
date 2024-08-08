((w) => {
  /**
   * @typedef {Object} ErrorDetails
   * @property {string} domain
   * @property {string} errortext
   * @property {string} url
   * @property {string} line
   * @property {string} datetime
   * @property {string} os
   * @property {string} browser
   */

  /**
   * This is the controller class for the client. It is used to watch for errors and send them to the server.
   */
  class ShadowWatcher {
    /**
     * Creates an instance of ShadowWatcher.
     * @param {string} token
     */
    constructor(token, reportingEndpoint) {
      this.token = token;
      this.reportingEndpoint = reportingEndpoint;
      this.sendLog = this.sendLog.bind(this);
      this.handleError = this.handleError.bind(this);
    }

    /**
     * Guesses the OS of the client.
     * @returns {string}
     */
    determineOS() {
      if (navigator.userAgent.indexOf("Windows") > -1) {
        return "windows";
      } else if (navigator.userAgent.indexOf("Mac") > -1) {
        return "mac";
      } else if (navigator.userAgent.indexOf("Linux") > -1) {
        return "linux";
      } else {
        return "unknown";
      }
    }

    /**
     * Guesses the browser of the client.
     * @returns {string}
     */
    determineBrowser() {
      if (navigator.userAgent.indexOf("Chrome") > -1) {
        return "chrome";
      } else if (navigator.userAgent.indexOf("Firefox") > -1) {
        return "firefox";
      } else if (navigator.userAgent.indexOf("Safari") > -1) {
        return "safari";
      } else if (navigator.userAgent.indexOf("Edge") > -1) {
        return "edge";
      } else {
        return "unknown";
      }
    }

    /**
     * Sends the log to the server.
     * @param {ErrorDetails} details
     */
    async SendLog(details) {
      const response = await fetch(this.reportingEndpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(details),
      });
      if (!response.ok) {
        console.error("Error sending log to server", response);
      } else {
        console.log("Log sent to server");
      }
    }

    /**
     * Handles errors from the client.
     * @param {ErrorDetails} details
     */
    handleError(error, url, line) {
      /** @type {ErrorDetails} */
      const details = {
        domain: w.location.hostname,
        error,
        url,
        line,
        datetime: new Date().toISOString(),
        os: this.determineOS(),
        browser: this.determineBrowser(),
      };
      this.SendLog(details);
    }
  }

  const shadowWatcher = new ShadowWatcher(token, reportingEndpoint);

  w.onerror = (error, url, line) => shadowWatcher.handleError(error, url, line);
})(window);
