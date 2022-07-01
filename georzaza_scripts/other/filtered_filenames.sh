#!/bin/bash
for file in $(find . -type f);do printf "%s:%s\n" ${file:2} `echo "$file" | awk 'BEGIN{FS="/"}{print $NF}'`; done | grep -v "png$\|\.qml$\|\.js$\|\.html$\|test\.go\>\|^\.git\|gitignore\|\.json$\|\.txt$\|\.md$\|LICENSE$\|TODO$\|\.cov$\|\.css$\|\.ethtest$\|\.sh$\|georzaza\|\.yml$\|Dockerfile\|^_data\/chain" | awk 'BEGIN {FS=":"} {printf "%50s\t\t%s\n",$1,$2}'  | sort -k2 > georzaza_scripts/filtered_filenames 