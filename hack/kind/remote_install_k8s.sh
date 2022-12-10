#!/usr/bin/env bash
remote=43.156.24.173
user=ubuntu
work_dir="/home/ubuntu"
key=$(cat ./.key)


upload_dir() {
  local local_dir=$1
  local remote_dir=$2
  sshpass -p "$key" scp -r "$local_dir" $user@$remote:"$remote_dir"
}

upload_file() {
  local local_dir=$1
  local remote_dir=$2
  sshpass -p "$key" scp -r "$local_dir" $user@$remote:"$remote_dir"
}

remote_exec() {
  local cmd=$1
  sshpass -p "$key" ssh $user@$remote "$cmd"
}
#main
#remote_exec 'sudo apt install docker.io containerd runc'
remote_exec "sudo mkdir -p $work_dir/dind"
remote_exec "sudo chown -R ubuntu:ubuntu $work_dir/dind"
upload_dir ../kind $work_dir/dind
remote_exec "ls $work_dir/dind/kind"
remote_exec "sudo bash $work_dir/dind/kind/run_dind.sh"
remote_exec 'sudo docker ps'