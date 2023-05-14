#!/usr/bin/env bash
declare -A repos=( ["nexusDatamodelTemplates"]="git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/datamodel-templates.git" ["nexusAppTemplates"]="git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/app-templates.git" ["nexusCompiler"]="git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/compiler.git" )
declare -A tags_repos

get_latest_tag(){
    git ls-remote -t --sort -v:refname $1 | awk '{print $2}' | awk -F'/' '{print $3}' | grep -E "\d\.\d+\.\d+$" | sort -V | tail -1 | tr -d '\n'
}

patch_values(){
    for key in ${!repos[@]}; do
         echo ${repos[$key]}
         tag=$(get_latest_tag ${repos[$key]})
         if [[ $tag != "" ]]; then
            tag_full="v${tag}"
            yq -i --arg tag $tag_full -Y '.'"$key"'.version=$tag'  pkg/common/values.yaml
         fi
    done
}

create_branch_and_commit(){
    git checkout periodic_update || git checkout -b periodic_update
    git add pkg/common/values.yaml && git commit -m "bump tags"
}

create_tag_if_not_exists(){
    current_tag=$(git ls-remote -t --sort -v:refname | awk '{print $2}' | awk -F'/' '{print $3}' | grep -E "\d\.\d+\.\d+$" | sort -V | tail -1 | tr -d '\n')
    new_tag=$(yq -r '.nexusCli.version' pkg/common/values.yaml | grep -E "\d\.\d+\.\d+$")
    if [[ ${new_tag:1} != $current_tag ]]; then
       echo "Creating tag because this is new ${new_tag:1}"
       git tag -a ${new_tag:1}
       git push origin ${new_tag:1}
    fi

}
func=$1
shift
$func "$@"

