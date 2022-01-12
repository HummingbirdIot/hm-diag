
case $1 in
  run | '' ) 
    echo "run"
    sleep 2
    echo "run over"
    ;;
  stop ) 
    echo "stopHummingbirdMiner" ;;
  restartMiner )
    echo "restartMiner" ;;
  minerLog )
    # minerLog <since time> <until time> <grep string>
    echo "minerLog" "$2" "$3" "$4" ;;  
  toUpdate )
    t=`date +%s`
    r=$((t%2))
    res=`[[ $r -eq 0 ]] && echo "yes" || echo "no"`
    echo ">>>state:${res}" ;;
  * ) 
    echo "unknown subcommand !"
esac
