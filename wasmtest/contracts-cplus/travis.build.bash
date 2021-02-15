#!/usr/bin/env bash
set -ev

oldir=$(pwd)
currentdir=$(dirname $0)
cd $currentdir

git clone --recursive https://github.com/Daironode/aingle-wasm-cdt-cpp
cd aingle-wasm-cdt-cpp; git checkout v1.1 -b testframe; bash compiler_install.bash > /dev/null;cd ../
compilerdir="./aingle-wasm-cdt-cpp/install/bin"

for f in $(ls *.cpp)
do
	$compilerdir/ont_cpp $f -lbase58 -lbuiltins -o  ${f%.cpp}.wasm
done

rm -rf aingle-wasm-cdt-cpp
mv *.wasm ../testwasmdata/
rm *.wasm.str
cp  *.avm ../testwasmdata/

cd $oldir
