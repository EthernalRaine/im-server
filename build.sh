#!/bin/bash

usage()
{
    echo "===========================HELP============================"
    echo "build.sh [-s verbosity configuration (silent/verbose)] [-b build configuration (full/stripped)]"
    exit 1
}

while getopts 's:b:h?' opt; do
    case "$opt" in
        s)
            mode=$OPTARG
            ;;
        b)
            echo "running application go build (${OPTARG})"
            
            oldnum=`cat version`  
            newnum=`expr $oldnum + 1`
            sed -i "s/$oldnum\$/$newnum/g" version
            
            rm src/utility/version
            cp version src/utility/

            cd src

            if [[ $OPTARG = "stripped" ]]
            then
                if [[ $mode = "silent" ]]
                then
                    go build -o ../build/im-next -ldflags '-s'
                fi

                if [[ $mode = "verbose" ]]
                then
                    go build -o ../build/im-next -ldflags '-s' -x -v
                fi
            fi

            if [[ $OPTARG = "full" ]]
            then
                if [[ $mode = "silent" ]]
                then
                    go build -o ../build/im-next
                fi

                if [[ $mode = "verbose" ]]
                then
                    go build -o ../build/im-next -x -v
                fi
            fi

            cd utility
            rm version
            echo "placeholder" > version

            echo "done building."
            echo
            exit 0
            ;;
        *)
            usage
            ;;
    esac
done
shift "$(($OPTIND -1))"
usage