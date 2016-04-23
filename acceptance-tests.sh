#!/usr/bin/env bash

function get_snapshot_plugin() {
  if ! grep -q vagrant-vbox-snapshot <(vagrant plugin list); then
    vagrant plugin install vagrant-vbox-snapshot
  fi
}

function up() {
  vagrant status | awk 'NR==3' | grep -q 'not created'
  vm_does_not_exist=$?

  vagrant box list | grep -q ubuntu/trusty64
  has_box=$?

  if [ $has_box -ne 0 ]; then
    vagrant box add ubuntu/trusty64
  fi

  vagrant up

  if [ $vm_does_not_exist -eq 0 ]; then
    get_snapshot_plugin
    vagrant snapshot take 'post-bootstrap'
  fi
}

function down() {
  vagrant halt
}

function restore() {
  get_snapshot_plugin
  if grep 'post-bootstrap' <(vagrant snapshot list); then
    vagrant snapshot go 'post-bootstrap'
  fi
}

function run() {
  up
  if [ -z "$1" ] || [ "$1" != "--no-clean" ]; then
    echo "Restoring to initial state"
    restore
  else
    echo "Skipping restore to initial state"
  fi
  vagrant ssh -c "
  cd /home/vagrant/.go/src/github.com/docker-exec/dexec
  go get
  go install
  bats .bats/dexec.bats"
  down
}

case $1 in
run)
  run $2;;
*)
  echo "Usage:"
  echo "   ./acceptance-tests.sh run [--no-clean]"
  ;;
esac
