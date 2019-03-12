#!/bin/sh

jmeter_root_folder=$1
report_folder_name=$2

if [ -z "${jmeter_root_folder}" ] || [ -z "${report_folder_name}" ] ; then
    echo "Invalid arguments"
else
    report_folder="${jmeter_root_folder}/reports/${report_folder_name}"
    report_log_xml="${report_folder}/log.csv"
    report_detail_folder="${report_folder}/detail"
    jmeter_log="${report_folder}/jmeter.log"

    echo "report_folder: ${report_folder}"
    echo "report_log_xml: ${report_log_xml}"
    echo "report_detail_folder: ${report_detail_folder}"
    echo "jmeter_log: ${jmeter_log}"

    rm -rf ${report_detail_folder}

    cmd="jmeter -g ${report_log_xml} -j ${jmeter_log} -o ${report_detail_folder}"

    echo $cmd

    eval "$cmd"

    cp ${report_detail_folder}/index.html ${report_detail_folder}/main.html

    report_path="/jmeter/reports/${report_folder_name}"

    echo "Check the detail report => ${report_path}/detail/index.html"
    echo "Check the result log => ${report_path}/log.csv"
fi
