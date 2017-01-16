# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"

  config.vm.synced_folder Dir.pwd, "/home/ubuntu/src/github.com/dtan4/ct2timer"

  config.vm.provision "shell", inline: <<-SHELL
    wget -q https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
    tar zxvf go1.7.4.linux-amd64.tar.gz -C /usr/local
    echo 'export GOROOT=/usr/local/go' >> /home/ubuntu/.bashrc
    echo 'export GOPATH=$HOME' >> /home/ubuntu/.bashrc
    echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> /home/ubuntu/.bashrc
    add-apt-repository ppa:masterminds/glide
    apt update
    apt install -y cmake glide
  SHELL
end
