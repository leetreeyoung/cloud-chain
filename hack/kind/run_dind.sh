#!/usr/bin/env bash
docker rm -f node-0
docker run --privileged -itd --name node-0 docker:dind

syncPort=2022
apk --update add --no-cache curl tzdata openssh bash rsync openssl xxhash && rm -rf /var/cache/apk/*
sed -i "s/#Port.*/Port ${syncPort}/g" /etc/ssh/sshd_config
sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config
sed -i "s/#   StrictHostKeyChecking ask/StrictHostKeyChecking no/g" /etc/ssh/ssh_config
ssh-keygen -t dsa -P "" -f /etc/ssh/ssh_host_dsa_key
chmod 600 /etc/ssh/ssh_host_dsa_key
ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key
chmod 600 /etc/ssh/ssh_host_rsa_key
ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key
ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key
echo "root:$(cat /dev/random | sed 's/[^a-zA-Z0-9]//g' | strings -n 12 |head -n 1)" | chpasswd
ssh-keygen -t rsa -f /root/.ssh/id_rsa -P '' -C "root@host" -q
cp /root/.ssh/id_rsa.pub /root/.ssh/authorized_keys

/usr/sbin/sshd -D -e -f /etc/ssh/sshd_config 2>&1 &






