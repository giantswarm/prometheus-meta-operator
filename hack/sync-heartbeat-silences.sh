#!/bin/bash

# This script help with heartbeat and silences synchronization.
#
# Whenever a cluster is know to have issues, a silence is created in https://github.com/giantswarm/silences repository.
# In some cases this silence might apply to the whole cluster,
# when this is the case heartbeat in Opsgenie for the corresponding cluster
# should be disabled.
# When the silence is then removed the heartbeat in Opsgenie should be re-enabled.
#
# There is currently no automation for this, so this script is here to help humans to do their job right.
#
# Requirements to run this script :
# - heartbeatctl: https://github.com/giantswarm/heartbeatctl
# - silencectl: https://github.com/giantswarm/silencectl

command -V heartbeatctl
command -V silencectl

echo "> start"

echo "> re-enable disabled heartbeats (for cluster with no silences)"
for h in $(heartbeatctl list -s=DISABLED --no-headers|awk '{print $1}'); do grep -q $h <(silencectl list) || heartbeatctl enable $h; done

echo "> disable heartbeat (for cluster with silences)"
silencectl list | xargs heartbeatctl disable

echo "> end"
