# syntax=docker/dockerfile:1.7
# Botasaurus needs a real Chromium. Use the official chromium image base.
FROM python:3.11-slim-bookworm

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    DEBIAN_FRONTEND=noninteractive \
    PIP_NO_CACHE_DIR=1

# System deps: chromium + fonts + libs needed by headless Chromium
RUN apt-get update && apt-get install -y --no-install-recommends \
      chromium chromium-driver \
      ca-certificates fonts-liberation \
      libnss3 libxss1 libasound2 libatk1.0-0 libatk-bridge2.0-0 \
      libdrm2 libgbm1 libgtk-3-0 libxcomposite1 libxdamage1 \
      libxfixes3 libxrandr2 libxshmfence1 libxkbcommon0 \
      curl xvfb tini \
   && rm -rf /var/lib/apt/lists/*

# Tell botasaurus where chromium lives (system chromium instead of auto-download)
ENV CHROME_BIN=/usr/bin/chromium \
    CHROMEDRIVER_PATH=/usr/bin/chromedriver \
    HOME=/home/app

WORKDIR /app
COPY requirements.txt .
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY app.py .

RUN useradd -m -u 10001 app && mkdir -p $HOME && chown -R app:app /app $HOME
USER app

EXPOSE 9090
ENV PORT=9090
ENTRYPOINT ["/usr/bin/tini","--"]
CMD ["sh","-c","exec gunicorn -b 0.0.0.0:${PORT:-9090} -w 1 -k gthread --threads 2 --timeout 120 app:app"]
