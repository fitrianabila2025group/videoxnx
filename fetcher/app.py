"""
Botasaurus-powered HTML fetcher microservice.

Exposes a single endpoint:
    POST /fetch  { "url": "...", "wait": 4 }  ->  { "status": 200, "html": "..." }

The Go backend calls this service when SCRAPER_FETCHER_URL is configured. Botasaurus
drives an undetectable Chromium instance that solves Cloudflare challenges automatically.
"""
from __future__ import annotations

import logging
import os
import threading
import time
from typing import Optional

from flask import Flask, jsonify, request

# Botasaurus high-level driver
from botasaurus_driver import Driver

LOG = logging.getLogger("fetcher")
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s %(message)s")

app = Flask(__name__)

# A single Driver instance shared across requests; we serialize access with a lock
# because a Chromium tab is not thread-safe.
_driver_lock = threading.Lock()
_driver: Optional[Driver] = None
_last_used: float = 0.0
IDLE_RECYCLE_SECONDS = int(os.getenv("FETCHER_IDLE_RECYCLE_SECONDS", "900"))


def _get_driver() -> Driver:
    global _driver, _last_used
    if _driver is None or (time.time() - _last_used) > IDLE_RECYCLE_SECONDS:
        if _driver is not None:
            try:
                _driver.close()
            except Exception:
                pass
        LOG.info("starting new botasaurus driver")
        _driver = Driver(headless=True, block_images=True)
    _last_used = time.time()
    return _driver


@app.route("/healthz", methods=["GET"])
def healthz():
    return jsonify({"ok": True})


@app.route("/fetch", methods=["POST"])
def fetch():
    body = request.get_json(silent=True) or {}
    url = (body.get("url") or "").strip()
    wait = int(body.get("wait") or 4)
    if not url.startswith(("http://", "https://")):
        return jsonify({"error": "invalid url"}), 400

    LOG.info("fetch %s", url)
    with _driver_lock:
        driver = _get_driver()
        try:
            # google_get is botasaurus's helper that visits via a referrer that helps avoid bot walls;
            # falls back to direct get() if not available on this driver version.
            try:
                driver.google_get(url, bypass_cloudflare=True)
            except Exception:
                driver.get(url, bypass_cloudflare=True)
            if wait > 0:
                driver.short_random_sleep() if hasattr(driver, "short_random_sleep") else time.sleep(wait)
            html = driver.page_html
            current_url = getattr(driver, "current_url", url)
            return jsonify({"status": 200, "html": html, "url": current_url})
        except Exception as e:
            LOG.exception("fetch failed: %s", e)
            return jsonify({"status": 500, "error": str(e)}), 502


if __name__ == "__main__":
    port = int(os.getenv("PORT", "9090"))
    app.run(host="0.0.0.0", port=port)
