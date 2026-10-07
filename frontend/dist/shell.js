// The bearer token lives on this object, never in a Datastar signal.
// Signals are posted to the server with every action.
(function () {
  var desktop = {
    token: "",
    headers: function () {
      return { Authorization: "Bearer " + this.token };
    },
    bootstrapOptions: function () {
      return {
        headers: this.headers(),
        openWhenHidden: true,
        retry: "auto",
        retryMaxCount: 5,
        retryInterval: 1000,
        retryScaler: 2,
        retryMaxWait: 8000,
      };
    },
    connect: function () {
      document.dispatchEvent(new CustomEvent("desktop-connect"));
    },
    save: function (raw) {
      var token = String(raw == null ? "" : raw).trim();
      var input = document.getElementById("shell-token");
      if (!token) {
        fault("Token is empty.", "signed-out");
        return;
      }
      var saveToken = window.go && window.go.main && window.go.main.App && window.go.main.App.SaveToken;
      if (!saveToken) {
        fault("Desktop bridge is not available.", "signed-out");
        return;
      }
      saveToken(token).then(function () {
        desktop.token = token;
        if (input) input.value = "";
        desktop.connect();
      }, function (err) {
        fault(message(err), "signed-out");
      });
    },
  };
  window.desktop = desktop;

  function message(err) {
    if (!err) return "Something went wrong.";
    if (typeof err === "string") return err;
    if (err.message) return err.message;
    return String(err);
  }

  function fault(text, phase) {
    document.dispatchEvent(new CustomEvent("desktop-fault", {
      detail: { message: text, phase: phase || "offline" },
    }));
  }

  function whenDatastarReady() {
    return new Promise(function (resolve) {
      document.addEventListener("datastar-ready", function () { resolve(); }, { once: true });
    });
  }

  function bridge() {
    return new Promise(function (resolve, reject) {
      var tries = 0;
      (function poll() {
        var session = window.go && window.go.main && window.go.main.App && window.go.main.App.Session;
        if (session) {
          resolve(session);
          return;
        }
        tries += 1;
        if (tries > 50) {
          reject(new Error("Desktop bridge is not available."));
          return;
        }
        setTimeout(poll, 20);
      })();
    });
  }

  Promise.all([bridge(), whenDatastarReady()]).then(function (ready) {
    return ready[0]();
  }).then(function (session) {
    desktop.token = (session && session.token) || "";
    document.dispatchEvent(new CustomEvent("desktop-session", {
      detail: {
        apiBase: (session && session.apiBase) || "",
        hasToken: desktop.token !== "",
      },
    }));
  }, function (err) {
    fault(message(err), "offline");
  });
})();
