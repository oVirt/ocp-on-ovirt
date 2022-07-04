function send_event_to_ovirt(){
local install_state="Installed"
local build_id=${BUILD_ID}
jobspec=$(echo ${JOB_SPEC}|jq 'del(.decoration_config)')

if [ "$#" -eq 1  ] ; then
    install_phase=$1
fi

echo "Checking ovirt-engine ovirt-imageio-proxy"
curl --insecure \
--connect-timeout 2 \
--max-time 10 \
--retry 5  \
https://ovirt-engine.ocp-on-ovirt.gcp.devcluster.openshift.com:54323 -vv || true

#take the last 7 chars from the id and convert it to int
printf -v build_id '%d\n' $((10#${build_id: -7})) # ${build_id: -7}
epoch=$(date +'%s')
cat <<__EOF__ > ${SHARED_DIR}/event.xml
<event>
  <description>Openshift CI - cluster installation;${OCP_CLUSTER};${install_phase:=};${jobspec}</description>
  <severity>normal</severity>
  <origin>openshift-ci</origin>
  <custom_id>$((${epoch:=}+${build_id:=}))</custom_id>
</event>
__EOF__


  curl --insecure  \
  --request POST \
  --header "Accept: application/xml"  \
  --header "Content-Type: application/xml" \
  -u "${OVIRT_ENGINE_USERNAME}:${OVIRT_ENGINE_PASSWORD}" \
  -d @${SHARED_DIR}/event.xml \
  ${OVIRT_ENGINE_URL}/events || true
}