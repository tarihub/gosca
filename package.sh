# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

gosca_version=`./bin/gosca -v`
echo "build version: $gosca_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
      gosca_dir_name="${gosca_version}_${os}_${arch}"
      gosca_path="./packages/${gosca_version}_${os}_${arch}"

      if [ "x${os}" = x"windows" ]; then
          if [ ! -f "./gosca_${os}_${arch}.exe" ]; then
              continue
          fi
          mkdir ${gosca_path}
          mv ./gosca_${os}_${arch}.exe ${gosca_path}/gosca.exe
      else
          if [ ! -f "./gosca_${os}_${arch}" ]; then
              continue
          fi
          mkdir ${gosca_path}
          mv ./gosca_${os}_${arch} ${gosca_path}/gosca
      fi

      # packages
      cd ./packages
      if [ "x${os}" = x"windows" ]; then
          zip -rq ${gosca_dir_name}.zip ${gosca_dir_name}
      else
          tar -zcf ${gosca_dir_name}.tar.gz ${gosca_dir_name}
      fi
      cd ..
      rm -rf ${gosca_path}
    done
done

cd -