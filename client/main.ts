import { type ErrorDetails } from "./.generated/types";

declare const REPORTING_ENDPOINT: string;

((w) => {

  if (typeof REPORTING_ENDPOINT === "undefined") {
    throw new Error("ShadowWatcher: REPORTING_ENDPOINT is not defined.");
  }

  const self = document.currentScript as HTMLScriptElement;
  const token: string | null = self.getAttribute("data-token");
  if (token === null) {
    throw new Error("ShadowWatcher: Token is not defined.");
  }

  /**
   * This is the controller class for the client. It is used to watch for errors and send them to the server.
   */
  class ShadowWatcher {

    token: string;
    reportingEndpoint: string;

    constructor(token: string, reportingEndpoint: string) {
      this.token = token;
      this.reportingEndpoint = reportingEndpoint;

      this.sendLog = this.sendLog.bind(this);
      this.handleError = this.handleError.bind(this);
    }

    /**
     * Guesses the OS of the client.
     * @returns {string}
     */
    determineOS(): string {
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
    determineBrowser(): string {
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

    async sendLog(details: ErrorDetails) {
      const response = await fetch(this.reportingEndpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-ACCESS-TOKEN": this.token,
        },
        body: JSON.stringify(details),
      });
      if (!response.ok) {
        console.error("Error sending log to server", response);
      } else {
        console.log("Log sent to server");
      }
    }

    handleError(
      message: string | Event,
      url: string,
      source: string | undefined,
      line: number | undefined,
      column: number | undefined,
      error: Error | undefined
    ) {
      if (typeof message !== "string") {
        message = error?.message || "";
      }
      const details: ErrorDetails = {
        domain: w.location.hostname,
        errorText: message,
        url,
        filename: source || "",
        line: line || 0,
        column: column || 0,
        datetime: new Date().toISOString(),
        userAgent: navigator.userAgent,
      };
      this.sendLog(details);
    }
  }

  const shadowWatcher = new ShadowWatcher(
    token,
    REPORTING_ENDPOINT
  );

  onerror = (message, source, lineno, colno, error) => {
    shadowWatcher.handleError(
      message,
      w.URL.toString(),
      source,
      lineno,
      colno,
      error
    );
  };

  console.log("ShadowWatcher loaded");
})(window);
