#!/bin/sh

dir=$(pwd)
base_dir=$HOME/data
mkdir -p ${base_dir}/coco; cd ${base_dir}/coco
wget https://s3-us-west-2.amazonaws.com/detectron/coco/coco_annotations_minival.tgz; unzip coco_annotations_minival.tgz
wget http://images.cocodataset.org/zips/train2014.zip; unzip train2014.zip
wget http://images.cocodataset.org/zips/val2014.zip; unzip val2014.zip
wget http://images.cocodataset.org/annotations/annotations_trainval2014.zip; unzip annotations_trainval2014.zip
cd $dir
