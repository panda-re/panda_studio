#!/bin/bash

apt install -y netcat

python3 /panda_studio/backend/runPANDAReplay.py | nc -q 3 ip 42069

