#!/bin/bash

users=(eslutsky rgolangh gzaidman dfediuck)

for user in ${users[@]} ; do
echo -e "adding user $user \n"
useradd $user

echo -e "adding ssh key"
mkdir /home/${user}/.ssh/
chmod 700 /home/${user}/.ssh/
touch /home/${user}/.ssh/authorized_keys
chmod 600 /home/${user}/.ssh/authorized_keys
curl https://github.com/${user}.keys >/home/${user}/.ssh/authorized_keys
chown -R ${user}:${user} /home/${user}
usermod -a -G wheel ${user}
echo -e "note: \n\tupdate /etc/sudoers to allow wheel users,escalation with out entering password"
done
