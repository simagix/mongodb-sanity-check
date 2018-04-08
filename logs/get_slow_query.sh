if [ "$1" == "" ]; then
    echo "usage: $0 log_file"
    exit
fi

grep ' .*ms$' $1 | awk '{t = $1; $1 = $NF;$2 = t;  print}' | sort -rn | head
