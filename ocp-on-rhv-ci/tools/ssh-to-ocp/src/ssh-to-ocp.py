#!/usr/bin/python

from os import system,environ,path
from glob import glob

#change.log
# ocp-on-rhv cluster vm ssh tool

#ver 0.2 - eslutsky@redhat.com
# - support to run from mutilply users
# - added dialog to select which cluster to connect to

#ver 0.3 - eslutsky@redhat.com
# - added vm count for each cluster

LEASE_FILE="/var/lib/dnsmasq/net-1*.leases"
SSH_OPTS="-oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null"
SSH_USER="core"
SSH_KEY="~/id_rsa"
TMUX_SESSION_NAME="ocp-on-rhv VMs"

class Leases(object):
    def __init__(self):
        self.leases=[]
        self.get_lease_files()

    def get_dialog_options(self):
        lst = []
        for cnt,fpath in enumerate(self.leases):
            num_lines = sum(1 for line in open(fpath))
            lst.append('"%d" "cluster - ovirt-1%d (%d running VMs)"' % (cnt,cnt,num_lines))

        return  " ".join(lst)

    def get_file_by_id(self,id):
        return self.leases[id]

    def get_lease_files(self):
         self.leases=glob("%s*" % LEASE_FILE)

def get_dialog_options():
    return LEASES.get_dialog_options()

def get_file_by_id(id):
    return LEASES.get_file_by_id(int(id))


def kill_session():
    system("tmux kill-session -t \"%s\"" % (TMUX_SESSION_NAME) )

def attach_session():
    system("tmux attach -t \"%s\"" % (TMUX_SESSION_NAME))
def new_session():
    system("tmux new -s \"%s\" -d" % (TMUX_SESSION_NAME))

def open_dialog():
    cmd_line = """
    set -o allexport;
    exec 3>&1 ;
    result=$(dialog  \
           --backtitle "Cluster Selection" \
           --title "Menu"  \
           --clear   \
           --cancel-label "Exit" \
           --menu "Please select:" 20 0 15 \
            %s 2>&1 1>&3);
    echo $result > ~/result
    """ % get_dialog_options()
    print "running " + cmd_line
    system(cmd_line)

    return open( path.expanduser("~/result")).read()

def run_tmux(cmd="",window_name=""):

    cmd_line = "tmux new-window -t \"%s\" -n %s \"%s\"" % (TMUX_SESSION_NAME,window_name,cmd)
    print cmd_line
    system(cmd_line)

def get_ssh(vm_address):
    return "ssh %s %s@%s -i %s" % (SSH_OPTS,SSH_USER,vm_address,SSH_KEY)

def lease_parser(file_name):
    fields = ("time","mac","ipaddress","vmname","mac2")
    db = []

    with open (file_name) as file:
        for line in file:
            row={}
            for c,f in enumerate(line.split(" ")):
                if f == "*" :
                    row [ fields[ c ] ] = "bootstrap"
                else:
                    row [ fields[ c ] ] = f
            db.append(row)
            yield row

    #return db
            #print data.split("\n")


    #return db
            #print data.split("\n")

LEASES = Leases()

if __name__ == "__main__":
    while [ True ] :
        lease_id = open_dialog()

        kill_session()
        new_session()
        for x in lease_parser( get_file_by_id(lease_id)  ):
            cmd=get_ssh(x['ipaddress'])
            run_tmux(cmd,x['vmname'])
            #print x['vmname']

        attach_session()
   # run_tmux()
#generator to iterate over the leases file


