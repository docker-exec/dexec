#!/usr/bin/env bash

function get_cwd() {
  pushd $(dirname ${0}) >/dev/null
  script_path=$(pwd -P)
  popd >/dev/null
  echo "${script_path}"
}

function get_snapshot_plugin() {
  if ! grep -q vagrant-vbox-snapshot <(vagrant plugin list); then
    vagrant plugin install vagrant-vbox-snapshot
  fi
}

function up() {
  vagrant status | awk 'NR==3' | grep -q 'not created'
  vm_does_not_exist=$?

  vagrant box list | grep -q ubuntu/bionic64
  has_box=$?

  if [[ ${has_box} -ne 0 ]]; then
    vagrant box add ubuntu/bionic64
  fi

  vagrant up

  if [[ ${vm_does_not_exist} -eq 0 ]]; then
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
  local cwd="$(get_cwd)"
  pushd "${cwd}" >/dev/null
  up
  if [[ -z "$1" ]] || [[ "$1" != "--no-clean" ]]; then
    echo "Restoring to initial state"
    restore
  else
    echo "Skipping restore to initial state"
  fi
  vagrant ssh -c "
  cd /home/vagrant/.go/src/github.com/docker-exec/dexec
  go get
  go install
  bats .acceptance-tests/bats/dexec.bats"
  down
  popd >/dev/null
}

function destroy() {
  local cwd="$(get_cwd)"
  pushd "${cwd}" >/dev/null
  echo "Destroying vagrant box"
  vagrant destroy -f

  popd >/dev/null
}

case $1 in
run)
  run $2;;
destroy)
  destroy;;
*)
  echo "Usage:"
  echo "   $0 run [--no-clean]"
  echo "   $0 destroy"
  ;;
esac
