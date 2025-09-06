#!/usr/bin/env bash
# make_service.sh — install systemd unit files from repo/service_related
# Usage:
#   ./make_service.sh --user            # install into ~/.config/systemd/user/
#   sudo ./make_service.sh --system usrname  # install into /etc/systemd/system/ and set User=usrname when enabling
set -euo pipefail

MODE=""
TARGET_USER=""
SERVICE_DIR="service"

if [[ ${1:-} == "--user" ]]; then
  MODE="user"
  TARGET_USER="${SUDO_USER:-${USER}}"
elif [[ ${1:-} == "--system" ]]; then
  MODE="system"
  TARGET_USER="${2:-}"
  if [[ -z "${TARGET_USER}" ]]; then
    echo "Usage: sudo ./make_service.sh --system <username>"
    exit 2
  fi
else
  echo "Usage: ./make_service.sh --user   (or)  sudo ./make_service.sh --system <username>"
  exit 2
fi

if [[ ! -d "${SERVICE_DIR}" ]]; then
  echo "No ${SERVICE_DIR}/ directory found. Create it and drop your .service files there."
  exit 1
fi

if [[ "${MODE}" == "user" ]]; then
  TARGET_HOME="$(eval echo "~${TARGET_USER}")"
  DEST_DIR="${TARGET_HOME}/.config/systemd/user"
  mkdir -p "${DEST_DIR}"
  cp -v "${SERVICE_DIR}"/* "${DEST_DIR}/"
  if [[ "${EUID}" -eq 0 && -n "${SUDO_USER:-}" ]]; then
    chown -R "${TARGET_USER}:${TARGET_USER}" "${DEST_DIR}"
    echo "Installed user units to ${DEST_DIR} (owned by ${TARGET_USER})"
    echo "Run as the target user:"
    echo "  systemctl --user daemon-reload"
    echo "  systemctl --user enable --now <unit>"
  else
    echo "Installed user units to ${DEST_DIR}"
    echo "Now run:"
    echo "  systemctl --user daemon-reload"
    echo "  systemctl --user enable --now <unit>"
  fi
else
  # system mode: requires root; copy to /etc/systemd/system
  DEST_DIR="/etc/systemd/system"
  cp -v "${SERVICE_DIR}"/* "${DEST_DIR}/"
  systemctl daemon-reload
  echo "Installed system units to ${DEST_DIR}"
  echo "You can enable them for user ${TARGET_USER} using:"
  echo "  sudo systemctl enable --now uptimejson-startup.service"
  echo "  sudo systemctl enable --now uptimejson-shutdown.service"
  # Note: system units with User= must have the correct User= line already in the unit files
fi
