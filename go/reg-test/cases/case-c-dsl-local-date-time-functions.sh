# See also the system date command:
# export TZ=America/Sao_Paulo; date -j -f "%Y-%m-%d %H:%M:%S %Z" "2017-02-19 00:30:00" +%s
# export TZ=America/Sao_Paulo; date -r  86400 +"%Y-%m-%d %H:%M:%S %Z"

export TZ=America/Sao_Paulo
#echo TZ=$TZ
run_mlr --opprint put '$b=localtime2sec($a); $c=sec2localtime($b); $d=sec2localdate($b)' <<EOF
a=2017-02-18 23:00:00
a=2017-02-18 23:59:59
a=2017-02-19 00:00:00
a=2017-02-19 00:30:00
a=2017-02-19 01:00:00
a=2017-10-14 23:00:00
a=2017-10-14 23:59:59
a=2017-10-15 00:00:00
a=2017-10-15 00:30:00
a=2017-10-15 01:00:00
EOF
export TZ=

export TZ=America/Sao_Paulo
#echo TZ=$TZ
run_mlr --opprint put '$b=localtime2sec($a); $c=sec2localtime($b); $d=sec2localdate($b)' <<EOF
a=2017-02-14 00:00:00
a=2017-02-15 00:00:00
a=2017-02-16 00:00:00
a=2017-02-17 00:00:00
a=2017-02-18 00:00:00
a=2017-02-19 00:00:00
a=2017-02-20 00:00:00
a=2017-10-12 00:00:00
a=2017-10-13 00:00:00
a=2017-10-14 00:00:00
a=2017-10-15 00:00:00
a=2017-10-16 00:00:00
a=2017-10-17 00:00:00
a=2017-10-18 00:00:00
a=2017-10-19 00:00:00
EOF
export TZ=

export TZ=America/Sao_Paulo
#echo TZ=$TZ
run_mlr --opprint put '$b=strptime_local($a, "%Y-%m-%d %H:%M:%S"); $c=strftime_local($b, "%Y-%m-%d %H:%M:%S")' <<EOF
a=2017-02-18 23:00:00
a=2017-02-18 23:59:59
a=2017-02-19 00:00:00
a=2017-02-19 00:30:00
a=2017-02-19 01:00:00
a=2017-10-14 23:00:00
a=2017-10-14 23:59:59
a=2017-10-15 00:00:00
a=2017-10-15 00:30:00
a=2017-10-15 01:00:00
EOF
export TZ=

export TZ=America/Sao_Paulo
#echo TZ=$TZ
run_mlr --opprint put '$b=strptime_local($a, "%Y-%m-%d %H:%M:%S"); $c=strftime_local($b, "%Y-%m-%d %H:%M:%S")' <<EOF
a=2017-02-18 23:00:00
a=2017-02-18 23:59:59
a=2017-02-19 00:00:00
a=2017-02-19 00:30:00
a=2017-02-19 01:00:00
a=2017-10-14 23:00:00
a=2017-10-14 23:59:59
a=2017-10-15 00:00:00
a=2017-10-15 00:30:00
a=2017-10-15 01:00:00
EOF
export TZ=
