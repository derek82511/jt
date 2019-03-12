#!/bin/sh

jvm_args=$1
script_name=$2
jmeter_root_folder=$3
report_folder_name=$4
config_args=$5
remote_host_config=$6

if [ -z "${jvm_args}" ] || [ -z "${script_name}" ] || [ -z "${jmeter_root_folder}" ] || [ -z "${report_folder_name}" ] ; then
    echo "Invalid arguments"
else
    script_file="${jmeter_root_folder}/scripts/${script_name}"
    report_folder="${jmeter_root_folder}/reports/${report_folder_name}"
    report_log_xml="${report_folder}/log.csv"
    report_detail_folder="${report_folder}/detail"
    jmeter_log="${report_folder}/jmeter.log"

    echo "jvm_args: ${jvm_args}"
    echo "script_file: ${script_file}"
    echo "report_folder: ${report_folder}"
    echo "report_log_xml: ${report_log_xml}"
    echo "report_detail_folder: ${report_detail_folder}"
    echo "jmeter_log: ${jmeter_log}"

    mkdir ${report_folder}
    mkdir ${report_detail_folder}

    cmd="JVM_ARGS=\"${jvm_args}\" jmeter ${config_args} -n -t ${script_file} ${remote_host_config} -l ${report_log_xml} -j ${jmeter_log} -e -o ${report_detail_folder}"

    echo $cmd

    eval "$cmd"

    cp ${report_detail_folder}/index.html ${report_detail_folder}/main.html

    report_path="/jmeter/reports/${report_folder_name}"

    echo "Finished testing."
    echo "Check the detail report => ${report_path}/detail/index.html"
    echo "Check the result log => ${report_path}/log.csv"
fi
