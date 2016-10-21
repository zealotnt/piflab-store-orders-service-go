#!/bin/bash

cat >> ~/.bashrc << EOF
ApiRun(){
gin -p 80 run
}
EOF

source ~/.bashrc