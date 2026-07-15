#!/usr/bin/env bash
# Create custom labels for voice-realtime. Safe to re-run (gh label create fails gracefully).
set -euo pipefail
REPO="${1:-lixuanqun/voice-realtime}"

create() {
  gh label create "$1" --color "$2" --description "$3" -R "$REPO" 2>/dev/null || \
    gh label edit "$1" --color "$2" --description "$3" -R "$REPO"
}

# Area
create "area/gateway"       "1d76db" "WebSocket gateway, routing, session, auth"
create "area/protocol"      "5319e7" "OpenAI Realtime mapping, audio codec"
create "area/provider"      "0e8a16" "Provider plugin framework and registry"
create "area/observability" "fbca04" "Logging, metrics, tracing"
create "area/docs"          "0075ca" "Documentation and tutorials"
create "area/ci"            "bfd4f2" "CI/CD and test infrastructure"

# Provider
create "provider/zhipu"      "c5def5" "Zhipu GLM-Realtime"
create "provider/stepfun"  "bfe5bf" "StepFun Realtime"
create "provider/volcengine" "d93f0b" "Volcengine Doubao dialogue"
create "provider/bailian"    "fef2c0" "Alibaba Bailian multimodal"
create "provider/new"        "e99695" "Request new cloud provider"

# Roadmap
create "roadmap/phase-0" "ededed" "Phase 0: scaffold and docs"
create "roadmap/phase-1" "c2e0c6" "Phase 1: zhipu + stepfun proxy"
create "roadmap/phase-2" "f9d0c4" "Phase 2: volcengine translate"
create "roadmap/phase-3" "f9c513" "Phase 3: bailian state machine"
create "roadmap/phase-4" "b60205" "Phase 4: production hardening"

# Priority
create "priority/high"   "b60205" "Blocks mainline or integration"
create "priority/medium" "fbca04" "Normal schedule"
create "priority/low"    "c2e0c6" "Nice to have"

# Type
create "type/rfc" "d4c5f9" "Architecture or protocol RFC"

# Status
create "status/blocked"     "e11d21" "Blocked by external dependency"
create "status/in-progress" "1d76db" "Actively being worked on"

echo "Labels created/updated for $REPO"
