#!/bin/sh
pid=`pgrep gopa`
if test -z "$pid"
then
  echo "GOPA IS NOT RUNNING"
else
  echo "KILL GOPA PID $pid"
  kill -QUIT $pid
fi
