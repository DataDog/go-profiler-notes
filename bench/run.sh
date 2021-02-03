os_arch() {
  echo "$(uname | tr '[:upper:]' '[:lower:]')_$(uname -m)"
}

#go run . \
  #-workloads mutex,chan \
  #-ops 100000 \
  #-blockprofilerates 0,1,10,100,1000,10000,100000,1000000 \
  #-runs 20 \
  #-bufsizes 0 \
  #-depths 16 \
  #> "block_$(os_arch).csv"

go run . \
  -workloads chan \
  -ops 100000 \
  -blockprofilerates 0,1,10,100,1000,10000,100000,1000000 \
  -bufsizes 128 \
  -runs 20 \
  -depths 16 \
  > "block_bufchan_$(os_arch).csv"
