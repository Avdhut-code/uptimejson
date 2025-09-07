#!/usr/bin/env bash
# install.sh — install uptimejson binary + config + manpage (no systemd units)
# Usage: ./install.sh [--force]
set -euo pipefail
FORCE=false

while [[ ${1:-} != "" ]]; do
  case "$1" in
    --force) FORCE=true ;;
    -h|--help)
      cat <<EOF
install.sh installs uptimejson for a single user.
Options:
  --force   : overwrite existing config.json
  -h/--help : show this help
EOF
      exit 0
      ;;
    *) echo "Unknown arg: $1"; exit 2 ;;
  esac
  shift
done

TARGET_USER="${SUDO_USER:-${USER}}"
TARGET_HOME="$(eval echo "~${TARGET_USER}")"
BIN_NAME="uptimejson"
REPO_ROOT="$(pwd)"

echo "Installing for user: ${TARGET_USER}"
echo "Target HOME: ${TARGET_HOME}"

# 1) Build: produce binary named 'uptimejson' if not present
if [[ ! -x "${BIN_NAME}" ]]; then
  if command -v go >/dev/null 2>&1 && [[ -f "main.go" ]]; then
    echo "Building ${BIN_NAME}..."
    go build -o "${BIN_NAME}" || { echo "go build failed"; exit 1; }
  else
    echo "No built binary '${BIN_NAME}' and no Go source found. Place binary named '${BIN_NAME}' or install Go."
    exit 1
  fi
fi

# 2) Install binary to /usr/local/bin
if [[ -w /usr/local/bin ]] || [[ "${EUID}" -eq 0 ]]; then
  echo "Installing binary to /usr/local/bin/"
  install -Dm755 "${BIN_NAME}" /usr/local/bin/${BIN_NAME}
else
  echo "Elevated privileges required to place binary in /usr/local/bin/"
  sudo install -Dm755 "${BIN_NAME}" /usr/local/bin/${BIN_NAME}
fi

# # 3) Install default config into user's config dir (only if missing unless --force)
USER_CONFIG_DIR="${TARGET_HOME}/.config/uptimejson"
USER_CONFIG_PATH="${USER_CONFIG_DIR}/config.json"
# TEMPLATE_CONFIG="files/config.json"

# mkdir -p "${USER_CONFIG_DIR}"
# if [[ ! -f "${TEMPLATE_CONFIG}" ]]; then
#   echo "Warning: template ${TEMPLATE_CONFIG} not found. Skipping config copy."
# else
#   if [[ -f "${USER_CONFIG_PATH}" && "${FORCE}" != true ]]; then
#     echo "User config exists at ${USER_CONFIG_PATH} — not overwriting (use --force)."
#   else
#     echo "Installing default config to ${USER_CONFIG_PATH}"
#     if [[ "${EUID}" -eq 0 && -n "${SUDO_USER:-}" ]]; then
#       sudo -u "${TARGET_USER}" install -Dm644 "${TEMPLATE_CONFIG}" "${USER_CONFIG_PATH}"
#     else
#       install -Dm644 "${TEMPLATE_CONFIG}" "${USER_CONFIG_PATH}"
#     fi
#   fi
# fi

# 4) Ensure log directory exists
LOG_DIR="${TARGET_HOME}/.local/share/uptimejson"
if [[ "${EUID}" -eq 0 && -n "${SUDO_USER:-}" ]]; then
  sudo -u "${TARGET_USER}" mkdir -p "${LOG_DIR}"
else
  mkdir -p "${LOG_DIR}"
fi

if [[ "${EUID}" -eq 0 ]]; then
  chown -R "${TARGET_USER}:${TARGET_USER}" "${LOG_DIR}" 2>/dev/null || true
fi

# 5) Manpage (if provided)
MAN_PAGE_SRC="files/uptimejson.1"
if [[ -f "${MAN_PAGE_SRC}" ]]; then
  if [[ -w /usr/share/man/man1 ]] || [[ "${EUID}" -eq 0 ]]; then
    sudo install -Dm644 "${MAN_PAGE_SRC}" /usr/share/man/man1/uptimejson.1
    sudo gzip -f /usr/share/man/man1/uptimejson.1
    echo "Manpage installed to /usr/share/man/man1/uptimejson.1.gz"
  else
    echo "Manpage requires sudo. Installing with sudo..."
    sudo install -Dm644 "${MAN_PAGE_SRC}" /usr/share/man/man1/uptimejson.1
    sudo gzip -f /usr/share/man/man1/uptimejson.1
  fi
else
  echo "No man page found at ${MAN_PAGE_SRC}, skipping man install."
fi

cat <<EOF

Install complete.
- binary -> /usr/local/bin/uptimejson
- config -> ${USER_CONFIG_PATH} (created if missing)
- log dir -> ${LOG_DIR}

Service files NOT installed. Put your unit files in repo/service_related/ and use make_service.sh to install them when ready.

Next (manual) steps for per-user service (if you want it):
  systemctl --user daemon-reload
  systemctl --user enable --now uptimejson-startup.service
  systemctl --user enable --now uptimejson-shutdown.service

If you ran this script with sudo, run the daemon-reload & enable commands as the target user.
EOF
