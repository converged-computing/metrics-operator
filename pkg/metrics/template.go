/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package metrics

var (

	// TemplateConvertHostnames assumes a hostnames.txt to write to hostlist.txt
	TemplateConvertHostnames = `
# openmpi is evil and we need the ip addresses
echo "Starting to look for ip addresses..."
for h in $(cat ./hostnames.txt); do
	if [[ "$h" == "" ]]; then
	  continue
	fi
	address=""
	# keep trying until we have an ip address
	while [ "$address" == "" ]; do
		address=$(getent hosts $h | awk '{ print $1 }')
	done
	echo "${address}" >> ./hostlist.txt
done 
num_address=$(cat hostlist.txt | wc -l)
echo "Done finding ${num_address} ip addresses"		
`
)
