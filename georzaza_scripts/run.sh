#!/bin/bash
cd ../
cp georzaza_scripts/find_candidate_same_files.py .
echo "Executing find_candidate_same_files.py"
python3 find_candidate_same_files.py
echo "Giving execution rights to common_lines.sh"
chmod +x common_lines.sh
echo "Executing common_lines.sh"
./common_lines.sh

# run those manually
#a=$(echo "CommonLines File1 File2 File1#Lines File2#Lines" && for file in `find . -type f`; do echo -e `cat $file | awk NF | wc -l` "\t" ${${file:2}%/*}.go "\t" ${${file:2}#*/}.go "\t" `echo ../${${file:2}%/*}.go | sed --expression='s/-/\//g' | xargs cat | wc -l` "\t" `echo ../${${file:2}#*/}.go | sed --expression='s/-/\//g' | xargs cat | wc -l`; done | sort -nrk1 | sed --expression='s/-/\//g')
#echo && echo $a | column -t --output-separator="  " | grep -v "_test\.go\>" 