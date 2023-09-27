/*
Copyright 2023 Lawrence Livermore National Security, LLC
 (c.f. AUTHORS, NOTICE.LLNS, COPYING)

SPDX-License-Identifier: MIT
*/

package addons

import (
	"fmt"
	"strings"

	api "github.com/converged-computing/metrics-operator/api/v1alpha2"
	"github.com/converged-computing/metrics-operator/pkg/metadata"
	"github.com/converged-computing/metrics-operator/pkg/specs"
	"k8s.io/apimachinery/pkg/util/intstr"
	jobset "sigs.k8s.io/jobset/api/jobset/v1alpha2"
)

// Flux Framework provides fully hierarchical graphs scheduler and resource manager
const (
	fluxIdentifier = "workload-flux"
)

type FluxFramework struct {
	SpackView

	// Target is the name of the replicated job to customize entrypoint logic for
	// This is what determines the size of the cluster
	target string

	// ContainerTarget is the name of the container to add flux to
	containerTarget string

	// mount is the location to install flux to
	mount     string
	pods      int32
	jobname   string
	namespace string

	// flux user and id need to match between containers (both are created)
	fluxUser      string
	fluxUid       string
	quorum        string
	tasks         int32
	optionFlags   string
	submitCommand string
	preCommand    string

	// volumeName to provide for the empty volume
	volumeName     string
	interactive    bool
	connectTimeout string
	debugZeroMQ    bool
	logLevel       string
	queuePolicy    string
	serviceName    string
	launcherLetter string
	workerLetter   string
	workerIndex    string
	launcherIndex  string
}

func (m FluxFramework) Family() string {
	return AddonFamilyWorkload
}

// Validate we have an executable provided, and args and optional
func (a *FluxFramework) Validate() bool {
	return true
}

// GetAddFluxUser gets string text to add the flux user
// this might need to vary depending on the OS
// We might also just be irresponsible and run as root :)
func (a *FluxFramework) getAddFluxUser() string {
	return fmt.Sprintf(`useradd -ms /bin/bash -u %s %s`, a.fluxUid, a.fluxUser)
}

func (m FluxFramework) AssembleVolumes() []specs.VolumeSpec {
	return m.GetSpackViewVolumes()
}

// Set custom options / attributes for the addon metric
func (a *FluxFramework) SetOptions(metric *api.MetricAddon, set *api.MetricSet) {

	a.EntrypointPath = "/metrics_operator/flux-entrypoint.sh"
	a.image = "ghcr.io/rse-ops/spack-flux-rocky-view:tag-8"
	a.SetDefaultOptions(metric)
	a.Mount = "/opt/share"
	a.VolumeName = "flux-volume"
	a.Identifier = fluxIdentifier
	a.fluxUser = "flux"
	a.fluxUid = "1004"
	a.interactive = false
	a.connectTimeout = "5s"
	a.debugZeroMQ = false
	a.logLevel = "6"
	a.pods = set.Spec.Pods
	a.jobname = set.Name
	a.namespace = set.Namespace
	a.serviceName = set.Spec.ServiceName
	a.queuePolicy = "fcfs"
	a.SpackViewContainer = "flux-framework"
	a.launcherIndex = "0"
	a.workerIndex = "0"
	a.launcherLetter = "l"
	a.workerLetter = "w"
	a.quorum = fmt.Sprintf("%d", a.pods)
	a.submitCommand = "submit"

	pc, ok := metric.Options["preCommand"]
	if ok {
		a.preCommand = pc.StrVal
	}
	wi, ok := metric.Options["workerIndex"]
	if ok {
		a.workerIndex = wi.StrVal
	}
	li, ok := metric.Options["launcherIndex"]
	if ok {
		a.launcherIndex = li.StrVal
	}
	mount, ok := metric.Options["mount"]
	if ok {
		a.mount = mount.StrVal
	}
	submit, ok := metric.Options["submit"]
	if ok {
		a.submitCommand = submit.StrVal
	}
	tasks, ok := metric.Options["tasks"]
	if ok {
		a.tasks = tasks.IntVal
	}
	fluxUid, ok := metric.Options["fluxUid"]
	if ok {
		a.fluxUid = fluxUid.StrVal
	}
	fluxuser, ok := metric.Options["fluxUser"]
	if ok {
		a.fluxUser = fluxuser.StrVal
	}

	workdir, ok := metric.Options["workdir"]
	if ok {
		a.workdir = workdir.StrVal
	}
	logLevel, ok := metric.Options["logLevel"]
	if ok {
		a.logLevel = logLevel.StrVal
	}
	target, ok := metric.Options["target"]
	if ok {
		a.target = target.StrVal
	}
	ctarget, ok := metric.Options["containerTarget"]
	if ok {
		a.containerTarget = ctarget.StrVal
	}
	image, ok := metric.Options["image"]
	if ok {
		a.image = image.StrVal
	}
	quorum, ok := metric.Options["quorum"]
	if ok {
		a.quorum = quorum.StrVal
	}
	ct, ok := metric.Options["connectTimeout"]
	if ok {
		a.connectTimeout = ct.StrVal
	}
	opts, ok := metric.Options["optionFlags"]
	if ok {
		a.optionFlags = opts.StrVal
	}
	interactive, ok := metric.Options["interactive"]
	if ok {
		if interactive.StrVal == "yes" || interactive.StrVal == "true" {
			a.interactive = true
		}
	}
	zmq, ok := metric.Options["debugZeroMQ"]
	if ok {
		if zmq.StrVal == "yes" || zmq.StrVal == "true" {
			a.debugZeroMQ = true
		}
	}

	// Create setup logic for flux from the view
	a.setSetup()
}

// generateRange is a shared function to generate a range string
func generateRange(size int32, start int32) string {
	var rangeString string
	if size == 1 {
		rangeString = fmt.Sprintf("%d", start)
	} else {
		rangeString = fmt.Sprintf("%d-%d", start, (start+size)-1)
	}
	return rangeString
}

// setSetup assumes flux installed in the view (/opt/view/bin)) and runs additional setup
// This includes generating the broker config, the curve certificate, and other config assets
func (a *FluxFramework) setSetup() {

	// fluxRoot for the view is in /opt/view/lib
	fluxRoot := "/opt/view"

	// Generate hostlists, this is the lead broker
	leadBroker := fmt.Sprintf("%s-%s-%s-0", a.jobname, a.launcherLetter, a.launcherIndex)
	workers := fmt.Sprintf("%s-%s-%s-[%s]", a.jobname, a.workerLetter, a.workerIndex, generateRange(a.pods-1, 0))
	hosts := fmt.Sprintf("%s,%s", leadBroker, workers)
	fqdn := fmt.Sprintf("%s.%s.svc.cluster.local", a.serviceName, a.namespace)

	// These shouldn't be formatted in block
	defaultBind := "tcp://eth0:%p"
	defaultConnect := "tcp://%h" + fmt.Sprintf(".%s:", fqdn) + "%p"

	setup := `#!/bin/sh
fluxuser=%s
fluxuid=%s
fluxroot=%s

# The mount for the view will be at the user defined mount / view
mount="%s/view"

echo "Hello I am hostname $(hostname) running setup."
# We only want one host to generate a certificate
mainHost="%s"

# Always use verbose, no reason to not here
echo "Flux username: ${fluxuser}"
echo "Flux install root: ${fluxroot}"
export fluxroot

# Add flux to the path
export PATH=/opt/view/bin:$PATH

# Cron directory
mkdir -p $fluxroot/etc/flux/system/cron.d
mkdir -p $fluxroot/var/lib/flux

# These actions need to happen on all hosts
mkdir -p $fluxroot/etc/flux/system
hosts="%s"
echo "flux R encode --hosts=${hosts} --local"
flux R encode --hosts=${hosts} --local > ${fluxroot}/etc/flux/system/R

echo
echo "📦 Resources"
cat ${fluxroot}/etc/flux/system/R

mkdir -p $fluxroot/etc/flux/imp/conf.d/

cat <<EOT >> ${fluxroot}/etc/flux/imp/conf.d/imp.toml
[exec]
allowed-users = [ "${fluxuser}", "root" ]
allowed-shells = [ "${mount}/libexec/flux/flux-shell" ]
EOT

echo
echo "🦊 Independent Minister of Privilege"
cat ${fluxroot}/etc/flux/imp/conf.d/imp.toml

# Write the broker configuration
mkdir -p ${fluxroot}/etc/flux/config
cat <<EOT >> ${fluxroot}/etc/flux/config/broker.toml
[exec]
imp = "${mount}/libexec/flux/flux-imp"

[access]
allow-guest-user = true
allow-root-owner = true

# Point to resource definition generated with flux-R(1).
[resource]
path = "${mount}/etc/flux/system/R"

[bootstrap]
curve_cert = "${mount}/etc/curve/curve.cert"
default_port = 8050
default_bind = "%s"
default_connect = "%s"
hosts = [
	{ host="${hosts}"},
]
[archive]
dbpath = "${mount}/var/lib/flux/job-archive.sqlite"
period = "1m"
busytimeout = "50s"

[sched-fluxion-qmanager]
queue-policy = "%s"
EOT

echo
echo "🐸 Broker Configuration"
cat ${fluxroot}/etc/flux/config/broker.toml

# If we are communicating via the flux uri this service needs to be started
chmod u+s ${fluxroot}/libexec/flux/flux-imp
chmod 4755 ${fluxroot}/libexec/flux/flux-imp
chmod 0644 ${fluxroot}/etc/flux/imp/conf.d/imp.toml

# The rundir needs to be created first, and owned by user flux
# Along with the state directory and curve certificate
mkdir -p ${fluxroot}/run/flux ${fluxroot}/etc/curve

# Generate the certificate (ONLY if the lead broker)
mkdir -p ${fluxroot}/etc/curve

if [[ "$(hostname)" == "$mainHost" ]]; then
    echo "I am the main host, generating shared certificate"
    $fluxroot/bin/flux keygen ${fluxroot}/etc/curve/curve.cert

	# Remove group and other read
	chmod o-r ${fluxroot}/etc/curve/curve.cert
	chmod g-r ${fluxroot}/etc/curve/curve.cert
	
	# Either the flux user owns the instance, or root
	# We must get the correct flux user id - this user needs to own
	# the run directory and these others
	chown -R ${fluxuid} ${fluxroot}/etc/curve/curve.cert
	
	echo
	echo "✨ Curve certificate"
	cat ${fluxroot}/etc/curve/curve.cert	
fi
`

	setup = fmt.Sprintf(
		setup,
		a.fluxUser,
		a.fluxUid,
		fluxRoot,
		a.Mount,
		leadBroker,
		hosts,
		defaultBind,
		defaultConnect,
		a.queuePolicy,
	)
	a.Setup = setup
}

// Exported options and list options
func (a *FluxFramework) Options() map[string]intstr.IntOrString {
	options := a.DefaultOptions()
	options["mount"] = intstr.FromString(a.mount)
	options["quorum"] = intstr.FromString(a.quorum)
	options["fluxUser"] = intstr.FromString(a.fluxUser)
	options["fluxUid"] = intstr.FromString(a.fluxUid)
	options["fluxUid"] = intstr.FromString(a.fluxUid)
	options["pods"] = intstr.FromInt(int(a.pods))
	options["connectTimeout"] = intstr.FromString(a.connectTimeout)
	options["logLevel"] = intstr.FromString(a.logLevel)
	options["jobname"] = intstr.FromString(a.jobname)
	options["namespace"] = intstr.FromString(a.namespace)
	options["serviceName"] = intstr.FromString(a.serviceName)
	options["queuePolicy"] = intstr.FromString(a.queuePolicy)
	options["launcherIndex"] = intstr.FromString(a.launcherIndex)
	options["launcherLetter"] = intstr.FromString(a.launcherLetter)
	options["workerIndex"] = intstr.FromString(a.workerIndex)
	options["workerLetter"] = intstr.FromString(a.workerLetter)
	options["submitCommand"] = intstr.FromString(a.submitCommand)
	return options
}

// CustomizeEntrypoint scripts
func (a *FluxFramework) CustomizeEntrypoints(
	cs []*specs.ContainerSpec,
	rjs []*jobset.ReplicatedJob,
) {
	for _, rj := range rjs {

		// Only customize if the replicated job name matches the target
		if a.target != "" && a.target != rj.Name {
			continue
		}
		a.customizeEntrypoint(cs, rj)
	}

}

// CustomizeEntrypoint for a single replicated job
// This is the portion that customizes our application to be run / submit by flux instead of by itself :)
func (a *FluxFramework) customizeEntrypoint(
	cs []*specs.ContainerSpec,
	rj *jobset.ReplicatedJob,
) {

	// Generate addon metadata
	meta := Metadata(a)

	interactive := ""
	if a.interactive {
		interactive = "-Sbroker.rc2_none"
	}
	zeromq := ""
	if a.debugZeroMQ {
		zeromq = "-Stbon.zmqdebug=1"
	}

	// This assumes a certain launcher letter for now
	// TODO allow to customize letter
	leadBroker := fmt.Sprintf("%s-%s-%s-0", a.jobname, a.launcherLetter, a.launcherIndex)

	// Watch only works with submit
	watch := ""
	if strings.Contains(a.submitCommand, "submit") {
		watch = "--watch"
	}

	// Prepare flags for flux. First, add big N if it's > pods, OR not set
	flags := ""
	if (a.tasks != 0 && a.tasks > a.pods) || a.tasks == 0 {
		flags = fmt.Sprintf(" -N %d", a.pods)
	}
	// Little n only gets added if it is set
	if a.tasks != 0 {
		flags += fmt.Sprintf(" -n %d %s -vvv", a.tasks, a.optionFlags)
	} else {
		flags += fmt.Sprintf(" %s -vvv", a.optionFlags)
	}

	// This should be run after the pre block of the script
	preBlock := `
echo "%s"

%s

# Try to support debian / rocky flavor
# This is the weakest point - it takes a long time to install with dnf
/usr/bin/yum install munge -y || apt-get install -y munge || echo "Issue installing munge, might already be installed."
systemctl enable munge || service munge start || echo "Issue starting munge, might already be started."

# Ensure the flux volume addition is complete.
wget -q https://github.com/converged-computing/goshare/releases/download/2023-09-06/wait-fs
chmod +x ./wait-fs
mv ./wait-fs /usr/bin/goshare-wait-fs
	
# Ensure spack view is on the path, wherever it is mounted
viewbase="%s"
viewroot=${viewbase}/view
software="${viewbase}/software"
viewbin="${viewroot}/bin"
fluxpath=${viewbin}/flux

# Important to add AFTER in case software in container duplicated
export PATH=$PATH:${viewbin}
	
# Wait for software directory, and give it time
goshare-wait-fs -p ${software}
	
# Wait for copy to finish
sleep 10
	
# Copy mount software to /opt/software
cp -R ${viewbase}/software /opt/software
	
# Wait for marker (from spack.go) to indicate copy is done
goshare-wait-fs -p ${fluxpath}
goshare-wait-fs -p ${viewbase}/metrics-operator-done.txt

# A small extra wait time to be conservative
sleep 5

# Prefix to run as root (which we will do first)
fluxuser="%s"
fluxuid="%s"

# Add a flux user (required) that should exist before pre-command
# This might vary between OS
# adduser --disabled-password --uid ${fluxuid} --gecos "" ${fluxuser} > /dev/null 2>&1 || echo "Issue adding ${fluxuser}"

# Ensure we use flux's python (TODO update this to use variable)
export PYTHONPATH=${viewroot}/lib/python3.11/site-packages
echo "PYTHONPATH is ${PYTHONPATH}"
echo "PATH is $PATH"

# Add fluxuser to sudoers
echo "${fluxuser} ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Put the state directory in /var/lib on shared view
export STATE_DIR=${viewroot}/var/lib/flux
mkdir -p ${STATE_DIR}

# Main host <name>-0 and the fully qualified domain name
mainHost="%s"

echo "👋 Hello, I'm $(hostname)"
echo "The main host is ${mainHost}"

workdir=$(pwd)
echo "The working directory is ${workdir}, contents include:"
ls -R ${workdir}

brokerOptions="-Scron.directory=/etc/flux/system/cron.d \
  -Stbon.fanout=256 \
  -Srundir=/run/flux %s \
  -Sstatedir=${STATE_DIR} \
  -Slocal-uri=local:///run/flux/local \
  -Stbon.connect_timeout=%s \
  -Sbroker.quorum=%s %s \
  -Slog-stderr-level=%s \
  -Slog-stderr-mode=local"

# Run an interactive cluster, giving no command to flux start
function run_interactive_cluster() {
    echo "🌀 ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions}"
    ${asFlux} flux broker --config-path /etc/flux/config ${brokerOptions}
}

flags="%s"
watch="%s"
submit="%s"

# We will copy the curve certificate if the lead, otherwise wait for it
curvepath=${viewroot}/etc/curve/curve.cert

# Start flux with the original entrypoint
if [ $(hostname) == "${mainHost}" ]; then

  # The main host needs to scp the curve.cert over to the others
  for host in $(cat ./hostlist.txt); do
      if [[ "$host" == "" ]]; then
	      continue
	  fi
	  if [[ "$host" == "${mainHost}" ]]; then
          continue
	  fi
      echo "Copying curve.cert to $host"
	  scp ${curvepath} $host:${curvepath}
  done
  echo "Command provided is: ${command}"
  if [ "${command}" == "" ]; then

    # An interactive job also doesn't require a command
    run_interactive_cluster
    
  else
     # TODO we can add --wrap here if needed
     echo "🌀 Submit Mode: flux start -o --config ${viewroot}/etc/flux/config ${brokerOptions} flux ${submit} ${flags} ${watch} --quiet -vvv ${command}"
     flux start -o --config ${viewroot}/etc/flux/config ${brokerOptions} flux ${submit} ${flags} --quiet ${watch} -vvv ${command}
  fi

# Block run by workers
else

# We basically sleep/wait until the lead broker is ready
echo "🌀 flux start -o --config ${viewroot}/etc/flux/config ${brokerOptions}"
goshare-wait-fs -p ${curvepath}

# We can keep trying forever, don't care if worker is successful or not
while true
  do
    flux start -o --config ${viewroot}/etc/flux/config ${brokerOptions}
    retval=$?
    echo "Return value for follower worker is ${retval}"
    echo "😪 Sleeping 15s to try again..."
    sleep 15
done
fi

echo "%s"
echo "%s"
`
	preBlock = fmt.Sprintf(
		preBlock,
		meta,
		a.preCommand,
		a.Mount,
		a.fluxUser,
		a.fluxUid,
		leadBroker,
		interactive,
		a.connectTimeout,
		a.quorum,
		zeromq,
		a.logLevel,
		flags,
		watch,
		a.submitCommand,
		metadata.CollectionStart,
		metadata.Separator,
	)

	// Flux needs this set to false
	setFQDN := false

	// We use container names to target specific entrypoint scripts here
	for _, containerSpec := range cs {

		// First check - is this the right replicated job?
		if containerSpec.JobName != rj.Name {
			continue
		}
		rj.Template.Spec.Template.Spec.SetHostnameAsFQDN = &setFQDN

		// Always copy over the pre block - we need the logic to copy software
		// Then we need to add the command, and finish with the full preBlock
		// The command is given to flux!
		command := containerSpec.EntrypointScript.Command
		containerSpec.EntrypointScript.Pre += "\n" + fmt.Sprintf("command='%s'", command) + "\n" + preBlock

		// Next check if we have a target set (for the container)
		if a.containerTarget != "" && containerSpec.Name != "" && a.containerTarget != containerSpec.Name {
			continue
		}

		// If the post command ends with sleep infinity, tweak it
		isInteractive, updatedPost := deriveUpdatedPost(containerSpec.EntrypointScript.Post)
		containerSpec.EntrypointScript.Post = updatedPost

		// We will never get to command, so just make it empty
		containerSpec.EntrypointScript.Command = ""

		// If is interactive, add back sleep infinity
		if isInteractive {
			containerSpec.EntrypointScript.Post += "\nsleep infinity\n"
		}
	}
}

func init() {
	base := AddonBase{
		Identifier: fluxIdentifier,
		Summary:    "hierarchical graph-based scheduler and resource manager",
	}
	app := ApplicationAddon{AddonBase: base}
	spack := SpackView{ApplicationAddon: app}
	flux := FluxFramework{SpackView: spack}
	Register(&flux)
}
