#!/bin/sh

# https://github.com/tensorflow/models/blob/master/research/object_detection/g3doc/preparing_inputs.md#generating-the-oxford-iiit-pet-tfrecord-files

dir=$(pwd)
base_dir=$HOME/data
mkdir -p ${base_dir}/pet; cd ${base_dir}/pet
wget http://www.robots.ox.ac.uk/~vgg/data/pets/data/images.tar.gz
wget http://www.robots.ox.ac.uk/~vgg/data/pets/data/annotations.tar.gz
tar -xvf annotations.tar.gz
tar -xvf images.tar.gz
cd $dir
