set -x
MR_TITLE="Dra"
if [[ $MR_TITLE =~ "Draft" || $MR_TITLE =~ "DRAFT" ]] ; then
        echo "skipping checks because MR is marked as draft"
        exit 0
fi