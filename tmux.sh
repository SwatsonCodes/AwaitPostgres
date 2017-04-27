#!/bin/bash

TMUX_SESSION=${TMUX_SESSION:-awaitpostgres}

echo $TMUX_SESSION

tmux has-session -t ${TMUX_SESSION}
if [ $? -ne 0 ]
then
  tmux new-session -s ${TMUX_SESSION} -d

  tmux split-window -d -t 0 "watch -n 1 kubectl get pods,services,jobs,deployments"

  tmux select-layout tiled
fi

tmux attach -t ${TMUX_SESSION}
