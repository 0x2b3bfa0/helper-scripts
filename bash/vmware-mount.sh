#!/bin/bash
#

sed -i 's/#user_allow_other/user_allow_other/g' /etc/fuse.conf

mount_dir="/mnt/hgfs"
host_str=".host:/"

/usr/bin/vmhgfs-fuse "$host_str" "$mount_dir" -o subtype=vmhgfs-fuse,allow_other,nonempty
