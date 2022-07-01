## Description:
## Get all the current folder's files recursively and for any (2 or more) same filenames under 
## different directories generate a script that will find and print the same lines of the files. 
## To avoid confusion: this/filepath/contains/the/filename.go 
##
## For details about how this script works see either the last "for" loop in this file
## or read the README.md file.


import os
from glob import glob

# Each filepath has the format of hello/world.go, this/filepath/contains/a/filename.go, etc
# The below dictionary will have it's keys being all the unique filenames, aka world.go and filename.go
# Each value of the dictionary will be an array with all the paths of the files that have the same name.
# example: for a given structure of 1/2/3.go 4/5/6.go 7/8/9.go a/b/3.go the below dictionary should finally become:
# { 
#   3.go: [1/2/3.go, a/b/3.go],
#   6.go: [4/5/6.go],
#   9.go: [7/8/9.go]
# }
uniq_filenames = {}


# get a list of all .go file paths recursively of the folder where this script is executed at
filepaths = [y for x in os.walk("./") for y in glob(os.path.join(x[0], '*.go'))]


# populate the uniq_filenames dictionary
for filepath in filepaths:
	filename = filepath.split("/")[-1]
	# if there is no suck key so far in our dictionary, add it, else append the filepath to it's list of values 
	if uniq_filenames.get(filename)	is None:
		uniq_filenames[filename] = [filepath]
	else:
		uniq_filenames[filename].append(filepath)


# create a shell script file
with open("common_lines.sh", "w") as handle:
	handle.write("#!/bin/bash\n")
	handle.write("# For details about this script either read the README.md file or the\n")
	handle.write("# relevant comment section of the python script find_candidate_same_files.py\n")
	handle.write("echo Generating a temporary folder\n")

	# create an output folder
	handle.write("mkdir -p georzaza_results\n")


# populate the script.
#
# The generated script will contain commands that print the same number of lines between any 2 filenames that 
# are the same but are saved under different filepaths, e.g. for files core/filter.go and event/filter/filter.go
#
# The results will be saved under the folder georzaza_results.
# To understand the format of the results generated take this for example:
# Files:  core/filter.go  and  event/filter/filter.go
# A folder named core-filter will be generated. 
# This folder will contain a file named event-filter-filter. 
# The latter file's contents will be the output of the command
# awk 'NR==FNR{arr[$0];next} $0 in arr' core/filter.go event/filter/filter.go
# which prints the number of lines of the core/filter.go that are present in the file event/filter/filter.go
#
# Note that for any 2 files with the same filename we only check them once. E.g. we will check core/filter.go 
# against event/filter/filter.go but not in the other way around.
# 

for filename,filepaths in uniq_filenames.items():
	if len(filepaths) > 1:
		for i in range(0, len(filepaths)-1):
			folder_to_save_under = "georzaza_results/" + filepaths[i].replace("./", "").replace("/", "-")[:-3]
			with open("common_lines.sh", "a") as handle:
				handle.write("mkdir -p " + folder_to_save_under + "\n")
			for j in range(i+1, len(filepaths)):
				if i == j:
					continue
				with open("common_lines.sh", "a") as handle:
					save_path = folder_to_save_under + "/" + filepaths[j].replace("./", "").replace("/", "-")[:-3]
					handle.write("awk 'NR==FNR{arr[$0];next} $0 in arr' %s %s > %s\n"% (filepaths[i], filepaths[j], save_path))


print("Generated script common_lines.sh")
