import { LogManager } from "aurelia-framework";
export var log = LogManager.getLogger("graph");

import { showMessage } from "./actions";

import { connectTo } from "aurelia-store";

var $ = require("jquery");

Date.prototype.timeNow = function () {
  return (
    (this.getHours() < 10 ? "0" : "") +
    this.getHours() +
    ":" +
    (this.getMinutes() < 10 ? "0" : "") +
    this.getMinutes() +
    ":" +
    (this.getSeconds() < 10 ? "0" : "") +
    this.getSeconds()
  );
};

@connectTo()
export class StatusGraphCustomElement {
  constructor(signaler) {
    var now = new Date().timeNow();

    this.currentMax = 60;
    this.source = "";
    this.apiPort;

    this.dataUpload = {
      x: [now],
      y: [0],
      mode: "lines+markers",
      name: "upload",
      line: {
        color: "rgb(219, 64, 82)",
        width: 1,
        shape: "spline",
      },
    };
    this.dataDownload = {
      x: [now],
      y: [0],
      mode: "lines+markers",
      name: "download",
      line: {
        color: "rgb(55, 128, 191)",
        width: 1,
        shape: "spline",
      },
    };
    this.layout = {
      width: window.innerWidth,
      height: window.innerHeight * 0.8,
      xaxis: {
        autotick: true,
        ticks: "outside",
        tick0: 0,
        ticklen: 8,
        tickwidth: 4,
        tickcolor: "#000",
        nticks: 20,
        title: {
          text: "Time (seconds)",
          standoff: 20,
        },
      },
      yaxis: {
        autotick: true,
        rangemode: "tozero",
        ticks: "outside",
        tick0: 0,
        ticklen: 8,
        tickwidth: 4,
        nticks: 20,
        tickcolor: "#000",
        title: {
          text: "Speed (Kb/s)",
          standoff: 20,
        },
      },
    };

    this.updateDataTimer = setInterval(() => this.periodicDataUpdate(), 1100);
  }

  attached() {
    this.gd = document.getElementById("gd");
    Plotly.newPlot(this.gd, [this.dataUpload, this.dataDownload], this.layout, {
      responsive: true,
    });
  }

  detached() {
    clearInterval(this.updateDataTimer);
  }

  periodicDataUpdate() {
    var $tab = $("status-graph");
    if ($tab.is(":visible") !== true) return; // skip update

    fetch(this.source)
      .then((response) => {
        if (!response.ok) {
          showMessage(`HTTP Error Status: ${response.status}`, "error", 1000);
          return cb({
            data: [],
          });
        }

        return response.json();
      })
      .then((obj) => {
        var up = null;
        var dw = null;
        log.info(obj);
        obj = obj.data;

        for (var i = 0; i < obj.length; i++) {
          if (obj[i].attribute == "Current Download Speed") {
            dw = ~~obj[i].value;
            log.info("dw", dw);
          }
          if (obj[i].attribute == "Current Upload Speed") {
            up = ~~obj[i].value;
            log.info("up", up);
          }
        }

        this.periodicGraphUpdate({
          upload: up !== null ? up : 0,
          download: dw !== null ? dw : 0,
        });
      })
      .catch((error) => {
        showMessage(error, "error", 1000);
        this.periodicGraphUpdate({
          upload: 0,
          download: 0,
        });
      });
}

  periodicGraphUpdate(data) {
    var now = new Date().timeNow();
    var up = data.upload;
    var dw = data.download;

    // ensure the data is valid
    if (up === undefined || up === null || dw === undefined || up === undefined)
      return;

    // discard old values
    /*var currMax = Math.max(1, (this.dataUpload.x.length - this.currentMax) / 2);
    if (this.dataUpload.x.length > this.currentMax) {
      for (var i = 0; i < currMax; i++) {
        this.dataUpload.y.shift();
        this.dataDownload.y.shift();
        this.dataUpload.x.shift();
        this.dataDownload.x.shift();
      }
    }

    this.dataUpload.x.push(now);
    this.dataDownload.x.push(now);

    this.dataUpload.y.push(up);
    this.dataDownload.y.push(dw);*/

    Plotly.extendTraces(
      this.gd,
      {
        x: [[now], [now]],
        y: [[up], [dw]],
      },
      [0, 1]
    );
  }

  stateChanged(newState, oldState) {
    this.apiPort = newState.port;
    this.source = `http://127.0.0.1:${this.apiPort}/api/v1/client/statistics/data`;
  }
}
