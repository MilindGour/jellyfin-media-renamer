#!/usr/bin/env bash

# This script is the starting point of jmr
# It starts both backend and frontend in production

# Start backend
./backend/jmr &
BACKEND_PID=$!

# Start frontend
PORT=7750 $(which node) ./frontend &
FRONTEND_PID=$!

trap "kill $BACKEND_PID $FRONTEND_PID" SIGINT SIGTERM
wait $FRONTEND_PID
