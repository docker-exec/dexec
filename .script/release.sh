#!/usr/bin/env bash

app_name=dexec
release_version=1.0.8
build_path=$(pwd -P)/build

rm -rf ${build_path}
mkdir ${build_path}
gox -output="${build_path}/${app_name}_${release_version}_{{.OS}}_{{.Arch}}/${app_name}" -os='linux windows darwin' -arch='386 amd64'
for bin_path in $(find -E ${build_path} -type d -depth 1 | perl -pe "s|^${build_path}\/||"); do
    target_file=$(perl -pe 's/^(.*)\.[^\.]$/$1/' <<<${bin_path})
    target_os=$(perl -pe 's/^\w+_[\d\.]+_(\w+)_\w+$/$1/' <<<${bin_path})
    target_arch=$(perl -pe 's/^\w+_[\d\.]+_\w+_(\w+)$/$1/' <<<${bin_path})
    case ${target_os} in
        windows)
            pushd ${build_path}/${bin_path} >/dev/null
            zip -9 -q ${build_path}/${target_file}.zip ${app_name}.exe
            popd >/dev/null
            rm -rf ${build_path}/${bin_path}
            ;;
        *)
            tar -czf ${build_path}/${target_file}.tar.gz -C ${build_path}/${bin_path} ${app_name} && rm -rf ${build_path}/${bin_path}
            ;;
    esac
done
