#!/bin/sh
pid=`pgrep gopa`
if test -z "$pid"
then
  echo "GOPA IS NOT RUNNING"
else
  echo "KILL GOPA PID:\n$pid"
  ps x|grep gopa|grep -v grep |awk '{print $1}'|xargs kill -QUIT
fi
