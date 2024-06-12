json=$(curl https://gcr.io/v2/nsx-sm/nexus/nexus-cli/tags/list | jq  -r '.manifest[] | .tag | select(.[]=="latest") | .[]')
declare -a all_apps_array=($(echo $json | tr "\n" " "))
for arrayo in ${all_apps_array[@]}; do 
   echo "----"
   echo $arrayo
done